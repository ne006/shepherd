package utils

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get caller information")
	}

	tests := []struct {
		name     string
		filePath string
		result   bool
	}{
		{"Non-existing file", fmt.Sprintf("/tmp/inexistent_file_%v", time.Now()), false},
		{"Existing file", filename, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, FileExists(tt.filePath))
		})

	}
}

func TestWritePidFile(t *testing.T) {
	tests := []struct {
		name string

		pidFilePath string
		pid         int
		wantErr     bool
	}{
		{"Valid path", "tmp/test_process.pid", 123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.NotEqual(t, nil, WritePidFile(tt.pidFilePath, tt.pid))
				assert.NoFileExists(t, tt.pidFilePath)
			} else {
				assert.Equal(t, nil, WritePidFile(tt.pidFilePath, tt.pid))
				assert.FileExists(t, tt.pidFilePath)
				if content, err := os.ReadFile(tt.pidFilePath); err != nil {
					t.Errorf("Could not read %s: %s", tt.pidFilePath, err)
				} else {
					assert.Equal(t, string(content), strconv.Itoa(tt.pid))
				}
			}
			os.Remove(tt.pidFilePath)
		})
	}
}
