package net

import (
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	externalip "github.com/glendc/go-external-ip"
)

// PublicIP uses several external services to get the public IP address of the system running it, using the
// package github.com/GlenDC/go-external-ip. The function takes an IP protocol version (4 or 6) as input and
// returns the public IP address as a string or an error.
//
// Parameters:
//
// protocol: A uint representing the IP protocol version (4 or 6).
//
// Returns:
//
// string: The public IP address of the system in string format.
// error: An error if the function fails to retrieve the public IP address.
//
// Example:
//
// protocol := uint(4) // or 6 for IPv6
// ip, err := PublicIP(protocol)
//
//	if err != nil {
//	  log.Fatalf("failed to get public IP address: %v", err)
//	}
//
// log.Printf("Public IP address: %s\n", ip)
func PublicIP(protocol uint) (string, error) {

	// Create the default consensus with the default config
	// and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)

	if err := consensus.UseIPProtocol(protocol); err != nil {
		return "", fmt.Errorf("failed to get public IP address: %v", err)
	}

	// Retrieve the external IP address.
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", fmt.Errorf("failed to get public IP address: %v", err)
	}

	// Return the IP address in string format.
	return ip.String(), nil
}

// DownloadFile downloads a file from the provided URL and saves it to the specified location on the local
// filesystem. The function takes the source URL and the destination path as inputs and returns the path
// where the file was saved or an error.
//
// Parameters:
//
// url: A string representing the URL of the file to be downloaded.
// dest: A string representing the destination path where the file should be saved on the local filesystem.
//
// Returns:
//
// string: The path where the downloaded file was saved.
// error: An error if the function fails to download the file.
//
// Example:
//
// url := "http://example.com/path/to/file"
// dest := "/path/to/save/location"
// file, err := DownloadFile(url, dest)
//
//	if err != nil {
//	  log.Fatalf("failed to download file: %v", err)
//	}
//
// log.Printf("File downloaded to: %s\n", file)
func DownloadFile(url string, dest string) (string, error) {
	resp, err := grab.Get(dest, url)
	if err != nil {
		return resp.Filename, fmt.Errorf("failed to download %s to %s: %v", url, dest, err)
	}

	return resp.Filename, nil
}
