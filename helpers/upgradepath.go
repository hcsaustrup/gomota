package helpers

import (
	"github.com/hashicorp/go-version"
)

type UpgradePath struct {
	path []*version.Version
}

func (up *UpgradePath) Add(rawVersion string) error {
	version, err := version.NewVersion(rawVersion)
	if err != nil {
		return err
	}
	up.path = append(up.path, version)
	return nil
}

func (up *UpgradePath) FindNextByString(rawVersion string) (*version.Version, error) {
	version, err := version.NewVersion(rawVersion)
	if err != nil {
		return nil, nil
	}
	return up.FindNext(version), nil
}

func (up *UpgradePath) FindNext(version *version.Version) *version.Version {
	for _, pathVersion := range up.path {
		if version.LessThan(pathVersion) {
			return pathVersion
		}
	}
	return nil
}
