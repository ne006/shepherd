package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkdir_Create(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		setEnvVar bool
		wantErr   bool
	}{
		{"Default path", "", false, false},
		{"Default path from env", "tmp/shepherd", true, false},
		{"Valid path", "tmp/", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wd Workdir

			if tt.path == "" {
				wd = NewWorkdir()
			} else {
				if tt.setEnvVar {
					t.Setenv("SHEPHERD_WORKDIR", tt.path)
					wd = NewWorkdir()
				} else {
					wd = NewWorkdir()
					wd.Path = tt.path
				}
			}

			needsCleanup := true

			if FileExists(wd.Path) {
				needsCleanup = false
			}

			err := wd.Create()

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, true, FileExists(wd.Path))
			}

			if needsCleanup {
				os.Remove(wd.Path)
			}
		})
	}
}
