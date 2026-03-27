package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			return false
		}
	} else {
		return true
	}
}

func WritePidFile(pidFilePath string, pid int) error {
	pidFileDir := filepath.Dir(pidFilePath)

	if err := os.MkdirAll(pidFileDir, 0766); err != nil {
		return err
	}

	if err := os.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0766); err != nil {
		return err
	}

	return nil
}
