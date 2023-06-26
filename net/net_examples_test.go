package net_test

import (
	"log"
	"net"

	netutils "github.com/l50/goutils/v2/net"
)

func ExampleDownloadFile() {
	url := "http://example.com/path/to/file"
	dest := "/path/to/save/location"
	file, err := netutils.DownloadFile(url, dest)

	if err != nil {
		log.Fatalf("failed to download file: %v", err)
	}
	_ = file
}

func ExamplePublicIP() {
	protocol := uint(4) // or 6 for IPv6
	ip, err := netutils.PublicIP(protocol)
	if err != nil {
		log.Fatalf("failed to get public IP address: %v", err)
	}

	if net.ParseIP(ip) == nil {
		log.Fatal("invalid IP address received")
	}

	_ = ip
}
