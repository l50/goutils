package chrome_test

import (
	"context"
	"os"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/l50/goutils/web"
	"github.com/l50/goutils/web/chrome"
)

func TestGetContext(t *testing.T) {
	tests := []struct {
		name   string
		driver chrome.Driver
		want   context.Context
	}{
		{
			name:   "Get existing context",
			driver: chrome.Driver{Context: context.Background()},
			want:   context.Background(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.driver.GetContext(); got != tc.want {
				t.Errorf("GetContext() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSetContext(t *testing.T) {
	tests := []struct {
		name     string
		driver   chrome.Driver
		newCtx   context.Context
		expected context.Context
	}{
		{
			name:     "Set new context",
			driver:   chrome.Driver{Context: context.Background()},
			newCtx:   context.TODO(),
			expected: context.TODO(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.driver.SetContext(tc.newCtx)
			if got := tc.driver.GetContext(); got != tc.expected {
				t.Errorf("SetContext() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name             string
		headless         bool
		ignoreCertErrors bool
	}{
		{
			name:             "Initialize browser",
			headless:         true,
			ignoreCertErrors: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			browser, err := chrome.Init(tc.headless, tc.ignoreCertErrors)
			if err != nil {
				t.Errorf("failed to initialize chrome: %v", err)
			}
			if browser.Driver == nil {
				t.Error("browser driver is nil")
			}
		})
	}
}

func TestNavigate(t *testing.T) {
	tests := []struct {
		name       string
		headless   bool
		ignoreCert bool
		url        string
	}{
		{
			name:       "Navigate with headless and ignoreCert",
			headless:   true,
			ignoreCert: true,
			url:        "https://google.com",
		},
		{
			name:       "Navigate with headless and not ignoreCert",
			headless:   true,
			ignoreCert: false,
			url:        "https://google.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the chrome browser
			browser, err := chrome.Init(tc.headless, tc.ignoreCert)
			if err != nil {
				t.Fatalf("failed to initialize a chrome browser: %v", err)
			}
			defer web.CancelAll(browser.Cancels...)

			site := web.Site{
				LoginURL: tc.url,
				Session: web.Session{
					Driver: browser.Driver,
				},
			}

			initialLoginActions := []chrome.InputAction{
				{
					Description: "Navigate to the login page",
					Action:      chromedp.Navigate(tc.url),
				},
			}

			waitTime, err := web.GetRandomWait(2, 6)
			if err != nil {
				t.Errorf("failed to create random wait time: %v", err)
			}

			if err := chrome.Navigate(site, initialLoginActions, waitTime); err != nil {
				t.Errorf("failed to navigate to %s: %v", site.LoginURL, err)
			}
		})
	}
}

func TestScreenShot(t *testing.T) {
	tests := []struct {
		name       string
		headless   bool
		ignoreCert bool
		url        string
		filename   string
	}{
		{
			name:       "Take screenshot",
			filename:   "test.png",
			headless:   true,
			ignoreCert: true,
			url:        "https://google.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the chrome browser
			browser, err := chrome.Init(true, true)
			if err != nil {
				t.Fatalf("failed to initialize a chrome browser: %v", err)
			}
			defer web.CancelAll(browser.Cancels...)

			// Set up the site with the browser's driver
			site := web.Site{
				LoginURL: tc.url,
				Session: web.Session{
					Driver: browser.Driver,
				},
			}

			// Navigation actions to set up for the screenshot
			initialLoginActions := []chrome.InputAction{
				{
					Description: "Navigate to the login page",
					Action:      chromedp.Navigate(tc.url),
				},
			}

			waitTime, err := web.GetRandomWait(2, 6)
			if err != nil {
				t.Errorf("failed to create random wait time: %v", err)
			}

			if err := chrome.Navigate(site, initialLoginActions, waitTime); err != nil {
				t.Errorf("failed to navigate to %s: %v", site.LoginURL, err)
			}

			if err := chrome.ScreenShot(site, tc.filename); err != nil {
				t.Errorf("failed to take screenshot: %v", err)
			} else {
				// Check if the file was created
				if _, err := os.Stat(tc.filename); os.IsNotExist(err) {
					t.Errorf("screenshot file was not created")
				} else {
					// Clean up the screenshot file after test
					if err := os.Remove(tc.filename); err != nil {
						t.Errorf("failed to delete screenshot file: %v", err)
					}
				}
			}
		})
	}
}

func TestGetPageSource(t *testing.T) {
	tests := []struct {
		name       string
		headless   bool
		ignoreCert bool
		url        string
	}{
		{
			name:       "Get page source",
			headless:   true,
			ignoreCert: true,
			url:        "https://google.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the chrome browser
			browser, err := chrome.Init(tc.headless, tc.ignoreCert)
			if err != nil {
				t.Fatalf("failed to initialize a chrome browser: %v", err)
			}
			defer web.CancelAll(browser.Cancels...)

			// Set up the site with the browser's driver
			site := web.Site{
				LoginURL: tc.url,
				Session: web.Session{
					Driver: browser.Driver,
				},
			}

			// Navigation actions to set up for the test
			initialLoginActions := []chrome.InputAction{
				{
					Description: "Navigate to the page",
					Action:      chromedp.Navigate(tc.url),
				},
			}

			waitTime, err := web.GetRandomWait(2, 6)
			if err != nil {
				t.Errorf("failed to create random wait time: %v", err)
			}

			if err := chrome.Navigate(site, initialLoginActions, waitTime); err != nil {
				t.Errorf("failed to navigate to %s: %v", site.LoginURL, err)
			}

			source, err := chrome.GetPageSource(site)
			if err != nil {
				t.Errorf("failed to get page source: %v", err)
			}

			// Check if page source is not empty
			if len(source) == 0 {
				t.Errorf("page source is empty")
			}
		})
	}
}
