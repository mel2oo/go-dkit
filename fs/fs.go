package fs

import (
	"fmt"
	"os"
)

func Mkdir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		}
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

func MkdirIf(path string, uid, gid int) error {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	if !stat.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	os.Chmod(path, os.ModePerm)
	os.Chown(path, uid, gid)
	return nil
}
