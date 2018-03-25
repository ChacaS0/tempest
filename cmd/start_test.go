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
