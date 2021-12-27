package utils

import (
	"testing"
)

func TestPublicIP(t *testing.T) {
	protocols := []uint{4, 6}
	for _, protocol := range protocols {
		ip, err := PublicIP(protocol)
		if err != nil {
			t.Fatalf("Unable to return public IPv%v address: %v, TestPublicIP() failed.", protocol, err)
		} else {
			t.Logf("IPv%v address retrieved successfully: %s\n", protocol, ip)
			return
		}
	}
}
