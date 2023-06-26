package cloudflare_test

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/l50/goutils/v2/cloudflare"
	"github.com/stretchr/testify/assert"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestGetDNSRecords(t *testing.T) {
	tests := []struct {
		name         string
		cf           cloudflare.Cloudflare
		responseBody string
		expectedErr  error
	}{
		{
			name: "valid case",
			cf: cloudflare.Cloudflare{
				CFApiKey: "valid_key",
				CFEmail:  "valid_email@example.com",
				CFZoneID: "valid_zone",
				Email:    "notification@example.com",
				Client: &MockClient{
					MockDo: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(strings.NewReader(`{"result": [{"name": "example.com", "content": "192.0.2.0"}]}`)),
						}, nil
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := cloudflare.GetDNSRecords(tc.cf)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Here you have to provide a response in the format that your function expects
	jsonResp := `{
		"result": [
			{
				"name": "your-dns-name",
				"content": "your-dns-content"
			}
		]
	}`

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(jsonResp)),
		Header:     make(http.Header),
	}, nil
}
