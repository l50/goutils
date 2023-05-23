package chrome

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
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

// SetContext sets the context associated with the Driver.
func (d *Driver) SetContext(ctx context.Context) {
	d.Context = ctx
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

// Init returns a chrome browser instance.
func Init(headless bool, ignoreCertErrors bool) (web.Browser, error) {
	var cancels []func()

	browser := web.Browser{
		Cancels: cancels,
	}

	options := []chromedp.ExecAllocatorOption{}
	setChromeOptions(&browser, headless, ignoreCertErrors, &options)

	driver, ok := browser.Driver.(*Driver)
	if !ok {
		err := errors.New("driver is not of type *ChromeDP")
		return web.Browser{}, err
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
		case *page.EventJavascriptDialogOpening:
			go func() {
				err := chromedp.Run(chromeDriver.Context,
					page.HandleJavaScriptDialog(true))

				if err != nil {
					log.Fatal(err)
				}
			}()
		case *network.EventRequestWillBeSent:
			// Check if we have been redirected
			// if so, change the URL that we are tracking.
			if msg.RedirectResponse != nil && site.Debug {
				fmt.Printf("Encountered redirect: %s\n", msg.RedirectResponse.URL)
			}
		case *network.EventResponseReceived:
			if site.Debug {
				fmt.Printf("Response URL: %s\n Response Headers: %s\n Response Status Code: %d\n",
					msg.Response.URL, msg.Response.Headers, msg.Response.Status)
			}
		}
	})

	// Perform actions sequentially
	for i, inputAction := range actions {
		actionType := fmt.Sprintf("%T", inputAction.Action)
		if inputAction.Description != "" && site.Debug {
			fmt.Printf("Executing action #%d:\n Description: %s\nType: %s", i+1, inputAction.Description, actionType)
		}
		if site.Debug {
			fmt.Printf("Executing action #%d:\n Type: %s", i+1, actionType)
		}

		if err := chromedp.Run(chromeDriver.Context, chromedp.Tasks{
			inputAction.Action,
			chromedp.Sleep(waitTime),
		}); err != nil {
			return err
		}
	}

	return nil
}

// ScreenShot takes a screenshot of the input `targetURL`
// and saves it to `imgPath`.
func ScreenShot(site web.Site, imgPath string) error {
	var screenshot []byte
	// Convert the driver to a chrome-specific *Driver instance
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	if err := chromedp.Run(chromeDriver.Context, takeSS(100, &screenshot)); err != nil {
		fmt.Errorf("failed to take screenshot: %v", err)
		return err
	}

	if err := os.WriteFile(imgPath, screenshot, 0644); err != nil {
		fmt.Errorf("failed to write screenshot to disk: %v", err)
		return err
	}

	return nil

}

func takeSS(quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {

			// force viewport emulation
			err := emulation.SetDeviceMetricsOverride(1280, 800, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(100).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  1280,
					Height: 800,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
