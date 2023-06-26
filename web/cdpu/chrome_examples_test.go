package cdpu_test

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/l50/goutils/v2/web"
	"github.com/l50/goutils/v2/web/cdpu"
)

func ExampleCheckElement() {
	// Initialize the chrome browser
	browser, err := cdpu.Init(true, true)
	if err != nil {
		log.Fatalf("failed to initialize a chrome browser: %v", err)
	}

	defer web.CancelAll(browser.Cancels...)

	url := "https://somesite.com/login"

	// Set up the site with the browser's driver
	site := web.Site{
		LoginURL: url,
		Session: web.Session{
			Driver: browser.Driver,
		},
	}

	// Define the XPath for the element to check
	elementXPath := "//button[@id='login']"

	// Create a done channel
	done := make(chan error)

	// Call the function in a goroutine and wait for result
	go func() {
		err := cdpu.CheckElement(site, elementXPath, done)
		if err != nil {
			log.Printf("failed to execute CheckElement: %v", err)
		}
	}()

	// Handle the result from the done channel
	select {
	case err := <-done:
		if err != nil {
			log.Printf("Element found or another error occurred: %v", err)
		}
	case <-time.After(10 * time.Second):
		log.Println("Timeout exceeded while waiting for element check")
	}
}

func ExampleDriver_GetContext() {
	d := &cdpu.Driver{}
	ctx := d.GetContext()

	if ctx == nil {
		log.Fatalf("context is nil")
	}
}

func ExampleDriver_SetContext() {
	d := &cdpu.Driver{}
	newCtx := context.Background()
	d.SetContext(newCtx)

	if d.GetContext() != newCtx {
		log.Fatalf("failed to set new context")
	}
}

func ExampleInit() {
	browser, err := cdpu.Init(true, true)
	if err != nil {
		log.Fatalf("failed to initialize a chrome browser: %v", err)
	}

	_ = browser
}

func ExampleInputAction() {
	action := cdpu.InputAction{
		Description: "Type in search box",
		Selector:    "#searchbox",
		Action:      chromedp.SendKeys("#searchbox", "example search"),
	}

	_ = action
}

func ExampleGetPageSource() {
	site := web.Site{
		// initialize site
	}

	source, err := cdpu.GetPageSource(site)
	if err != nil {
		log.Fatalf("failed to get page source: %v", err)
	}

	_ = source
}

func ExampleNavigate() {
	actions := []cdpu.InputAction{
		// initialize actions
	}

	site := web.Site{
		// initialize site
	}

	if err := cdpu.Navigate(site, actions, 1000); err != nil {
		log.Fatalf("failed to navigate site: %v", err)
	}
}

func ExampleScreenShot() {
	site := web.Site{
		// initialize site
	}

	if err := cdpu.ScreenShot(site, "/path/to/save/image.png"); err != nil {
		log.Fatalf("failed to capture screenshot: %v", err)
	}
}
