package process

import (
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/golang-mods/exerrors"
	"github.com/golang-mods/serrors"
	"github.com/golang-mods/sorted"
	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/golang-mods/tango/internal/process/binary"
	"github.com/golang-mods/tango/internal/process/gomod"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/golang-mods/tango/internal/process/manifestfile"
	"github.com/spf13/afero"
)

var ErrNotFoundReferenceVersion = errors.New("not found reference version")

func applyGoModRequiresToManifestFile(
	logger *slog.Logger,
	requires []*gomod.Require,
	tools []manifest.ManifestTool,
) error {
	referenceTools := make([]*manifest.ManifestTool, 0, len(tools))
	for i := range tools {
		tool := &tools[i]
		if tool.Version.IsReference() {
			referenceTools = append(referenceTools, tool)
		}
	}

	var errs []error
	i, j := 0, 0
	for i < len(requires) && j < len(referenceTools) {
		require := requires[i]
		tool := referenceTools[j]

		if result := manifest.ComapreName(require, tool); result == 0 || strings.HasPrefix(tool.Path, require.Path+"/") {
			if err := tool.Version.UnmarshalText([]byte(require.Version.Original())); err != nil {
				errs = append(errs, err)
			} else {
				logger.Debug("Found reference version",
					"path", tool.Path,
					"reference", require.Path,
					"version", tool.Version,
				)
			}

			j++
		} else if result < 0 {
			i++
		} else {
			errs = append(errs, serrors.Format("%w", ErrNotFoundReferenceVersion)("path", tool.Path))
			j++
		}
	}

	return errors.Join(errs...)
}

func applyManifestFileToManifestLockFile(
	logger *slog.Logger,
	command *gocommand.Memorized,
	manifestFile *manifestfile.ManifestFile[manifest.Manifest],
	manifestLockFile *manifestfile.ManifestFile[manifest.ManifestLock],
) error {
	tools, toolsToAdd, lockTools, lockToolsToRemove := sorted.IntersectionXorWith(
		manifestFile.Manifest().Tools,
		manifestLockFile.Manifest().Tools,
		manifest.ComapreName,
	)

	if err := exerrors.ParallelMap0(tools, func(tool manifest.ManifestTool, i int) error {
		lockTool := &lockTools[i]
		if !tool.Version.Check(&lockTool.Version) {
			version, err := latestVersionThatSatisfiesConstraints(command, tool.Path, &tool.Version)
			if err != nil {
				return err
			}

			lockTool.Version = *version
			manifestLockFile.Updated()
		}

		if !slices.Equal(tool.Env, lockTool.Env) {
			lockTool.Env = tool.Env
			manifestLockFile.Updated()
		}
		return nil
	}); err != nil {
		return err
	}

	if len(toolsToAdd) > 0 {
		addLockTools, err := exerrors.ParallelMap(toolsToAdd, func(tool manifest.ManifestTool, _ int) (manifest.ManifestLockTool, error) {
			version, err := latestVersionThatSatisfiesConstraints(command, tool.Path, &tool.Version)
			if err != nil {
				return manifest.ManifestLockTool{}, err
			}

			return manifest.ManifestLockTool{
				ToolPath: tool.ToolPath,
				Version:  *version,
				Tool:     tool.Tool,
			}, nil
		})
		if err != nil {
			return err
		}

		lockTools = sorted.UnionWith(lockTools, addLockTools, manifest.ComapreName)
		manifestLockFile.Updated()
	}

	if len(lockToolsToRemove) > 0 {
		manifestLockFile.Updated()
	}

	manifestLockFile.Manifest().Tools = lockTools

	return nil
}

func applyManifestLockFileToBinariesDirectoy(
	logger *slog.Logger,
	fs afero.Fs,
	command *gocommand.Memorized,
	directory string,
	manifestLockFile *manifestfile.ManifestFile[manifest.ManifestLock],
	binaries []*binary.Binary,
) error {
	tools, toolsToAdd, intersectionBinaries, binariesToRemove := sorted.IntersectionXorWith(
		manifestLockFile.Manifest().Tools,
		binaries,
		manifest.ComapreName,
	)

	for i, tool := range tools {
		binary := intersectionBinaries[i]

		if !equalToolToBinary(&tool, binary) {
			binariesToRemove = append(binariesToRemove, binary)
			toolsToAdd = append(toolsToAdd, tool)
		}
	}

	slices.SortFunc(binariesToRemove, manifest.ComapreName)
	slices.SortFunc(toolsToAdd, manifest.ComapreName)

	return errors.Join(
		removeBinaries(logger, fs, binariesToRemove),
		installTools(logger, command, directory, toolsToAdd),
	)
}

func removeBinaries(logger *slog.Logger, fs afero.Fs, binaries []*binary.Binary) error {
	return exerrors.Map0(binaries, func(binary *binary.Binary, _ int) error {
		logger.Info("removing", "name", binary.Path)
		return fs.Remove(binary.Version.File)
	})
}

func installTools(
	logger *slog.Logger,
	command *gocommand.Memorized,
	directory string,
	tools []manifest.ManifestLockTool,
) error {
	logDownloading := func(pkg string) {
		logger.Info("downloading", "pacakge", pkg)
	}

	return exerrors.Map0(tools, func(tool manifest.ManifestLockTool, _ int) error {
		pkg := manifest.ToPackage(tool.Path, &tool.Version)
		logger.Info("installing", "package", pkg)

		return command.Install(directory, []string{pkg}, tool.Env, logDownloading)
	})
}

func equalToolToBinary(tool *manifest.ManifestLockTool, binary *binary.Binary) bool {
	return tool.Version.Equal(binary.Version.Module.Version) &&
		sorted.Contains(binary.Version.Builds, tool.Env) &&
		binary.IsValidFileName()
}
