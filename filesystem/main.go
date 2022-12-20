package filesystem

import (
	"os"
)

const (
	RootPath = "."
	DataPath = "./data"
)

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0777)

	if err == nil {
		return nil
	} else {
		return err
	}
}
