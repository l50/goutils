package chrome_test

import (
	"testing"

	"github.com/l50/goutils/v2/web/chrome"
)

func TestInit(t *testing.T) {
	headless := false
	ignoreCertErrors := true
	browser, err := chrome.Init(headless, ignoreCertErrors)
	if err != nil {
		t.Errorf("Failed to initialize browser: %v", err)
	}

	if browser.Driver == nil {
		t.Error("Browser driver is nil")
	}
}
