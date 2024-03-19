package tango

import (
	"errors"

	"github.com/golang-mods/sorted"
	"github.com/golang-mods/tango/internal/process"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/samber/lo"
)

func (manager *Manager) Uninstall(packages []string) error {
	if len(packages) > 0 {
		manager.logger.Warn("packages is empty")
	}

	pkgs := lo.Map(packages, func(pkg string, _ int) manifest.ManifestTool {
		return manifest.ManifestTool{ToolPath: manifest.NewToolPath(pkg)}
	})

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

		tools, notInstalled := sorted.XorWith(file.Manifest().Tools, pkgs, manifest.ComapreName)
		for _, pkg := range notInstalled {
			manager.logger.Warn("not installed", "package", pkg.Path)
		}

		file.Manifest().Tools = tools
		file.Updated()
	}

	return errors.Join(processer.ProcessManifest(), processer.ProcessBinaries())
}
