package utils

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
