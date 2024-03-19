package gocommand

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-mods/exerrors"
	"github.com/spf13/afero"
)

var _ Command = (*Mock)(nil)

var extension = ""

type Mock struct {
	Fs afero.Fs
}

func NewMock(fs afero.Fs) Command { return &Mock{Fs: fs} }

func (mock Mock) Install(directory string, packages []string, env []string, logDownloading func(string)) error {
	_, err := mock.Fs.Stat(directory)
	if os.IsNotExist(err) {
		if err := mock.Fs.MkdirAll(directory, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return exerrors.Map0(packages, func(pkg string, _ int) error {
		slashIndex := strings.LastIndex(pkg, "/")

		name := filepath.Join(directory, pkg[slashIndex+1:])
		file, err := mock.Fs.Create(name + extension)
		if err != nil {
			return err
		}
		return file.Close()
	})
}

func (mock Mock) Versions(pkg string) (*Versions, error) {
	return &Versions{
		Path:     pkg,
		Versions: []*semver.Version{semver.MustParse("v1.0.0")},
	}, nil
}

func (Mock) BinaryVersion(name string) (*BinaryVersion, error) {
	return &BinaryVersion{
		Path: "mock/mock",
		Module: &Module{
			Name:    "mock/mock",
			Version: semver.MustParse("v1.0.0"),
		},
	}, nil
}
