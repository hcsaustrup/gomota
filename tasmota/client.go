package tasmota

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type TasmotaClient struct {
	Hostname string
	Username string
	Password string
	_logger  *logrus.Entry
}

func (c *TasmotaClient) logger() *logrus.Entry {
	if c._logger == nil {
		c._logger = logrus.WithField("hostname", c.Hostname)
	}
	return c._logger
}

func (c *TasmotaClient) RunCommand(command string, response interface{}) error {
	client := &http.Client{
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	req.Header.Add("Authorization", network.BasicAuthHeader(opts.username, opts.password))
		// 	return nil
		// },
	}

	params := url.Values{}
	params.Add("cmnd", command)
	params.Add("user", c.Username)
	params.Add("password", c.Password)

	url := fmt.Sprintf("http://%s/cm?%s", c.Hostname, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// req.Header.Add("Authorization", network.BasicAuthHeader(c.Username, c.Password))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body")
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to decode status response: %v", err)
	}

	logrus.WithField("url", url).Debugf("Request completed")
	return nil
}
