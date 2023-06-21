package cloudflare

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// HTTPClient defines the behavior of an HTTP client.
//
// **Attributes:**
//
// Do: Sends an HTTP request and returns the HTTP response or an error.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Cloudflare represents information needed to interface
// with the Cloudflare API.
//
// **Attributes:**
//
// CFApiKey: Cloudflare API key.
// CFEmail: Email associated with the Cloudflare account.
// CFZoneID: Zone ID of the domain on Cloudflare.
// Email: Email address for notifications.
// Endpoint: API endpoint for Cloudflare.
// Client: HTTP client for making requests.
type Cloudflare struct {
	CFApiKey string
	CFEmail  string
	CFZoneID string
	Email    string
	Endpoint string
	Client   HTTPClient
}

// GetDNSRecords retrieves the DNS records from Cloudflare for a
// specified zone ID using the provided Cloudflare credentials.
// It makes a GET request to the Cloudflare API, reads the
// response, and prints the 'name' and 'content' fields of
// each DNS record.
//
// **Parameters:**
//
// cf: A Cloudflare struct containing the necessary credentials
// (email, API key) and the zone ID for which the DNS records
// should be retrieved.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to
// get the DNS records.
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
