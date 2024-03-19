package process

import (
	"github.com/golang-mods/tango/internal/process/gomod"
	"github.com/spf13/afero"
)

func readGoModFile(fs afero.Fs, name string) ([]*gomod.Require, error) {
	file, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return gomod.Decode(file)
}
