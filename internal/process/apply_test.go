package process

import (
	"io"
	"log/slog"
	"slices"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-mods/tango/internal/process/gomod"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestApplyGoModRequiresToManifestFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	requires := []*gomod.Require{
		{
			ToolPath: manifest.NewToolPath("example.com/aaa/cmd"),
			Version:  lo.Must(semver.NewVersion("v1")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/bbb/cmd"),
			Version:  lo.Must(semver.NewVersion("v2")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/ccc/cmd"),
			Version:  lo.Must(semver.NewVersion("v3")),
		},
	}

	toolsWithVersion := []manifest.ManifestTool{
		{
			ToolPath: manifest.NewToolPath("example.com/aaa/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("v1")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/bbb/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("v2")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/ccc/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("v3")),
		},
	}

	tools := []manifest.ManifestTool{
		{
			ToolPath: manifest.NewToolPath("example.com/aaa/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("reference")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/bbb/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("reference")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/ccc/cmd"),
			Version:  *lo.Must(manifest.NewConstraints("reference")),
		},
	}

	expected := []manifest.ManifestTool{
		{
			ToolPath: manifest.NewToolPath("example.com/aaa/cmd"),
			Version:  *lo.Must(manifest.NewReferenceConstraintsWithVersion("v1")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/bbb/cmd"),
			Version:  *lo.Must(manifest.NewReferenceConstraintsWithVersion("v2")),
		},
		{
			ToolPath: manifest.NewToolPath("example.com/ccc/cmd"),
			Version:  *lo.Must(manifest.NewReferenceConstraintsWithVersion("v3")),
		},
	}

	testCases := []struct {
		desc     string
		requires []*gomod.Require
		tools    []manifest.ManifestTool
		expected []manifest.ManifestTool
		err      error
	}{
		{
			desc:     "one",
			requires: []*gomod.Require{requires[0]},
			tools:    []manifest.ManifestTool{tools[0]},
			expected: []manifest.ManifestTool{expected[0]},
			err:      nil,
		},
		{
			desc:     "many",
			requires: requires,
			tools:    slices.Clone(tools),
			expected: expected,
			err:      nil,
		},
		{
			desc:     "requires one",
			requires: []*gomod.Require{requires[1]},
			tools:    []manifest.ManifestTool{toolsWithVersion[0], tools[1], toolsWithVersion[2]},
			expected: []manifest.ManifestTool{toolsWithVersion[0], expected[1], toolsWithVersion[2]},
			err:      nil,
		},
		{
			desc:     "tools one",
			requires: requires,
			tools:    []manifest.ManifestTool{tools[1]},
			expected: []manifest.ManifestTool{expected[1]},
			err:      nil,
		},
		{
			desc: "empty",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.ErrorIs(t, tC.err, applyGoModRequiresToManifestFile(logger, tC.requires, tC.tools))
			assert.Equal(t, tC.expected, tC.tools)
		})
	}
}
