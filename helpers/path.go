package helpers

import (
	"os"
	"path/filepath"
)

func GetPathOfSelf() (dir string, err error) {
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}
