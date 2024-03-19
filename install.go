package tango

import (
	"errors"

	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/sorted"
	"github.com/golang-mods/tango/internal/process"
	"github.com/golang-mods/tango/internal/process/manifest"
)

type InstallOption func(*installOptions)

func InstallOptionEnv(env []string) InstallOption {
	return func(options *installOptions) {
		options.env = env
	}
}

type installOptions struct {
	env []string
}

func createInstallOptions(options ...InstallOption) *installOptions {
	result := &installOptions{}
	for _, apply := range options {
		apply(result)
	}

	return result
}

func (manager *Manager) InstallAll() error {
	return manager.Install(nil)
}

func (manager *Manager) Install(packages []string, options ...InstallOption) error {
	pkgs, err := exerrors.Map(packages, func(pkg string, _ int) (manifest.ManifestTool, error) {
		return manifest.ParseArgument(command, pkg)
	})
	if err != nil {
		return err
	}
	if err := manifest.SortNamers(pkgs); err != nil {
		return err
	}

	internalOptions := createInstallOptions(options...)
	if len(internalOptions.env) > 0 {
		for _, pkg := range pkgs {
			pkg.Env = internalOptions.env
		}
	}

	processer, err := process.Open(manager.logger, fs, command, manager.current)
	if err != nil {
		return err
	}
	defer processer.Close()

	if len(pkgs) > 0 {
		file, err := processer.ManifestFile()
		if err != nil {
			return err
		}

		dupTools, _, dupPkgs, addPkgs := sorted.IntersectionXorWith(file.Manifest().Tools, pkgs, manifest.ComapreName)
		for i := range len(dupPkgs) {
			tool := &dupTools[i]
			pkg := dupPkgs[i]
			tool.Version = pkg.Version
			tool.Env = pkg.Env
		}

		file.Manifest().Tools = sorted.UnionWith(dupTools, addPkgs, manifest.ComapreName)
		file.Updated()
	}

	return errors.Join(processer.ProcessManifest(), processer.ProcessBinaries())
}
