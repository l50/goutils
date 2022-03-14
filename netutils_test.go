package utils

import (
	"testing"
)

func TestPublicIP(t *testing.T) {
	protocols := []uint{4, 6}
	for _, protocol := range protocols {
		_, err := PublicIP(protocol)
		if err != nil {
			t.Fatalf("unable to return public IPv%v address - TestPublicIP() failed: %v", protocol, err)
		}
	}
}
