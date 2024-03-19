package gocommand

import "github.com/Masterminds/semver/v3"

type Command interface {
	Install(directory string, packages []string, env []string, logDownloading func(string)) error
	Versions(pkg string) (*Versions, error)
	BinaryVersion(name string) (*BinaryVersion, error)
}

type Versions struct {
	Path     string
	Versions []*semver.Version
}

type BinaryVersion struct {
	File         string
	GoVersion    *semver.Version
	Path         string
	Module       *Module
	Dependencies []*Module
	Builds       []string
}

type Module struct {
	Name    string
	Version *semver.Version
	Hash    string
}
