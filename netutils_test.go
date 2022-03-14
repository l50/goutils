package utils

import (
	"testing"
)

func TestPublicIP(t *testing.T) {
	protocols := []uint{4, 6}
	for _, protocol := range protocols {
		_, err := PublicIP(protocol)
		// A lot of networks aren't using IPv6 - let's avoid false positives
		if err != nil && protocol == 4 {
			t.Fatalf("unable to return public IPv%v address - TestPublicIP() failed: %v", protocol, err)
		}
	}
}
