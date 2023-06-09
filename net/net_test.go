package net_test

import (
	"testing"

	"github.com/l50/goutils/v2/net"
)

func TestDownloadFile(t *testing.T) {
	testCases := []struct {
		name    string
		url     string
		dest    string
		want    string
		wantErr bool
	}{
		{
			name:    "Successfully download file",
			url:     "https://raw.githubusercontent.com/l50/helloworld/master/hello.go",
			dest:    "/tmp/hello.go",
			want:    "/tmp/hello.go",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := net.DownloadFile(tc.url, tc.dest)
			if (err != nil) != tc.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("DownloadFile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPublicIP(t *testing.T) {
	testCases := []struct {
		name     string
		protocol uint
		wantErr  bool
	}{
		{
			name:     "IPv4 public IP",
			protocol: 4,
			wantErr:  false,
		},
		{
			name:     "IPv6 public IP",
			protocol: 6,
			wantErr:  true, // Considering that a lot of networks aren't using IPv6
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := net.PublicIP(tc.protocol)
			if (err != nil) != tc.wantErr {
				t.Errorf("PublicIP() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
