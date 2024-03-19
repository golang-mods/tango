package binary

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/spf13/afero"
)

type Binary struct {
	manifest.ToolPath
	Version *gocommand.BinaryVersion
}

func (binary *Binary) IsValidFileName() bool {
	fileName := filepath.Base(binary.Version.File)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	name := path.Base(binary.Version.Path)

	return fileName == name
}

func ReadBinariesDirectory(fs afero.Fs, command gocommand.Command, directory string) ([]*Binary, error) {
	files, err := afero.ReadDir(fs, directory)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	binaries, err := exerrors.ParallelMap(files, func(file os.FileInfo, _ int) (*Binary, error) {
		name := filepath.Join(directory, file.Name())
		version, err := command.BinaryVersion(name)
		if err != nil {
			return nil, err
		}

		return &Binary{
			ToolPath: manifest.NewToolPath(version.Path),
			Version:  version,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	if err := manifest.SortNamers(binaries); err != nil {
		return nil, err
	}

	return binaries, nil
}
