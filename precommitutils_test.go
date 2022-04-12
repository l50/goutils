package utils

import "testing"

func TestInstallGoPCDeps(t *testing.T) {
	if err := InstallGoPCDeps(); err != nil {
		t.Fatal(err)
	}
}

func TestInstallPCHooks(t *testing.T) {
	if err := InstallPCHooks(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePCHooks(t *testing.T) {
	if err := UpdatePCHooks(); err != nil {
		t.Fatal(err)
	}
}

// Currently not running because I need to figure out how to stop
// this test from calling itself repeatedly.

// func TestRunPCHooks(t *testing.T) {
// 	seconds := 1
// 	timeout := time.After(time.Duration(seconds) * time.Second)
// 	errors := make(chan error)
// 	done := make(chan bool)
// 	go func() {
// 		if err := RunPCHooks(); err != nil {
// 			errors <- fmt.Errorf("received error from RunHooks(): %v", err)
// 		}
// 		done <- true
// 	}()

// 	select {
// 	case <-timeout:
// 		fmt.Printf("timed out TestRunPCHooks() after %ds to stop the test from infinitely calling itself",
// 			seconds)
// 	case <-done:
// 	case err := <-errors:
// 		t.Fatalf("test fail - TestRunPCHooks() in precommitutils_test.go: %v", err)
// 	}
// }
