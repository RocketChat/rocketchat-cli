package native

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rocketchat/booster/internal"
)

var ErrNodeNotInstalled = errors.New("nodejs not installed")

// the node version manager

type N struct {
	version     string
	arch        string
	installDir  string
	environment []string

	binPath string
	global  bool
}

type NpmRunner interface {
	Run(args ...string) ([]byte, []byte, error)
}

func NewNodeManager(global bool, version string, rootDir string) (*N, error) {

	if version[0] != 'v' {
		version = "v" + version
	}

	n := &N{
		global:  global,
		version: version,
		arch:    runtime.GOARCH,
	}

	n.installDir = filepath.Join(rootDir, "n", version, runtime.GOOS, runtime.GOARCH)

	return n, n.findInstall()
}

func (n *N) Npm() NpmRunner {
	npm := *n
	npm.binPath = filepath.Join(filepath.Dir(n.binPath), "npm")
	return &npm
}

func (n *N) Install() error {
	tmpDir, err := internal.MkdirTemp()
	if err != nil {
		return fmt.Errorf("failed to create temporary directory to install nodejs: %v", err)
	}

	url, filename := n._assets()

	archivePath := filepath.Join(tmpDir, filename)
	fmt.Println(tmpDir)

	err = internal.DownloadSilent(url, archivePath)
	if err != nil {
		return err
	}

	cmd := exec.Command("xz", "--decompress", archivePath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	archivePath = strings.TrimSuffix(archivePath, ".xz")

	cmd = exec.Command("tar", "xf", archivePath)
	cmd.Dir = tmpDir
	if err = cmd.Run(); err != nil {
		return err
	}

	archivePath = strings.TrimSuffix(archivePath, ".tar")

	toInstall := []string{"share", "lib", "include", "bin"}

	var dst = n.installDir

	if n.global {
		dst = "/usr/local"
	}

	for _, loc := range toInstall {
		if err := internal.DumbInstall(filepath.Join(dst, loc), filepath.Join(archivePath, loc)); err != nil {
			return err
		}
	}

	return nil
}

func (n *N) EnsureInstalled() error {
	if n.version == n.Version() {
		return nil
	}

	return n.Install()
}

func (n *N) findInstall() error {
	if n.global {
		binPath, err := exec.LookPath("node")
		if err != nil && errors.Is(err, exec.ErrNotFound) {
			return ErrNodeNotInstalled
		} else if err != nil {
			return fmt.Errorf("unknown error trying to detect nodejs global installation: %v", err)
		}

		n.binPath = binPath

		return nil
	}

	n.binPath = filepath.Join(n.installDir, "bin/node")

	path := os.Getenv("PATH")

	n.environment = append(n.environment, fmt.Sprintf("PATH=%s:%s", filepath.Dir(n.binPath), path))

	return nil
}

func (n *N) Run(args ...string) (stdout []byte, stderr []byte, err error) {
	stdout, stderr, err = internal.RunCommand(append([]string{n.binPath}, args...), n.environment)

	return
}

func (n *N) RunOnce(args ...string) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, wait internal.WaitCommandFn, err error) {
	stdin, stdout, stderr, wait, err = internal.RunCommandOnce(append([]string{n.binPath}, args...), n.environment)

	return
}

func (n *N) Version() string {
	out, _, _ := n.Run("--version")
	if len(out) == 0 {
		return string(out)
	}
	return string(out[:len(out)-1])
}

func (n *N) _assets() (url string, filename string) {
	filename = fmt.Sprintf("node-%s-%s-%s.tar.xz", n.version, runtime.GOOS, runtime.GOARCH)
	url = fmt.Sprintf("https://nodejs.org/download/release/%s/%s", n.version, filename)
	return
}
