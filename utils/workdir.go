package utils

import (
	"os"
	"path/filepath"
)

func defaultWdPath() string {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return ".shepherd"
	} else {
		return filepath.Join(homeDir, ".shepherd")
	}
}

func NewWorkdir() Workdir {
	return Workdir{
		Path: defaultWdPath(),
	}
}

type Workdir struct {
	Path string
}

func (wd *Workdir) Create() error {
	return os.MkdirAll(wd.Path, 0766)
}
