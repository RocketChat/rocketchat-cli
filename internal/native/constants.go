package native

import (
	"os"
	"path"
)

var SystemPathPrefix = ""

func init() {
	SystemPathPrefix = path.Join(os.Getenv("HOME"), ".roctl")
}
