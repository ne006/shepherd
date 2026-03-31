package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkdir_Create(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Default path", "", false},
		{"Valid path", "tmp/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wd Workdir

			if tt.path == "" {
				wd = NewWorkdir()
			} else {
				wd = NewWorkdir()
				wd.Path = tt.path
			}

			err := wd.Create()

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, true, FileExists(wd.Path))
			}
		})
	}
}
