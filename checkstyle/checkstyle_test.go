package checkstyle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Process(t *testing.T) {
	files := []string{}
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	cs := Process(files)
	assert.Equal(t, len(cs), 5)
}
