package fs

import (
	"fmt"
	"os"
)

func Mkdir(path string, mode os.FileMode) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, mode)
		}
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

func MkdirIf(path string, mode os.FileMode, uid, gid int) error {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(path, mode); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	} else {
		if !stat.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
	}

	os.Chmod(path, mode)
	os.Chown(path, uid, gid)
	return nil
}
