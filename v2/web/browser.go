package web

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
)

// Driver is an interface that can be implemented for
// various browsers in go. It should service to capture
// information such as context and browser options.
type Driver interface {
	GetContext() context.Context
	SetContext(context.Context)
	GetOptions() []chromedp.ExecAllocatorOption
	SetOptions([]chromedp.ExecAllocatorOption)
}

// Browser defines parameters used
// for driving a web browser.
type Browser struct {
	Driver  interface{}
	Cancels []func()
}

// Site is used to define parameters
// for interacting with web applications.
type Site struct {
	LoginURL string
	Session  Session
	Debug    bool
}

// FormField contains a form field name and its associated selector.
type FormField struct {
	Name     string `json:"-"`
	Selector string `json:"-"`
}

// CancelAll cancels all input cancels.
func CancelAll(cancels []func()) {
	for _, cancel := range cancels {
		cancel()
	}
}

// cryptoRandIntn generates a random int in the range [0, n) using crypto/rand.
func cryptoRandIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid argument to cryptoRandIntn: %d <= 0", n)
	}

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return 0, err
	}

	val := binary.BigEndian.Uint64(b) // Convert bytes to uint64.

	// Return a value in the range [0, n).
	return int(val % uint64(n)), nil
}

// Wait is used to wait for a random period of time
// that is anchored to the input near value.
//
// It's useful to simulate a more human interaction
// while interfacing with a web application.
//
// Example:
// err = chromedp.Run(caldera.Driver.Context,
// chromedp.Navigate(caldera.URL),
// chromedp.Sleep(Wait(1000)),
// chromedp.SendKeys(userXPath, caldera.Creds.User),
// chromedp.Sleep(Wait(1000)),
// ...
func Wait(near float64) (time.Duration, error) {
	zoom := int(near / 10)
	x, err := cryptoRandIntn(zoom)
	if err != nil {
		return 0, err
	}
	x += int(0.95 * near)
	return time.Duration(x) * time.Millisecond, nil
}

// GetRandomWait returns a random duration between the specified minWait and maxWait durations.
// The function takes the minimum and maximum wait times as arguments, creates a new random
// number generator with a seed based on the current Unix timestamp, and calculates the random
// wait time within the given range.
//
// Example usage:
//
//	minWait := 2 * time.Second
//	maxWait := 6 * time.Second
//	randomWaitTime := GetRandomWait(minWait, maxWait)
//
// Parameters:
//
//	minWait: The minimum duration to wait.
//	maxWait: The maximum duration to wait.
//
// Returns:
//
//	time.Duration: A random duration between minWait and maxWait.
func GetRandomWait(minWait, maxWait time.Duration) time.Duration {
	// Create rng based on current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWaitTime := time.Duration(rng.Int63n(int64(maxWait-minWait))) + minWait
	return randomWaitTime
}
