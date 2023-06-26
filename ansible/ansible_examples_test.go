package ansible_test

import (
	"fmt"
	"io"
	"os"

	"github.com/l50/goutils/v2/ansible"
)

func ExamplePing() {
	hostsFile := "/path/to/your/hosts.ini"

	// We have to replace the standard output temporarily to capture it.
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Stdout = w

	if err := ansible.Ping(hostsFile); err != nil {
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
	// localhost | SUCCESS => {
	//     "changed": false,
	//     "ping": "pong"
	// }
}
