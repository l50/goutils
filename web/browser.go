package web

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
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

// loginOptions is a struct that holds the configurations for the login process.
//
// twoFacEnabled: Determines if two-factor authentication is enabled during login.
// logMeOut: Determines if the user is logged out after login.
type loginOptions struct {
	twoFacEnabled bool
	logMeOut      bool
}

// LoginOption is a type for functions that modify the login options.
// These functions take a pointer to a loginOptions struct and modify it in place.
type LoginOption func(*loginOptions)

// CancelAll executes all provided cancel functions. It is typically used for cleaning up or aborting operations
// that were started earlier and can be cancelled.
//
// Parameters:
//
// cancels: A slice of cancel functions, each of type func(). These are typically functions returned by context.WithCancel,
// or similar functions that provide a way to cancel an operation.
//
// Example:
//
// var cancels []func()
//
// ctx, cancel := context.WithCancel(context.Background())
// cancels = append(cancels, cancel)
//
// // Later, when all operations need to be cancelled:
// CancelAll(cancels)
//
// Note: The caller is responsible for handling any errors that may occur during the execution of the cancel functions.
func CancelAll(cancels ...func()) {
	for _, cancel := range cancels {
		cancel()
	}
}

// cryptoRandIntn generates a random int64 in the range [0, n) using crypto/rand.
func cryptoRandIntn(n int64) (int64, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid argument to cryptoRandIntn: %d <= 0", n)
	}

	val, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0, err
	}

	return val.Int64(), nil
}

// GetRandomWait returns a random duration in seconds between the specified minWait
// and maxWait durations. The function takes the minimum and maximum wait times as
// arguments, creates a new random number generator with a seed based on the current
// Unix timestamp, and calculates the random wait time within the given range.
//
// Parameters:
//
//	minWait: The minimum duration to wait.
//	maxWait: The maximum duration to wait.
//
// Returns:
//
//	time.Duration: A random duration between minWait and maxWait.
//	error: An error if the generation of the random wait time fails.
//
// Example:
//
//	minWait := 2
//	maxWait := 6
//	randomWaitTime, err := GetRandomWait(minWait, maxWait)
//
//	if err != nil {
//	  log.Fatalf("failed to generate random wait time: %v", err)
//	}
func GetRandomWait(min, max int) (time.Duration, error) {
	minWait := time.Duration(min) * time.Second
	maxWait := time.Duration(max) * time.Second

	if minWait < 0 {
		return 0, fmt.Errorf("minWait cannot be less than 0: %v", minWait)
	}
	if maxWait < minWait {
		return 0, fmt.Errorf("maxWait cannot be less than minWait: %v < %v", maxWait, minWait)
	}

	diff := maxWait - minWait
	randomValue, err := cryptoRandIntn(int64(diff))
	if err != nil {
		return 0, err
	}
	randomWaitTime := time.Duration(randomValue) + minWait
	return randomWaitTime, nil
}

// IsTwoFacEnabled checks if two-factor authentication is enabled in the provided login options.
//
// Parameters:
//
// opts: A pointer to a loginOptions instance.
//
// Returns:
//
// bool: A boolean indicating whether two-factor authentication is enabled.
//
// Example:
//
//	options := []web.LoginOption{
//	   web.WithTwoFacEnabled(false),
//	}
//
// loginOpts := web.SetLoginOptions(options...)
// isTwoFacEnabled := web.IsTwoFacEnabled(loginOpts) // isTwoFacEnabled is now false.
func IsTwoFacEnabled(opts *loginOptions) bool {
	return opts.twoFacEnabled
}

// IsLogMeOutEnabled checks if the option to log out the user after login is enabled in the provided login options.
//
// Parameters:
//
// opts: A pointer to a loginOptions instance.
//
// Returns:
//
// bool: A boolean indicating whether the user is to be logged out after login.
//
// Example:
//
//	options := []web.LoginOption{
//	   web.WithLogMeOut(false),
//	}
//
// loginOpts := web.SetLoginOptions(options...)
// isLogMeOutEnabled := web.IsLogMeOutEnabled(loginOpts) // isLogMeOutEnabled is now false.
func IsLogMeOutEnabled(opts *loginOptions) bool {
	return opts.logMeOut
}

// SetLoginOptions applies provided login options to a new loginOptions instance with default values
// and returns a pointer to this instance. This function is primarily used to configure login behavior
// in the LoginAccount function.
//
// Parameters:
//
// options: A variadic set of LoginOption functions. Each LoginOption is a function that takes a
// pointer to a loginOptions struct and modifies it in place.
//
// Returns:
//
// *loginOptions: A pointer to a loginOptions struct that has been configured with the provided options.
//
// Example:
//
//	options := []web.LoginOption{
//	   web.WithTwoFacEnabled(false),
//	   web.WithLogMeOut(true),
//	}
//
// loginOpts := web.SetLoginOptions(options...)
//
// // Now loginOpts.twoFacEnabled is false and loginOpts.logMeOut is true.
func SetLoginOptions(options ...LoginOption) *loginOptions {
	loginOpts := &loginOptions{
		twoFacEnabled: true, // Default value for twoFac
		logMeOut:      true, // Default value for logout
	}

	for _, option := range options {
		option(loginOpts)
	}

	return loginOpts
}

// Wait generates a random period of time anchored to a given input value.
//
// Parameters:
//
// near: A float64 value that serves as the base value for generating the random wait time.
//
// Returns:
//
// time.Duration: The calculated random wait time in milliseconds.
// error: An error if the generation of the random wait time fails.
//
// The function is useful for simulating more human-like interaction in the context of a web application.
// It first calculates a 'zoom' value by dividing the input 'near' by 10. Then, a random number is generated in the
// range of [0, zoom), and added to 95% of 'near'. This sum is then converted to a time duration in milliseconds and returned.
//
// Example:
//
// waitTime, err := Wait(1000.0)
//
//	if err != nil {
//	  log.Fatalf("failed to generate random wait time: %v", err)
//	}
//
// log.Printf("Generated random wait time: %v\n", waitTime)
func Wait(near float64) (time.Duration, error) {
	zoom := int64(near / 10)
	x, err := cryptoRandIntn(zoom)
	if err != nil {
		return 0, err
	}
	additionalWait := int64(0.95 * near)
	totalWait := x + additionalWait
	return time.Duration(totalWait) * time.Millisecond, nil
}

// WithLogout is a function that returns a LoginOption function which sets the logMeOut option.
//
// Parameters:
//
// enabled: Determines if the user should be logged out after login.
//
// Returns:
//
// LoginOption: A function that modifies the logMeOut option of a loginOptions struct.
func WithLogout(enabled bool) LoginOption {
	return func(opts *loginOptions) {
		opts.logMeOut = enabled
	}
}

// WithTwoFac is a function that returns a LoginOption function which sets the twoFacEnabled option.
//
// Parameters:
//
// enabled: Determines if two-factor authentication should be enabled during login.
//
// Returns:
//
// LoginOption: A function that modifies the twoFacEnabled option of a loginOptions struct.
func WithTwoFac(enabled bool) LoginOption {
	return func(opts *loginOptions) {
		opts.twoFacEnabled = enabled
	}
}
