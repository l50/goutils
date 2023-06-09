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
