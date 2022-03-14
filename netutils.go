package utils

import (
	"fmt"

	"github.com/fatih/color"
	externalip "github.com/glendc/go-external-ip"
)

// PublicIP uses several external services to get the public IP address of the
// system running it using https://pkg.go.dev/github.com/GlenDC/go-external-ip.
// It returns a public IP address or an error.
func PublicIP(protocol uint) (string, error) {

	// Create the default consensus with the default config
	// and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)

	if err := consensus.UseIPProtocol(protocol); err != nil {
		return "", fmt.Errorf(color.RedString("failed to get public IP address: %v", err))
	}

	// Retrieve the external IP address.
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", fmt.Errorf(color.RedString("failed to get public IP address: %v", err))
	}

	// Return the IP address in string format.
	return ip.String(), nil
}
