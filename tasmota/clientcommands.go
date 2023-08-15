package tasmota

import (
	"fmt"

	"saustrup.net/gomota/tasmota/models"
)

func (c *TasmotaClient) GetDetailedStatus() (*models.StatusResponse, error) {
	response := &models.StatusResponse{}
	err := c.RunCommand("status 0", &response)
	return response, err
}

func (c *TasmotaClient) GetFirmwareStatus() (*models.FirmwareStatusResponse, error) {
	response := &models.FirmwareStatusResponse{}
	err := c.RunCommand("status 2", &response)
	return response, err
}

func (c *TasmotaClient) SetUpgradeURL(url string) (*models.OTAURLResponse, error) {
	response := &models.OTAURLResponse{}
	err := c.RunCommand(fmt.Sprintf("otaurl %s", url), &response)
	return response, err
}

func (c *TasmotaClient) StartUpgrade() (*models.UpgradeResponse, error) {
	response := &models.UpgradeResponse{}
	err := c.RunCommand("upgrade 1", &response)
	return response, err
}
