package web

import (
	"context"
	"math/rand"
	"time"

	"github.com/l50/goutils/v2/web/chrome"
)

// Driver is an interface that can be implemented for
// various browsers in go. It should service to capture
// information such as context and browser options.
type Driver interface {
	GetContext() context.Context
}

// Browser defines parameters used
// for driving a web browser.
type Browser struct {
	Cancels []func()
	Driver  *chrome.Driver
}

// Site is used to define parameters
// for TTPs targeting web applications
// that require browser automation.
type Site struct {
	LoginURL string
	Session  *Session
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
func Wait(near float64) time.Duration {
	zoom := int(near / 10)
	x := rand.Intn(zoom) + int(0.95*near)
	return time.Duration(x) * time.Millisecond
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
