package manifest

import (
	"io"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
)

type ManifestLock struct {
	Tools []ManifestLockTool `toml:"tools,omitempty"`
}

type ManifestLockTool struct {
	ToolPath
	Version semver.Version `toml:"version"`
	Tool
}

func EncodeManifestLock(writer io.Writer, manifest *ManifestLock) error {
	return toml.NewEncoder(writer).Encode(manifest)
}

func DecodeManifestLock(reader io.Reader) (*ManifestLock, error) {
	var manifest ManifestLock

	if err := toml.NewDecoder(reader).Decode(&manifest); err != nil {
		return nil, err
	}

	if err := SortNamers(manifest.Tools); err != nil {
		return nil, err
	}

	for _, tool := range manifest.Tools {
		slices.Sort(tool.Env)
	}

	return &manifest, nil
}
