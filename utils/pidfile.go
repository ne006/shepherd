package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

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
