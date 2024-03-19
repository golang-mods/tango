package process

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/golang-mods/serrors"
	"github.com/golang-mods/tango/internal/gocommand"
	"github.com/golang-mods/tango/internal/process/manifest"
	"github.com/samber/lo"
	"github.com/spf13/afero"
)

var ErrNotFoundFile = errors.New("file is not found")

func findFileTowardAncestors(fs afero.Fs, current string, name string) (string, error) {
	for previous := ""; previous != current; current, previous = filepath.Dir(current), current {
		name := filepath.Join(current, name)
		if info, err := fs.Stat(name); err == nil && !info.IsDir() {
			return name, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", serrors.Format("%w", ErrNotFoundFile)("name", name)
}

var ErrNotFoundVersionThatSatisfiesConstraints = errors.New("not found version that satisfies constraints")

func latestVersionThatSatisfiesConstraints(
	command *gocommand.Memorized,
	path string,
	constraints *manifest.Constraints,
) (*semver.Version, error) {
	versions, err := command.MemorizedVersions(path)
	if err != nil {
		return nil, err
	}

	version, _, ok := lo.FindLastIndexOf(versions.Versions, constraints.Check)
	if !ok {
		return nil, serrors.Format("%w", ErrNotFoundVersionThatSatisfiesConstraints)(
			"path", path,
			"version", constraints,
		)
	}

	return version, nil
}
