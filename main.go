package main

import (
	"errors"
	"net"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"saustrup.net/gomota/helpers"
	"saustrup.net/gomota/network"
	"saustrup.net/gomota/tasmota"
)

type options struct {
	network     string
	username    string
	password    string
	upgradePath string
	debug       bool
}

func main() {

	opts := &options{}

	//---------------------------------------------------------------------------------------------------------------------
	// Parse options
	//---------------------------------------------------------------------------------------------------------------------

	flag.StringVar(&opts.network, "network", "10.69.1.0/24", "Network to scan in network/prefix notation")
	flag.StringVar(&opts.upgradePath, "upgrade-path", "1.0.11,3.9.22,4.2.0,5.14.0,6.7.1,7.2.0,8.5.1,9.1.0,10.1.0,11.1.0,13.1.0.1", "Firmware upgrade path")
	flag.StringVar(&opts.username, "username", "", "Tasmota username")
	flag.StringVar(&opts.password, "password", "", "Tasmota password")
	flag.BoolVar(&opts.debug, "debug", false, "Enable debugging")
	flag.Parse()

	if opts.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	upgradePath := &helpers.UpgradePath{}
	for _, versionString := range strings.Split(opts.upgradePath, ",") {
		if err := upgradePath.Add(versionString); err != nil {
			logrus.WithError(err).Fatalf("Failed to parse version in upgrade path: %s", versionString)
		}
	}

	//---------------------------------------------------------------------------------------------------------------------
	// Instantiate services
	//---------------------------------------------------------------------------------------------------------------------

	hosts, err := network.GetHostsFromCIDR(opts.network)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, host := range hosts {
		wg.Add(1)
		go (func(hostname string) {
			defer wg.Done()

			c := &tasmota.TasmotaClient{
				Hostname: hostname,
				Username: opts.username,
				Password: opts.password,
			}

			logger := logrus.WithField("hostname", hostname)

			if err := c.Upgrader(upgradePath); err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// A timeout error occurred
				} else if errors.Is(err, syscall.ECONNREFUSED) {
					logger.Warn("Not a Tasmota device")
				} else if errors.Is(err, syscall.ETIMEDOUT) {
					//
				} else if errors.Is(err, syscall.EHOSTUNREACH) {
					//
				} else {
					logger.WithError(err).Warn("Failed")
				}
			}

		})(host)
	}

	logrus.Info("Waiting for workers to finish.")
	wg.Wait()

	logrus.Info("All services stopped. See you next Wednesday.")
}
