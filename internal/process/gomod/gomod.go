package gomod

import (
	"io"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/tango/internal/process/manifest"
	"golang.org/x/mod/modfile"
)

type Require struct {
	manifest.ToolPath
	Version *semver.Version
}

func Decode(reader io.Reader) ([]*Require, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	file, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, err
	}

	requires, err := exerrors.Map(file.Require, func(require *modfile.Require, _ int) (*Require, error) {
		version, err := semver.NewVersion(require.Mod.Version)
		if err != nil {
			return nil, err
		}

		return &Require{
			ToolPath: manifest.NewToolPath(require.Mod.Path),
			Version:  version,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	if err := manifest.SortNamers(requires); err != nil {
		return nil, err
	}

	return requires, nil
}
