package cloudflare_test

import (
	"fmt"
	"io"
	"os"

	"github.com/l50/goutils/v2/cloudflare"
)

func ExampleGetDNSRecords() {
	// Mocked HTTP client
	mockClient := &MockHTTPClient{}

	cf := cloudflare.Cloudflare{
		CFApiKey: "your-api-key",
		CFEmail:  "your-email@example.com",
		CFZoneID: "your-zone-id",
		Client:   mockClient,
	}

	// We have to replace the standard output temporarily to capture it.
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Stdout = w

	if err := cloudflare.GetDNSRecords(cf); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Close the Pipe Writer to let the ReadString finish properly.
	w.Close()

	// Restore the standard output.
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Display the result we've captured.
	fmt.Print(string(out))

	// Output:
	// name:  your-dns-name
	// content:  your-dns-content
}
