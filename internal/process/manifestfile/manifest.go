package manifestfile

import (
	"errors"
	"io"
	"os"
	"sync/atomic"

	"github.com/spf13/afero"
)

type EncodeManifest[Manifest any] func(writer io.Writer, manifest *Manifest) error
type DecodeManifest[Manifest any] func(reader io.Reader) (*Manifest, error)

type ManifestFile[Manifest any] struct {
	manifest *Manifest
	encode   EncodeManifest[Manifest]
	file     afero.File
	updated  atomic.Bool
}

func Open[Manifest any](
	fs afero.Fs,
	decode DecodeManifest[Manifest],
	encode EncodeManifest[Manifest],
	name string,
) (*ManifestFile[Manifest], error) {
	return openManifestFile(fs, decode, encode, name, os.O_RDWR)
}

func OpenOrCreate[Manifest any](
	fs afero.Fs,
	decode DecodeManifest[Manifest],
	encode EncodeManifest[Manifest],
	name string,
) (*ManifestFile[Manifest], error) {
	return openManifestFile(fs, decode, encode, name, os.O_RDWR|os.O_CREATE)
}

func openManifestFile[Manifest any](
	fs afero.Fs,
	decode DecodeManifest[Manifest],
	encode EncodeManifest[Manifest],
	name string,
	flag int,
) (*ManifestFile[Manifest], error) {
	file, err := fs.OpenFile(name, flag, 0600)
	if err != nil {
		return nil, err
	}

	manifest, err := decode(file)
	if err != nil {
		return nil, err
	}

	return &ManifestFile[Manifest]{
		manifest: manifest,
		encode:   encode,
		file:     file,
	}, nil
}

func (manifestFile *ManifestFile[Manifest]) Close() error {
	var err error

	if manifestFile.updated.Load() {
		err = errors.Join(
			func() error { _, err := manifestFile.file.Seek(0, 0); return err }(),
			manifestFile.file.Truncate(0),
			manifestFile.encode(manifestFile.file, manifestFile.manifest),
		)
	}

	return errors.Join(err, manifestFile.file.Close())
}

func (manifestFile *ManifestFile[Manifest]) Updated() {
	manifestFile.updated.Store(true)
}

func (manifestFile *ManifestFile[Manifest]) NotUpdated() {
	manifestFile.updated.Store(false)
}

func (manifestFile *ManifestFile[Manifest]) Manifest() *Manifest {
	return manifestFile.manifest
}
