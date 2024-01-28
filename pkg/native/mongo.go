package native

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rocketchat/booster/internal"
)

type M struct {
	version    string
	installDir string

	binPath string

	cache VersionList

	repl replProcess
}

type replProcess struct {
	process *os.Process

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	wait func() error
}

type VersionList struct {
	Versions []struct {
		Version          string             `json:"version"`
		Lts              bool               `json:"lts_release"`
		ReleaseCandidate bool               `json:"release_candidate"`
		Production       bool               `json:"production_release"`
		Downloads        []DownloadArtifact `json:"downloads"`
	} `json:"versions"`
}

type DownloadArtifact struct {
	Arch    string `json:"arch"`
	Edition string `json:"edition"`
	Target  string `json:"target"`
	Archive struct {
		SHA256 string `json:"sha256"`
		Url    string `json:"url"`
	} `json:"archive"`
}

const versionsUrl = "https://downloads.mongodb.org/full.json"

const versionsFile = "mongo_versions.json"

var ErrArchitectureUnsupported = errors.New("mongodb architecture unknown")
var ErrMongodbVersionNotFound = errors.New("mongodb version not found")

func NewMongoManager(version, rootDir string) (*M, error) {
	installDir := filepath.Join(rootDir, "m", version, runtime.GOOS, runtime.GOARCH)
	m := &M{
		installDir: installDir,
		version:    version,
		binPath:    filepath.Join(installDir, "bin", "mongo"),
	}
	return m, m.init()
}

func (m *M) Install() error {
	var found bool

	target, err := m.getTarget()
	if err != nil {
		return err
	}

	for _, release := range m.cache.Versions {
		if m.version == release.Version {
			found = true
			for _, download := range release.Downloads {
				if download.Edition == "targeted" && download.Target == target && download.Arch == internal.ARCH {
					// handle this download
					return m.install(download)
				}
			}
		}
	}

	if found {
		return ErrArchitectureUnsupported
	}

	return ErrMongodbVersionNotFound
}

func (m *M) install(release DownloadArtifact) error {
	dir, err := internal.MkdirTemp()
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	archive := filepath.Join(dir, "mongo.tgz")

	if err = internal.DownloadSilent(release.Archive.Url, archive); err != nil {
		return err
	}

	f, err := os.Open(archive)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(data)

	if fmt.Sprintf("%x", hash) != release.Archive.SHA256 {
		return errors.New("downloaded file's hash does not match what is given")
	}

	cmd := exec.Command("tar", "xzf", archive, "--strip-components=1")
	cmd.Dir = dir

	if err = cmd.Run(); err != nil {
		return err
	}

	return internal.DumbInstall(filepath.Join(m.installDir, "bin"), filepath.Join(dir, "bin"))
}

// init sets up th cache and initial data
func (m *M) init() error {
	if stat, err := os.Stat(internal.CacheDir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(internal.CacheDir, 0750)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !stat.IsDir() {
			return errors.New("cache can not be saved since is not a directory")
		}
	}

	cacheFilename := filepath.Join(internal.CacheDir, versionsFile)

	var (
		cacheExists bool = true
		stat        fs.FileInfo
		err         error
	)

	if stat, err = os.Stat(cacheFilename); err != nil {
		if os.IsNotExist(err) {
			cacheExists = false
		} else {
			return err
		}
	}

	cacheFile, err := os.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0750)
	if err != nil {
		return err
	}

	var data VersionList

	if cacheExists && time.Since(stat.ModTime()) < (time.Hour*24) {
		if err = json.NewDecoder(cacheFile).Decode(&data); err != nil {
			return err
		}

		m.cache = data
		return nil
	}

	resp, err := http.Get(versionsUrl)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	_, err = cacheFile.Write(content)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return err
	}

	m.cache = data

	return nil
}

func (m *M) getTarget() (string, error) {
	os, err := internal.OsRelease()
	if err != nil {
		return "", err
	}

	switch os.Id {
	case "ubuntu":
		switch os.VersionId {
		case "20.04", "18.04", "18.10", "20.10", "21.04", "21.10", "22.04", "22.10":
			return os.Id + strings.ReplaceAll(os.VersionId, ".", ""), nil
		}
	case "debian":
		fallthrough
	case "rhel":
		return os.Id + os.VersionId, nil
	}

	return "", errors.New("not supported target")
}

func (m *M) Repl(code string) {
	m.repl.stdin.Write([]byte(code))
}

func (m *M) ReplStart() error {
	if m.repl.process != nil {
		return nil
	}

	var err error

	m.repl.stdin, m.repl.stdout, m.repl.stderr, m.repl.wait, err = m.runOnce()
	if err != nil {
		return err
	}

	go func() {
		var p []byte = make([]byte, 1024)
		for {
			_, _ = m.repl.stdout.Read(p)
			fmt.Print(string(p))
		}
	}()

	go func() {
		var p []byte = make([]byte, 1024)
		for {
			_, _ = m.repl.stderr.Read(p)
			fmt.Print(string(p))
		}
	}()

	return nil
}

type MongodRunner struct {
	m *M

	wait func() error
}

func (m *M) run(args ...string) (stdout []byte, stderr []byte, err error) {
	stdout, stderr, err = internal.RunCommand(append([]string{m.binPath}, args...), nil)

	return
}

func (m *M) runOnce(args ...string) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, wait internal.WaitCommandFn, err error) {
	stdin, stdout, stderr, wait, err = internal.RunCommandOnce(append([]string{m.binPath}, args...), nil)

	return
}

func (m *M) Mongod() *MongodRunner {
	mongod := *m
	mongod.binPath = filepath.Join(mongod.installDir, "bin", "mongod")
	return &MongodRunner{m: &mongod}
}

func (m *MongodRunner) Run(args ...string) (stdout []byte, stderr []byte, err error) {
	return m.m.run(args...)
}

func (m *MongodRunner) RunOnce(args ...string) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, wait internal.WaitCommandFn, err error) {
	return m.m.runOnce(args...)
}

func (m *MongodRunner) Version() (string, error) {
	out, serr, err := m.Run("--version")
	if err != nil {
		return "", fmt.Errorf("failed to get mongo version: stderr: %s, err: %v", serr, err)
	}

	var v []byte

	// "db version vX.Y.Z"
	for _, b := range out[12:] {
		if b == '\n' {
			break
		}

		v = append(v, b)
	}

	return string(v), err
}

func (m *MongodRunner) Start() error {
	// there won't be any configs today
	s, err := os.Stat("/data/db")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err = os.MkdirAll("/data/db", 0750); err != nil {
			return err
		}
	} else if !s.IsDir() {
		return errors.New("db dir is a file")
	}

	_, _, _, wait, err := m.m.runOnce("--dbpath", "/data/db", "--pidfilepath", filepath.Join(m.m.installDir, "mongod.pid"))
	if err != nil {
		return err
	}

	m.wait = wait

	return nil
}

func (m *MongodRunner) Stop() error {
	if m.wait == nil {
		return errors.New("mongod isn't running")
	}

	pidfile := filepath.Join(m.m.installDir, "mongod.pid")

	var err error

	defer func() {
		if err == nil {
			os.Remove(pidfile)
		}
	}()

	f, err := os.Open(pidfile)
	if err != nil {
		return err
	}

	pidStr, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	pid, _ := strconv.Atoi(string(pidStr))

	err = syscall.Kill(pid, syscall.SIGINT)
	if err != nil {
		return err
	}

	err = m.wait()

	m.wait = nil

	return err
}
