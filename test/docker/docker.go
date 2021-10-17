package ymirtestdocker

import (
	"errors"
	"fmt"
	"os"

	dc "github.com/ory/dockertest/v3/docker"
)

func FindTestNetwork(c *dc.Client) (*dc.Network, error) {
	nwName, found := os.LookupEnv("YMIR_CI_DOCKER_NETWORK")

	if !found {
		return nil, errors.New("YMIR_CI_DOCKER_NETWORK is not defined")
	}

	nws, err := c.ListNetworks()

	if err != nil {
		return nil, err
	}

	for _, nw := range nws {
		if nw.Name == nwName {
			return &nw, nil
		}
	}

	return nil, fmt.Errorf("YMIR_CI_DOCKER_NETWORK (%s) not found", nwName)
}