package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/fatih/color"
)

// KeeperLoggedIn returns true if keeper vault
// is logged in with the input `email`.
// Otherwise, it returns false.
func KeeperLoggedIn(email string) bool {
	fmt.Println(color.YellowString(
		"Checking if we are logged into Keeper vault"))
	loggedIn := "My Vault>"

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Get the keeper menu output and exit
	// Semgrep is falsely flagging this as a SQLi
	// nosemgrep
	_, err := script.Echo("q").Exec("keeper login " + email).Stdout()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if err != nil {
		fmt.Printf("failed to check login state "+
			"for keeper vault: %v", err)
		return false
	}

	// The output response has a ton of newlines -
	// split on newlines to make it easier to
	// get the auth URL.
	outSlice := strings.Split(string(out), "\n")

	for _, output := range outSlice {
		if strings.Contains(output, loggedIn) {
			return true
		}
	}

	return false

}

// RetrieveKeeperPW returns the password found at
// the specified input path.
func RetrieveKeeperPW(path string) (string, error) {
	fmt.Println(color.YellowString(
		"Retrieving %s from keeper", path))

	cmd := fmt.Sprintf("keeper find-password '%s'", path)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Get password
	_, err := script.Exec(cmd).Stdout()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if err != nil {
		return "", err
	}

	// Remove newlines from output
	return strings.ReplaceAll(string(out), "\n", ""), nil
}
