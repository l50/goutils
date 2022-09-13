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
			t.Errorf("unable to return public IPv%v address - TestPublicIP() failed: %v", protocol, err)
		}
	}
}

func TestDownloadFile(t *testing.T) {
	type args struct {
		url  string
		dest string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Successfully download file",
			args: args{
				url:  "https://raw.githubusercontent.com/l50/helloworld/master/hello.go",
				dest: "/tmp/hello.go",
			},
			want:    "/tmp/hello.go",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadFile(tt.args.url, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("error: DownloadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
