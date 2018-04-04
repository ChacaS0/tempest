package cmd

import (
	"testing"
)

// TestCallPurge is the test for callPurge(targets []string) error {}
// It is not a really sexy test, only checks if there are no errors when
// running the func
func TestCallPurge(t *testing.T) {
	// setting the place
	slTest := make([]string, 0)

	// call it with test mode on
	tTest = true
	if err := callPurge(slTest); err != nil {
		t.Log("[FAIL]:TESTMODE: An error occurred when trying to use ``callPurge(", slTest, ")\n\t->", err)
		t.Fail()
	}

	// call it without test mode on
	tTest = false
	if err := callPurge(slTest); err != nil {
		t.Log("[FAIL]:REGMODE: An error occurred when trying to use ``callPurge(", slTest, ")\n\t->", err)
		t.Fail()
	}
}

// TestHandleShutupMode checks if the func handles well the shutup mode.
// - use a test log file
// - check the creatation of log file if doesn't already exist, don't replace the existing, just append.
// - capture some stdout with test mode on and then check that the right values got written
func TestHandleShutupMode(t *testing.T) {
	// TODO:
}
