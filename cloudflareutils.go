package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// Cloudflare holds information needed to interface
// with the cloudflare API.
type Cloudflare struct {
	CFApiKey string
	CFEmail  string
	CFZoneID string
	Email    string
	Endpoint string
	Client   http.Client
}

// GetDNSRecords retrieves the DNS records from cloudflare.
func GetDNSRecords(cf Cloudflare) error {
	cf.Client = http.Client{}
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
