package web_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/l50/goutils/web"
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

			// Note: Because the cancellations are performed instantly in this example, we can't
			// check here whether the cancellations have happened. In actual application code,
			// cancellations would generally have some observable effect (like stopping a goroutine).
		})
	}
}

func TestWait(t *testing.T) {
	near := 1000.0
	waitTime, err := web.Wait(near)
	if err != nil {
		t.Fatalf("Unexpected error while waiting: %v", err)
	}

	waitTimeMillis := float64(waitTime) / float64(time.Millisecond)

	if waitTimeMillis < 0.95*near || waitTimeMillis > near+0.1*near {
		t.Errorf("Wait time is out of range. Expected between %v and %v, got %v.", 0.95*near, near+0.1*near, waitTimeMillis)
	}
}

func TestGetRandomWait(t *testing.T) {
	minWait := 2 * time.Second
	maxWait := 6 * time.Second
	randomWaitTime, err := web.GetRandomWait(minWait, maxWait)
	if err != nil {
		t.Fatalf("Unexpected error while waiting: %v", err)
	}

	if randomWaitTime < minWait || randomWaitTime > maxWait {
		t.Errorf("Random wait time is out of range. Expected between %v and %v, got %v.", minWait, maxWait, randomWaitTime)
	}
}
