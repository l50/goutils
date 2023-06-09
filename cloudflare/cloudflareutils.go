package cloudflare

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// HTTPClient is an interface defining the behavior of an HTTP client. This allows for more flexibility
// in HTTP interactions and facilitates testing by enabling the use of mock HTTP clients.
//
// The HttpClient interface includes a single method, Do, which sends an HTTP request and returns
// the HTTP response or an error.
//
// Parameters:
//
// req: The HTTP request to be sent.
//
// Returns:
//
// *http.Response: The HTTP response to the request.
// error: An error, if one occurred during the execution of the request.
//
// Example:
//
//	type MockClient struct {
//	    MockDo func(req *http.Request) (*http.Response, error)
//	}
//
//	func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
//	    return m.MockDo(req)
//	}
//
//	mockClient := &MockClient{
//	    MockDo: func(req *http.Request) (*http.Response, error) {
//	        return &http.Response{
//	            StatusCode: 200,
//	            Body:       io.NopCloser(strings.NewReader(`{"result": "ok"}`)),
//	        }, nil
//	    },
//	}
//
//	cf := Cloudflare{
//	    CFApiKey: "valid_key",
//	    CFEmail:  "valid_email@example.com",
//	    CFZoneID: "valid_zone",
//	    Email:    "notification@example.com",
//	    Client:   mockClient,
//	}
//
// err := GetDNSRecords(cf)
//
//	if err != nil {
//	    log.Fatalf("failed to get DNS records: %v", err)
//	}
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Cloudflare holds information needed to interface
// with the Cloudflare API.
//
// Parameters:
//
// CFApiKey: Cloudflare API key.
// CFEmail: Email address associated with the Cloudflare account.
// CFZoneID: Zone ID of the domain on Cloudflare.
// Email: Email address for notifications.
// Endpoint: API endpoint for Cloudflare.
// Client: HTTP client for making requests.
//
// Example:
//
//	cf := Cloudflare{
//	  CFApiKey: "your_api_key",
//	  CFEmail: "your_email@example.com",
//	  CFZoneID: "your_zone_id",
//	  Email: "your_notification_email@example.com",
//	  Endpoint: "",  // This will be set in the function.
//	  Client: http.Client{},
//	}
//
// err := GetDNSRecords(cf)
//
//	if err != nil {
//	  log.Fatalf("failed to get DNS records: %v", err)
//	}
type Cloudflare struct {
	CFApiKey string
	CFEmail  string
	CFZoneID string
	Email    string
	Endpoint string
	Client   HTTPClient
}

// GetDNSRecords retrieves the DNS records from Cloudflare for a specified zone ID using the provided Cloudflare credentials.
// It makes a GET request to the Cloudflare API, reads the response, and prints the 'name' and 'content' fields of each DNS record.
//
// Parameters:
//
// cf: A Cloudflare struct containing the necessary credentials (email, API key) and the zone ID for which the DNS records should be retrieved.
//
// Returns:
//
// error: An error if any issue occurs while trying to get the DNS records.
//
// Example:
//
//	cf := Cloudflare{
//	  CFEmail: "your-email@example.com",
//	  CFApiKey: "your-api-key",
//	  CFZoneID: "your-zone-id",
//	  Client: &http.Client{},
//	}
//
// err := GetDNSRecords(cf)
//
//	if err != nil {
//	  log.Fatalf("failed to get DNS records: %v", err)
//	}
func GetDNSRecords(cf Cloudflare) error {
	cf.Endpoint = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", cf.CFZoneID)

	req, err := http.NewRequest("GET", cf.Endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Email", cf.CFEmail)
	req.Header.Set("X-Auth-Key", cf.CFApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cf.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := gjson.GetBytes(body, "result")
	result.ForEach(func(key, value gjson.Result) bool {
		if !(value.Get("name").String()[0] == '*') {
			fmt.Println("name: ", value.Get("name").String())
			fmt.Println("content: ", value.Get("content").String())
		}
		return true
	})

	return nil
}
