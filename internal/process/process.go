package process

import (
	"errors"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/golang-mods/tango/internal/constant"
	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/golang-mods/tango/internal/memorize"
	"github.com/golang-mods/tango/internal/process/binary"
	"github.com/golang-mods/tango/internal/process/gomod"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/golang-mods/tango/internal/process/manifestfile"
	"github.com/spf13/afero"
)

type Processer struct {
	logger  *slog.Logger
	fs      afero.Fs
	command *gocommand.Memorized

	binariesDirectory string
	manifestFile      func() (*manifestfile.ManifestFile[manifest.Manifest], error)
	goModRequires     func() ([]*gomod.Require, error)
	manifestLockFile  func() (*manifestfile.ManifestFile[manifest.ManifestLock], error)
	binaries          func() ([]*binary.Binary, error)
}

func Open(
	logger *slog.Logger,
	fs afero.Fs,
	command *gocommand.Memorized,
	current string,
) (*Processer, error) {
	// find paths
	manifestFile, err := findFileTowardAncestors(fs, current, constant.ManifestFileName)
	if err != nil {
		return nil, err
	}
	rootDirectory := filepath.Dir(manifestFile)
	manifestLockFile := filepath.Join(rootDirectory, constant.ManifestLockFileName)
	binariesDirectory := filepath.Join(rootDirectory, constant.BinariesDirectoryName)

	// read in parallel
	manifestFileAsync := memorize.Async(func() (*manifestfile.ManifestFile[manifest.Manifest], error) {
		return manifestfile.Open(fs, manifest.DecodeManifest, manifest.EncodeManifest, manifestFile)
	})
	goModRequiresAsync := memorize.Async(func() ([]*gomod.Require, error) {
		name, err := findFileTowardAncestors(fs, rootDirectory, constant.GoModFile)
		if errors.Is(err, ErrNotFoundFile) {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		return readGoModFile(fs, name)
	})
	manifestLockFileAsync := memorize.Async(func() (*manifestfile.ManifestFile[manifest.ManifestLock], error) {
		return manifestfile.OpenOrCreate(
			fs,
			manifest.DecodeManifestLock,
			manifest.EncodeManifestLock,
			manifestLockFile,
		)
	})
	binariesAsync := memorize.Async(func() ([]*binary.Binary, error) {
		return binary.ReadBinariesDirectory(fs, command, binariesDirectory)
	})

	return &Processer{
		logger:            logger,
		fs:                fs,
		command:           command,
		binariesDirectory: binariesDirectory,
		manifestFile:      manifestFileAsync,
		goModRequires:     goModRequiresAsync,
		manifestLockFile:  manifestLockFileAsync,
		binaries:          binariesAsync,
	}, nil
}

func (processer *Processer) Close() error {
	return errors.Join(asyncClose(processer.manifestFile), asyncClose(processer.manifestLockFile))
}

func asyncClose[T io.Closer](async func() (T, error)) error {
	if file, err := async(); err != nil {
		return err
	} else {
		return file.Close()
	}
}

func (processer *Processer) ProcessManifest() error {
	requires, err := processer.goModRequires()
	if err != nil {
		return err
	}
	manifestFile, err := processer.manifestFile()
	if err != nil {
		return err
	}
	if err := applyGoModRequiresToManifestFile(
		processer.logger,
		requires,
		manifestFile.Manifest().Tools,
	); err != nil {
		manifestFile.NotUpdated()
		return err
	}

	manifestLockFile, err := processer.manifestLockFile()
	if err != nil {
		return err
	}
	if err := applyManifestFileToManifestLockFile(
		processer.logger,
		processer.command,
		manifestFile,
		manifestLockFile,
	); err != nil {
		manifestFile.NotUpdated()
		manifestLockFile.NotUpdated()
		return err
	}

	return nil
}

func (processer *Processer) ProcessBinaries() error {
	manifestLockFile, err := processer.manifestLockFile()
	if err != nil {
		return err
	}
	binaries, err := processer.binaries()
	if err != nil {
		return err
	}
	return applyManifestLockFileToBinariesDirectoy(
		processer.logger,
		processer.fs,
		processer.command,
		processer.binariesDirectory,
		manifestLockFile,
		binaries,
	)
}

func (processer *Processer) ManifestFile() (*manifestfile.ManifestFile[manifest.Manifest], error) {
	return processer.manifestFile()
}

func (processer *Processer) ManifestLockFile() (*manifestfile.ManifestFile[manifest.ManifestLock], error) {
	return processer.manifestLockFile()
}
