package web_test

import (
	"context"

	"github.com/l50/goutils/v2/web"
)

func ExampleCancelAll() {
	var cancels []func()
	_, cancel := context.WithCancel(context.Background())
	cancels = append(cancels, cancel)

	// Later, when all operations need to be cancelled:
	web.CancelAll(cancels...)
	// Output:
}

func ExampleGetRandomWait() {
	minWait := 2
	maxWait := 6
	randomWaitTime, _ := web.GetRandomWait(minWait, maxWait)
	_ = randomWaitTime
}

func ExampleIsTwoFacEnabled() {
	options := []web.LoginOption{
		web.WithTwoFac(false),
	}
	loginOpts := web.SetLoginOptions(options...)
	isTwoFacEnabled := web.IsTwoFacEnabled(loginOpts)
	_ = isTwoFacEnabled
}

func ExampleIsLogMeOutEnabled() {
	options := []web.LoginOption{
		web.WithLogout(false),
	}
	loginOpts := web.SetLoginOptions(options...)
	isLogMeOutEnabled := web.IsLogMeOutEnabled(loginOpts)
	_ = isLogMeOutEnabled
}

func ExampleSetLoginOptions() {
	options := []web.LoginOption{
		web.WithTwoFac(false),
		web.WithLogout(true),
	}
	loginOpts := web.SetLoginOptions(options...)
	_ = loginOpts // use loginOpts
}
