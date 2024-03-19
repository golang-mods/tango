package tango

import "github.com/Masterminds/semver/v3"

func (manager *Manager) InfoVersions(pkg string) ([]*semver.Version, error) {
	versions, err := command.MemorizedVersions(pkg)
	if err != nil {
		return nil, err
	}

	return versions.Versions, nil
}
