package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

type Pidfile struct {
	path string
}

func (pf *Pidfile) Write(pid int) error {
	pidFileDir := filepath.Dir(pf.path)

	if err := os.MkdirAll(pidFileDir, 0766); err != nil {
		return err
	}

	if err := os.WriteFile(pf.path, []byte(strconv.Itoa(pid)), 0766); err != nil {
		return err
	}

	return nil
}

func (pf *Pidfile) Remove() error {
	return os.Remove(pf.path)
}
