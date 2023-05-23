package chrome

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/l50/goutils/v2/web"
)

// Driver is used to interface with Google Chrome using go.
type Driver struct {
	Context context.Context
	Options *[]chromedp.ExecAllocatorOption
}

// GetContext returns the context associated with the Driver.
func (d *Driver) GetContext() context.Context {
	return d.Context
}

// setChromeOptions is used to set the chrome
// parameters required by ChromeDP.
func setChromeOptions(browser *web.Browser, headless bool, ignoreCertErrors bool, options *[]chromedp.ExecAllocatorOption) {
	*options = append(*options,
		chromedp.DisableGPU,
		chromedp.Flag("ignoreCertErrors", ignoreCertErrors),
		// Uncomment to prevent navigation to a new tab
		// chromedp.Flag("block-new-web-contents", true),
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
		chromedp.Flag("headless", headless),
	)
	browser.Driver = &Driver{
		Options: options,
	}
}

// Init returns a chrome browser for the TTP Runner to use.
func Init(headless bool, ignoreCertErrors bool) (web.Browser, error) {
	var cancels []func()

	browser := web.Browser{
		Cancels: cancels,
	}

	options := []chromedp.ExecAllocatorOption{}
	setChromeOptions(&browser, headless, ignoreCertErrors, &options)

	driver, ok := browser.Driver.(*Driver)
	if !ok {
		if err := errors.New("driver is not of type *ChromeDP"); err != nil {
			return web.Browser{}, err
		}
	}

	// Create contexts and their associated cancels.
	allocatorCtx, cancel := chromedp.NewExecAllocator(
		context.Background(), *driver.Options...)
	browser.Cancels = append([]func(){cancel}, cancels...)
	driver.Context, cancel = chromedp.NewContext(allocatorCtx,
		chromedp.WithLogf(log.Printf))
	browser.Cancels = append([]func(){cancel}, cancels...)

	return browser, nil
}

// InputAction contains selectors and actions to run
// with chrome.
type InputAction struct {
	Description string
	Selector    string
	Action      chromedp.Action
}

// GetPageSource retrieves the source code of the web page currently loaded in the site session.
// It returns the page source code as a string or an error if the function fails to obtain the source code.
//
// Args:
//
// site (web.Site): a web.Site instance containing the session information for the website to retrieve the source code from.
//
// Returns:
//
// (string): the source code of the web page loaded in the site session.
//
// (error): an error if the function fails to retrieve the source code or if the driver is not of the correct type.
func GetPageSource(site web.Site) (string, error) {
	// Convert the driver to a chrome-specific *Driver instance
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return "", errors.New("driver is not of type *Driver")
	}

	// Get the page source code
	var pageSource string
	err := chromedp.Run(chromeDriver.Context, chromedp.OuterHTML("html", &pageSource))

	return pageSource, err
}

// Navigate navigates an input site using the provided InputActions.
func Navigate(site web.Site, actions []InputAction, waitTime time.Duration) error {
	// Convert the driver to a chrome-specific *Driver instance
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	// Enable network events
	if err := chromedp.Run(chromeDriver.Context, network.Enable()); err != nil {
		return err
	}

	// Set up request logging
	chromedp.ListenTarget(chromeDriver.Context, func(ev interface{}) {
		switch msg := ev.(type) {
		case *network.EventRequestWillBeSent:
			fmt.Printf("Request: %s %s\n", msg.Request.Method, msg.Request.URL)
			// Check if we have been redirected
			// if so, change the URL that we are tracking.
			if msg.RedirectResponse != nil {
				fmt.Printf("Encountered redirect: %s\n", msg.RedirectResponse.URL)
			}
		case *network.EventResponseReceived:
			fmt.Printf("Response URL: %s\n Response Headers: %s\n Response Status Code: %d\n",
				msg.Response.URL, msg.Response.Headers, msg.Response.Status)
		}
	})

	// Perform actions sequentially
	for i, inputAction := range actions {
		actionType := fmt.Sprintf("%T", inputAction.Action)
		if inputAction.Description != "" {
			fmt.Printf("Executing action #%d:\n Description: %s\nType: %s", i+1, inputAction.Description, actionType)
		}
		fmt.Printf("Executing action #%d:\n Type: %s", i+1, actionType)

		if err := chromedp.Run(chromeDriver.Context, chromedp.Tasks{
			inputAction.Action,
			chromedp.Sleep(waitTime),
		}); err != nil {
			return err
		}
	}

	return nil
}
