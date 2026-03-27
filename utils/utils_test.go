package utils

import (
	"fmt"
	"runtime"
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.result, FileExists(testCase.filePath))
		})

	}
}
