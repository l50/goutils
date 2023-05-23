package web_test

import (
	"testing"
	"time"

	"github.com/l50/goutils/v2/web"
)

func TestWait(t *testing.T) {
	near := 1000.0
	waitTime := float64(web.Wait(near)) / float64(time.Millisecond)

	if waitTime < 0.95*near || waitTime > near+0.1*near {
		t.Errorf("Wait time is out of range. Expected between %v and %v, got %v.", 0.95*near, near+0.1*near, waitTime)
	}
}

func TestGetRandomWait(t *testing.T) {
	minWait := 2 * time.Second
	maxWait := 6 * time.Second
	randomWaitTime := web.GetRandomWait(minWait, maxWait)

	if randomWaitTime < minWait || randomWaitTime > maxWait {
		t.Errorf("Random wait time is out of range. Expected between %v and %v, got %v.", minWait, maxWait, randomWaitTime)
	}
}
