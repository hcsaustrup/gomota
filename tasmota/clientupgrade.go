package tasmota

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-version"
	"saustrup.net/gomota/helpers"
)

const (
	defaultVariant      = "tasmota"
	intermediateVariant = "lite"
)

func (c *TasmotaClient) Upgrader(upgradePath *helpers.UpgradePath) error {
	status, err := c.GetDetailedStatus()
	if err != nil {
		return err
	}

	// Add details to logger
	c._logger = c.logger().WithField("name", status.Status.DeviceName)

	versionRegexp := regexp.MustCompile(`(.*)\((\w+)\)`)
	versionMatches := versionRegexp.FindStringSubmatch(status.StatusFWR.Version)
	if len(versionMatches) == 0 {
		return fmt.Errorf("unexpected firmware pattern: %s", status.StatusFWR.Version)
	}
	currentVersion, err := version.NewVersion(versionMatches[1])
	if err != nil {
		return fmt.Errorf("failed to parse version: %s", versionMatches[1])
	}

	currentVariant := versionMatches[2]
	if currentVariant != defaultVariant && currentVariant != intermediateVariant {
		return fmt.Errorf("variant is neither default not intermediate - unsure how to handle this")
	}

	// Got upgrade path?
	nextVersion := upgradePath.FindNext(currentVersion)
	nextVariant := intermediateVariant
	if nextVersion == nil {
		if currentVariant == defaultVariant {
			c.logger().Infof("No upgrade path available for %s(%s)", currentVersion.String(), currentVariant)
			return nil
		}
		nextVersion = currentVersion
		nextVariant = defaultVariant
	}

	// Build OTA URL
	filename := "tasmota"
	if nextVariant != defaultVariant {
		filename += fmt.Sprintf("-%s", nextVariant)
	}
	filename += ".bin"
	if gzipConstraints, err := version.NewConstraint(">= 8.2"); err != nil {
		panic(err)
	} else if gzipConstraints.Check(currentVersion) {
		filename += ".gz"
	}
	otaUrl := fmt.Sprintf("http://ota.tasmota.com/tasmota/release-%s/%s", nextVersion.String(), filename)

	if otaurlResponse, err := c.SetUpgradeURL(otaUrl); err != nil {
		return err
	} else if otaurlResponse.OTAURL != otaUrl {
		return fmt.Errorf("failed to set OTA URL: %s", otaUrl)
	}

	if _, err := c.StartUpgrade(); err != nil {
		return err
	}

	c.logger().Infof("Upgrading: %s(%s) -> %s(%s)", currentVersion.String(), currentVariant, nextVersion.String(), nextVariant)

	return nil
}
