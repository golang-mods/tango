package process

import (
	"path/filepath"

	"github.com/golang-mods/tango/internal/constant"
	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/golang-mods/tango/internal/process/binary"
	"github.com/spf13/afero"
)

func ReadBinariesDirectory(fs afero.Fs, command gocommand.Command, current string) ([]*binary.Binary, error) {
	manifestFile, err := findFileTowardAncestors(fs, current, constant.ManifestFileName)
	if err != nil {
		return nil, err
	}
	rootDirectory := filepath.Dir(manifestFile)
	binariesDirectory := filepath.Join(rootDirectory, constant.BinariesDirectoryName)

	return binary.ReadBinariesDirectory(fs, command, binariesDirectory)
}
