package cdpu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/l50/goutils/v2/web"
)

// CheckElement checks if a web page element, identified by the provided XPath,
// exists within a specified timeout.
//
// **Note:** Ensure to handle the error sent to the "done" channel in a
// separate goroutine or after calling this function to avoid deadlock.
//
// **Parameters:**
//
// site: A web.Site struct representing the target site.
// elementXPath: A string representing the XPath of the target element.
// done: A channel through which the function sends an error if the
// element is found or another error occurs.
//
// **Returns:**
//
// error: An error if the element is found, the web driver is not of
// type *Driver, failed to create a random wait time, or another error occurs.
func CheckElement(site web.Site, elementXPath string, done chan error) error {
	// Create a new context with a timeout
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	ctx, cancel := context.WithTimeout(chromeDriver.GetContext(), 10*time.Second)
	defer cancel()

	actions := []InputAction{
		{
			Description: fmt.Sprintf("Check if the element with XPath %s is present", elementXPath),
			Selector:    elementXPath,
			Action: chromedp.ActionFunc(func(ctx context.Context) error {
				go func() {
					var nodes []*cdp.Node
					err := chromedp.Run(ctx, chromedp.Nodes(elementXPath, &nodes, chromedp.BySearch))
					if err == nil && len(nodes) > 0 {
						err = fmt.Errorf("%s account is locked out", site.Session.Credential.User)
					}
					done <- err
				}()

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 5):
				}
				return nil
			}),
			Context: ctx,
		},
	}

	randomWaitTime, err := web.GetRandomWait(2, 6)
	if err != nil {
		return fmt.Errorf("failed to create random wait time: %v", err)
	}

	return Navigate(site, actions, randomWaitTime)
}

// Driver is an interface to Google Chrome, containing a context.Context
// associated with this Driver and Options for the execution of Google Chrome.
//
// **Attributes:**
//
// Context: The context associated with this Driver.
// Options: The options for the execution of Google Chrome.
type Driver struct {
	Context context.Context
	Options *[]chromedp.ExecAllocatorOption
}

// GetContext retrieves the context associated with the Driver instance.
//
// **Returns:**
//
// context.Context: The context associated with this Driver.
func (d *Driver) GetContext() context.Context {
	return d.Context
}

// SetContext associates a new context with the Driver instance.
//
// **Parameters:**
//
// ctx (context.Context): The new context to be associated with this Driver.
func (d *Driver) SetContext(ctx context.Context) {
	d.Context = ctx
}

// Init initializes a chrome browser instance with the specified headless mode and
// SSL certificate error ignoring options, then returns the browser instance for
// further operations.
//
// **Parameters:**
//
// headless (bool): Whether or not the browser should be in headless mode.
// ignoreCertErrors (bool): Whether or not SSL certificate errors should be ignored.
//
// **Returns:**
//
// web.Browser: An initialized Browser instance.
// error: Any error encountered during initialization.
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

// InputAction represents selectors and actions to run with Chrome. It
// contains a description, a selector to find an element on the page, and
// a chromedp.Action which defines the action to perform on the selected
// element.
//
// **Attributes:**
//
// Description: A string that describes the action.
// Selector: The CSS selector of the element to perform the action on.
// Action: A chromedp.Action that defines what action to perform.
// Context: The context in which to execute the action.
type InputAction struct {
	Description string
	Selector    string
	Action      chromedp.Action
	Context     context.Context
}

// GetPageSource retrieves the HTML source code of the currently loaded
// page in the provided Site's session.
//
// **Parameters:**
//
// site (web.Site): The site whose source code is to be retrieved.
//
// **Returns:**
//
// string: The source code of the currently loaded page.
// error: An error if any occurred during source code retrieval.
func GetPageSource(site web.Site) (string, error) {
	// Convert the driver to a chrome-specific *Driver instance
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return "", errors.New("driver is not of type *Driver")
	}

	// Get the page source code
	var pageSource string
	err := chromedp.Run(chromeDriver.GetContext(), chromedp.OuterHTML("html", &pageSource))

	return pageSource, err
}

