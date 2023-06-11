package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/l50/goutils/v2/web"
)

func TestCancelAll(t *testing.T) {
	tests := []struct {
		name    string
		cancels []func()
	}{
		{
			name: "Test with actual cancellation",
			cancels: []func(){
				func() {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						<-ctx.Done()
						fmt.Println("Operation 1 cancelled")
					}()
					cancel()
				},
				func() {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						<-ctx.Done()
						fmt.Println("Operation 2 cancelled")
					}()
					cancel()
				},
				func() {
					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						<-ctx.Done()
						fmt.Println("Operation 3 cancelled")
					}()
					cancel()
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function being tested.
			web.CancelAll(tc.cancels...)
		})
	}
}

func TestGetRandomWait(t *testing.T) {
	tests := []struct {
		name    string
		minWait int
		maxWait int
		wantErr bool
	}{
		{
			name:    "normal case",
			minWait: 2,
			maxWait: 6,
		},
		{
			name:    "negative min wait",
			minWait: -2,
			maxWait: 6,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := web.GetRandomWait(tc.minWait, tc.maxWait)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetRandomWait() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func TestIsTwoFacEnabled(t *testing.T) {
	tests := []struct {
		name    string
		options []web.LoginOption
		want    bool
	}{
		{
			name:    "default values",
			options: nil,
			want:    true,
		},
		{
			name:    "two-factor disabled",
			options: []web.LoginOption{web.WithTwoFac(false)},
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loginOpts := web.SetLoginOptions(tc.options...)
			if got := web.IsTwoFacEnabled(loginOpts); got != tc.want {
				t.Errorf("IsTwoFacEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsLogMeOutEnabled(t *testing.T) {
	tests := []struct {
		name    string
		options []web.LoginOption
		want    bool
	}{
		{
			name:    "default values",
			options: nil,
			want:    true,
		},
		{
			name:    "logout disabled",
			options: []web.LoginOption{web.WithLogout(false)},
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loginOpts := web.SetLoginOptions(tc.options...)
			if got := web.IsLogMeOutEnabled(loginOpts); got != tc.want {
				t.Errorf("IsLogMeOutEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSetLoginOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     []web.LoginOption
		twoFacValue bool
		logoutValue bool
	}{
		{
			name:        "default values",
			options:     nil, // no options passed
			twoFacValue: true,
			logoutValue: true,
		},
		{
			name:        "two-factor disabled",
			options:     []web.LoginOption{web.WithTwoFac(false)},
			twoFacValue: false,
			logoutValue: true,
		},
		{
			name:        "logout disabled",
			options:     []web.LoginOption{web.WithLogout(false)},
			twoFacValue: true,
			logoutValue: false,
		},
		{
			name:        "both two-factor and logout disabled",
			options:     []web.LoginOption{web.WithTwoFac(false), web.WithLogout(false)},
			twoFacValue: false,
			logoutValue: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// This assumes there's an exported function in web package that allows checking the options
			// As noted above, this is just for demonstrating how the tests would be structured.
			// In a real project, you'd test exported functions that use the loginOptions.
			loginOpts := web.SetLoginOptions(tc.options...)
			if got := web.IsTwoFacEnabled(loginOpts); got != tc.twoFacValue {
				t.Errorf("got IsTwoFacEnabled() = %v, want %v", got, tc.twoFacValue)
			}
			if got := web.IsLogMeOutEnabled(loginOpts); got != tc.logoutValue {
				t.Errorf("got IsLogMeOutEnabled() = %v, want %v", got, tc.logoutValue)
			}
		})
	}
}

func TestWait(t *testing.T) {
	tests := []struct {
		name    string
		near    float64
		wantErr bool
	}{
		{
			name: "normal case",
			near: 1000.0,
		},
		{
			name:    "negative near",
			near:    -1000.0,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := web.Wait(tc.near)
			if (err != nil) != tc.wantErr {
				t.Errorf("Wait() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
