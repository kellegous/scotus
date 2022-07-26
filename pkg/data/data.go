package data

import (
	"io/fs"
	"os"
)

func EnsureDir(
	dir string,
	perm fs.FileMode,
	reset bool,
) error {
	if _, err := os.Stat(dir); err != nil {
		return os.MkdirAll(dir, perm)
	} else if reset {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
		return os.MkdirAll(dir, perm)
	}
	return nil
}
