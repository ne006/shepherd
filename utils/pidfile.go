package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

type Pidfile struct {
	Path string
}

func (pf *Pidfile) Write(pid int) error {
	pidFileDir := filepath.Dir(pf.Path)

	if err := os.MkdirAll(pidFileDir, 0766); err != nil {
		return err
	}

	if err := os.WriteFile(pf.Path, []byte(strconv.Itoa(pid)), 0766); err != nil {
		return err
	}

	return nil
}

func (pf *Pidfile) Remove() error {
	return os.Remove(pf.Path)
}