// Navigate performs the provided actions sequentially on the provided Site's
// session. It enables network events and sets up request logging.
//
// **Parameters:**
//
// site (web.Site): The site on which the actions should be performed.
// actions ([]InputAction): A slice of InputAction objects which define
// the actions to be performed.
// waitTime (time.Duration): The time to wait between actions.
//
// **Returns:**
//
// error: An error if any occurred during navigation.
func Navigate(site web.Site, actions []InputAction, waitTime time.Duration) error {
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	if err := enableNetwork(chromeDriver); err != nil {
		return err
	}

	setUpRequestLogging(chromeDriver, &site)

	for i, inputAction := range actions {
		actionType := fmt.Sprintf("%T", inputAction.Action)
		if inputAction.Description != "" && site.Debug {
			fmt.Printf("Executing action #%d:\n Description: %s\nType: %s", i+1, inputAction.Description, actionType)
		}
		if site.Debug {
			fmt.Printf("Executing action #%d:\n Type: %s", i+1, actionType)
		}

		ctx := inputAction.Context
		if ctx == nil {
			ctx = chromeDriver.GetContext()
		}

		if err := chromedp.Run(ctx, chromedp.Tasks{
			inputAction.Action,
			chromedp.Sleep(waitTime),
		}); err != nil {
			return err
		}
	}

	return nil
}

// ScreenShot captures a screenshot of the currently loaded page in the
// provided Site's session and writes the image data to the provided file path.
//
// **Parameters:**
//
// site (web.Site): The site whose page a screenshot should be taken of.
// imgPath (string): The path to which the screenshot should be saved.
//
// **Returns:**
//
// error: An error if any occurred during screenshot capturing or saving.
func ScreenShot(site web.Site, imgPath string) error {
	var screenshot []byte
	// Convert the driver to a chrome-specific *Driver instance
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	if err := chromedp.Run(chromeDriver.Context, takeSS(100, &screenshot)); err != nil {
		return fmt.Errorf("failed to take screenshot: %v", err)
	}

	if err := os.WriteFile(imgPath, screenshot, 0644); err != nil {
		return fmt.Errorf("failed to write screenshot to disk: %v", err)
	}

	return nil
}

// enableNetwork enables network events
func enableNetwork(chromeDriver *Driver) error {
	return chromedp.Run(chromeDriver.Context, network.Enable())
}

// setChromeOptions is used to set the chrome
// parameters required by ChromeDP.
func setChromeOptions(browser *web.Browser, headless bool,
	ignoreCertErrors bool, options *[]chromedp.ExecAllocatorOption) {
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

// setUpRequestLogging sets up request logging
func setUpRequestLogging(chromeDriver *Driver, site *web.Site) {
	chromedp.ListenTarget(chromeDriver.Context, func(ev interface{}) {
		switch msg := ev.(type) {
		case *page.EventJavascriptDialogOpening:
			go func() {
				if err := chromedp.Run(chromeDriver.Context,
					page.HandleJavaScriptDialog(true)); err != nil {
					fmt.Printf("error handling JavaScript dialog: %v", err)
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

// SaveCookiesToDisk retrieves cookies from the current session and writes them to a file.
//
// **Parameters:**
//
// site (web.Site): The site from which to retrieve cookies.
// filePath (string): The file path where the cookies should be saved.
//
// **Returns:**
//
// error: An error if any occurred during cookie retrieval or file writing.
func SaveCookiesToDisk(site web.Site, filePath string) error {
	chromeDriver, ok := site.Session.Driver.(*Driver)
	if !ok {
		return errors.New("driver is not of type *Driver")
	}

	ctx := chromeDriver.GetContext()

	// Create an instance of GetCookiesParams
	cookiesParams := network.GetCookies()

	// Optional: Specify URLs if needed
	// cookiesParams.Urls = []string{"<your-url-here>"}

	// Use the new GetCookies method
	var cookies []*network.Cookie
	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = cookiesParams.Do(ctx)
		return err
	}))
	if err != nil {
		return err
	}

	// Marshal cookies into JSON
	cookieJSON, err := json.Marshal(cookies)
	if err != nil {
		return err
	}

	// Write the JSON to the specified file
	return os.WriteFile(filePath, cookieJSON, 0644)
}
