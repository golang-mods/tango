package tango

import (
	"github.com/golang-mods/sorted"
	"github.com/golang-mods/tango/internal/process"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/samber/lo"
)

func (manager *Manager) UpdateAll() error {
	return manager.Update(nil)
}

func (manager *Manager) Update(packages []string) error {
	pkgs := lo.Map(packages, func(pkg string, _ int) manifest.ToolPath { return manifest.NewToolPath(pkg) })

	processer, err := process.Open(manager.logger, fs, command, manager.current)
	if err != nil {
		return err
	}
	defer processer.Close()

	lockFile, err := processer.ManifestLockFile()
	if err != nil {
		return err
	}
	original := lockFile.Manifest().Tools
	if len(pkgs) > 0 {
		tools, notInstalled := sorted.XorWith(original, pkgs, manifest.ComapreName)
		lockFile.Manifest().Tools = tools
		for _, pkg := range notInstalled {
			manager.logger.Warn("not installed", "package", pkg.Path)
		}
	} else {
		lockFile.Manifest().Tools = nil
	}

	if err := processer.ProcessManifest(); err != nil {
		return err
	}

	tools, _ := sorted.IntersectionWith(lockFile.Manifest().Tools, pkgs, manifest.ComapreName)
	oldTools, newTools := sorted.IntersectionWith(original, tools, manifest.ComapreName)
	for i, oldTool := range oldTools {
		newTool := newTools[i]
		if oldTool.Version.Equal(&newTool.Version) {
			manager.logger.Info("already updated", "path", oldTool.Path, "version", oldTool.Version)
		} else {
			manager.logger.Info("update", "path", oldTool.Path, "old", oldTool.Version, "new", newTool.Version)
		}
	}

	return processer.ProcessBinaries()
}
