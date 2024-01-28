package internal

import (
	"os"
	"path/filepath"
)

var CacheDir = filepath.Join(os.Getenv("HOME"), ".cache/booster")
