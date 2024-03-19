package manifest

import (
	"path"

	"github.com/Masterminds/semver/v3"
)

type Tool struct {
	Env []string `toml:"env,omitempty"`
}

type ToolPath struct {
	Path string `toml:"path"`
}

func NewToolPath(path string) ToolPath {
	return ToolPath{Path: path}
}

func (p ToolPath) Name() string {
	return p.Path
}

func (p ToolPath) BinaryName() string {
	return path.Base(p.Path)
}

func ToPackage(path string, version *semver.Version) string {
	return path + "@v" + version.String()
}
