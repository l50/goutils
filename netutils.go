package utils

import (
	"fmt"
	"os"

	externalip "github.com/glendc/go-external-ip"
)

func errCheck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// PublicIP uses several external services to get the public IP address of the
// system running it using https://pkg.go.dev/github.com/GlenDC/go-external-ip.
// It returns a public IP address or an error.
func PublicIP(protocol uint) (string, error) {

	// Create the default consensus with the default config
	// and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)
	err := consensus.UseIPProtocol(protocol)
	if err != nil {
		return "", err
	}

	// Retrieve the external IP address.
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", err
	}

	// Return the IP address in string format.
	return ip.String(), nil
}
