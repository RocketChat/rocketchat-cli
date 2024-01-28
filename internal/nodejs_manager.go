package internal

/*
manages nodejs installations
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/rocketchat/booster/internal/progress"
	"github.com/mholt/archiver"
)

type nodeManager struct {
	dst         string
	l           *Logger
	environment map[string]string
}

func NewNodeManager(dir string, l *Logger) *nodeManager {
	return &nodeManager{
		dst: dir,
		l:   l,
	}
}

func (m *nodeManager) Download(version string) error {
	url := m.buildUrl(version)

	m.l.Infof("downloading nodejs version %s", version)
	m.l.Debugf("url: %s", url)

	tmpdir, err := os.MkdirTemp("", "booster")
	if err != nil {
		m.l.Errorf(err, "failed to create temporary directory for file download")
		return err
	}

	m.l.Debugf("temp directory %s", tmpdir)

	filename := path.Join(tmpdir, path.Base(url))

	m.l.Debugf("downloading to file: %s", filename)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	m.l.Debugf("downloading %d bytes", resp.ContentLength)

	defer resp.Body.Close()

	bar, err := progress.NewWriteProgressBar("Downloading nodejs", resp.ContentLength, f, &progress.ProgressBarOptions{BarWidth: 70})
	if err != nil {
		return err
	}

	io.Copy(bar, resp.Body)

	m.l.Infof("nodejs downloading finished")
	m.l.Infof("installing nodejs")

	dir := path.Join(m.dst, "nodejs", version)
	if err = os.MkdirAll(dir, 0777); err != nil {
		m.l.Errorf(err, "failed to create nodejs installation directory")
		return err
	}

	if err = archiver.Unarchive(filename, dir); err != nil {
		m.l.Errorf(err, "failed to extract nodejs archive")
		return err
	}
	m.l.Infof("nodejs successfully installed")

	return nil
}

func (m *nodeManager) buildUrl(version string) string {
	return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-linux-%s.tar.xz", version, version, runtime.GOARCH)

}
