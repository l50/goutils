package chrome

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/l50/goutils/web"
)

// Driver is used to interface with Google Chrome using go.
//
// Contains:
//
// Context: A context.Context that's associated with this Driver.
//
// Options: Options for execution of Google Chrome.
type Driver struct {
	Context context.Context
	Options *[]chromedp.ExecAllocatorOption
}

// GetContext returns the context associated with the Driver.
//
// Returns:
//
// context.Context: The context that's associated with this Driver.
func (d *Driver) GetContext() context.Context {
	return d.Context
}

// SetContext sets the context associated with the Driver.
//
// Parameters:
//
// ctx (context.Context): The new context to be associated with this Driver.
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
//
// This function initializes a chrome browser instance with the specified headless mode and SSL certificate error ignoring options.
// The browser instance is then returned for further operations.
//
// Parameters:
//
// headless (bool): Whether or not the browser should be in headless mode.
//
// ignoreCertErrors (bool): Whether or not SSL certificate errors should be ignored.
//
// Returns:
//
// web.Browser: A Browser instance which has been initialized.
//
// error: Any encountered error during initialization.
func Init(headless bool, ignoreCertErrors bool) (web.Browser, error) {
	var cancels []func()

	browser := web.Browser{
		Cancels: cancels,
	}

	options := []chromedp.ExecAllocatorOption{}
	setChromeOptions(&browser, headless, ignoreCertErrors, &options)

	// Create contexts and their associated cancels.
	allocatorCtx, cancel := chromedp.NewExecAllocator(
		context.Background(), *browser.Driver.GetOptions())
	browser.Cancels = append([]func(){cancel}, cancels...)
	browser.Driver.SetContext(chromedp.NewContext(allocatorCtx,
		chromedp.WithLogf(log.Printf)))

	return browser, nil
}

// InputAction contains selectors and actions to run with chrome.
//
// Contains:
//
// Description: A description of this action.
//
// Selector: A selector to find an element on the page.
//
// Action: An chromedp.Action which defines the action to perform on the selected element.
type InputAction struct {
	Description string
	Selector    string
	Action      chromedp.Action
}

// GetPageSource retrieves the source code of the web page currently loaded in the site session.
//
// This function will return the HTML source code of the currently loaded page in the provided Site's session.
//
// Parameters:
//
// site (web.Site): The site whose source code is to be retrieved.
//
// Returns:
//
// string: The source code of the currently loaded page.
//
// error: An error if any occurred during source code retrieval.
func GetPageSource(site web.Site) (string, error) {
	chromeDriver := site.Session.Driver
	// Get the page source code
	var pageSource string
	err := chromedp.Run(chromeDriver.GetContext(), chromedp.OuterHTML("html", &pageSource))

	return pageSource, err
}

// Navigate navigates an input site using the provided InputActions.
//
// This function will perform the provided actions sequentially on the provided Site's session. It enables network events and sets up request logging.
//
// Parameters:
//
// site (web.Site): The site on which the actions should be performed.
//
// actions ([]InputAction): A slice of InputAction objects which define the actions to be performed.
//
// waitTime (time.Duration): The time to wait between actions.
//
// Returns:
//
// error: An error if any occurred during navigation.
func Navigate(site web.Site, actions []InputAction, waitTime time.Duration) error {
	chromeDriver := site.Session.Driver

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

// ScreenShot takes a screenshot of the input `targetURL` and saves it to `imgPath`.
//
// This function captures a screenshot of the currently loaded page in the provided Site's session and writes the image data to the provided file path.
//
// Parameters:
//
// site (web.Site): The site whose page a screenshot should be taken of.
//
// imgPath (string): The path to which the screenshot should be saved.
//
// Returns:
//
// error: An error if any occurred during screenshot capturing or saving.
func ScreenShot(site web.Site, imgPath string) error {
	var screenshot []byte
	chromeDriver := site.Session.Driver

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
