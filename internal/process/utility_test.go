package process

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFindFileTowardAncestors(t *testing.T) {
	fs := afero.NewMemMapFs()

	name := "file"
	current := filepath.Join("a", "b")
	bFile := filepath.Join(current, name)
	aFile := filepath.Join(filepath.Dir(current), name)

	assert.Nil(t, fs.MkdirAll(current, 0600))

	{
		root, err := findFileTowardAncestors(fs, current, name)
		assert.Empty(t, root)
		assert.Error(t, err)
	}

	{
		file, err := fs.Create(aFile)
		assert.Nil(t, err)
		file.Close()

		root, err := findFileTowardAncestors(fs, current, name)
		assert.Nil(t, err)
		assert.Equal(t, aFile, root)
	}

	{
		file, err := fs.Create(bFile)
		assert.Nil(t, err)
		file.Close()

		root, err := findFileTowardAncestors(fs, current, name)
		assert.Nil(t, err)
		assert.Equal(t, bFile, root)
	}
}
