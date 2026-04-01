package utils

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPidfile_Write(t *testing.T) {
	tests := []struct {
		name string

		path    string
		pid     int
		wantErr bool
	}{
		{"Valid path", "tmp/test_process.pid", 123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := Pidfile{Path: tt.path}

			if tt.wantErr {
				assert.NotEqual(t, nil, pf.Write(tt.pid))
				assert.NoFileExists(t, tt.path)
			} else {
				assert.Equal(t, nil, pf.Write(tt.pid))
				assert.FileExists(t, tt.path)
				if content, err := os.ReadFile(tt.path); err != nil {
					t.Errorf("Could not read %s: %s", tt.path, err)
				} else {
					assert.Equal(t, string(content), strconv.Itoa(tt.pid))
				}
			}
			os.Remove(tt.path)
		})
	}
}

func TestPidfile_Remove(t *testing.T) {
	tests := []struct {
		name string

		path    string
		wantErr bool
	}{
		{"Valid path", "tmp/test_process.pid", false},
		{"Invalid path", "tmp/non-existent-file", true},
	}

	os.WriteFile("tmp/test_process.pid", []byte("123"), 0766)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := Pidfile{Path: tt.path}

			err := pf.Remove()

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoFileExists(t, tt.path)
			}
		})
	}
}
