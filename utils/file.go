package utils

import (
	"errors"
	"os"
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
