package internal

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func DownloadSilent(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("nodejs version not found")
	}

	dir := filepath.Dir(path)

	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}

		err := os.MkdirAll(dir, 0750)
		if err != nil {
			return err
		}
	}

	w, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0750)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return w.Sync()
}

func MkdirTemp() (string, error) {
	dir, err := os.MkdirTemp("", "booster")
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(dir)
}

// not a full implementation of install, of course
func DumbInstall(dst, src string) (err error) {
	var links map[string]string = make(map[string]string)

	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)

		info, _ := d.Info()

		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			linkSrc, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			if relSrc, err := filepath.Rel(src, linkSrc); err == nil {
				links[filepath.Join(dst, relSrc)] = filepath.Join(dst, rel)
				return nil
			}
		}

		if d.IsDir() {
			if err := os.MkdirAll(filepath.Join(dst, rel), info.Mode().Perm()); err != nil {
				return err
			}
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}

		defer srcFile.Close()

		dstFile, err := os.OpenFile(filepath.Join(dst, rel), os.O_CREATE|os.O_WRONLY, info.Mode().Perm())
		if err != nil {
			return err
		}

		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}

		return dstFile.Sync()
	})

	if err != nil {
		return
	}

	for src, dst := range links {
		if err = os.Symlink(src, dst); err != nil && os.IsExist(err) {
			os.Remove(dst)
			err = os.Symlink(src, dst)
			if err != nil {
				break
			}
		} else if err != nil {
			break
		}
	}

	return
}

func removeSurroundingQuotes(str string) string {
	if str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}

	return str
}

type _OsRelease struct {
	Name            string
	Id              string
	Version         string
	VersionId       string
	IdLike          string
	PrettyName      string
	VersionCodename string
}

/*
NAME="Ubuntu"
VERSION="20.04.6 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.6 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal
*/
func OsRelease() (_OsRelease, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return _OsRelease{}, err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return _OsRelease{}, err
	}

	str := string(content)

	var o _OsRelease

	for _, line := range strings.Split(str, "\n") {
		keywords := strings.Split(line, "=")
		switch keywords[0] {
		case "NAME":
			o.Name = removeSurroundingQuotes(keywords[1])
		case "VERSION":
			o.Version = removeSurroundingQuotes(keywords[1])
		case "ID":
			o.Id = removeSurroundingQuotes(keywords[1])
		case "ID_LIKE":
			o.IdLike = removeSurroundingQuotes(keywords[1])
		case "PRETTY_NAME":
			o.PrettyName = removeSurroundingQuotes(keywords[1])
		case "VERSION_ID":
			o.VersionId = removeSurroundingQuotes(keywords[1])
		case "VERSION_CODENAME":
			o.VersionCodename = removeSurroundingQuotes(keywords[1])
		}
	}

	return o, nil
}

var ARCH = map[string]string{"amd64": "x86_64", "arm64": "aarch64"}[runtime.GOARCH]

func RunCommand(args []string, environment []string) (stdout []byte, stderr []byte, err error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = environment

	stdout, err = cmd.Output()
	if err != nil {
		e, _ := err.(*exec.ExitError)
		if e == nil {
			return
		}
		stderr = e.Stderr
		return
	}

	return
}

type WaitCommandFn func() error

func RunCommandOnce(args []string, environment []string) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, wait WaitCommandFn, err error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = environment

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()

	wait = cmd.Wait

	return
}
