package spotbugs

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

	sb := Process(files)
	assert.Equal(t, len(sb), 2)
	assert.Equal(t, sb[0].FindBugsSummary.TotalBugs, 4)
	assert.Equal(t, sb[1].FindBugsSummary.TotalBugs, 0)
}
