package cmd

import (
	"os"
	"testing"
)

// TestGetPaths() test function for list/getPaths() (returnSlice []string, pathsError error) {}
func TestGetPaths(t *testing.T) {
	// presests

	// Slice we will use for tests
	slTest := []string{
		conf.Gobin,
		conf.Gopath,
	}

	tempestcfbup := setTestTempestcf(t, slTest)

	// try to use the func
	allPaths, err := getPaths()
	if err != nil {
		if err.Error() != "empty" {
			t.Log("[ERROR]:: Could not read from the file .tempestcf and returned with this error:\n\t->", err)
		} else {
			t.Log("[ERROR]::", Tempestcf, " is empty !!\n\t->", err)
		}
		t.Fail()
	}

	// Check if the paths we get with getPaths() are the same as slTest
	if !SameSlices(allPaths, slTest) {
		t.Log("[FAIL]:: getPAths() is not returning the right data:\n\t-> Tempestcf:", Tempestcf, "\n\t-> getPaths():", allPaths, "\n\tslTest:", slTest)
		t.Fail()
	}

	// Fallback
	fbTestTempestcf(t, tempestcfbup)
}

// TestPrintList checks if printList() error {} prints well the slice
// given by getPaths() ([]string, error) {}
func TestPrintList(t *testing.T) {
	// // HINT: Change stdout to a variable to check the result?
	// // 		Will probably be in a []bytes

	// bup of current Tempestcf
	tempestcfbup := Tempestcf
	// new testing Tempestcf
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	tmpcf, errCreate := os.OpenFile(Tempestcf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errCreate != nil {
		t.Log("[ERROR]:: No file create, but we tried !! #sadface:\n\t->", Tempestcf, "\n\t->", errCreate)
		Tempestcf = tempestcfbup
		t.FailNow()
	}
	defer tmpcf.Close()

	// Before adding lines, test the case of empty .tempestcf
	emptyOutput := captureStdout(func() {
		if err := printList(); err != nil {
			t.Log("[ERROR]:: Can't printList()")
			t.Fail()
		}
	})

	// Verify output
	wantedOutput := ":: No target set yet\n:: Suggestion - Run: \n\ttempest help add\nFor more information about adding targets!\n"
	if wantedOutput != emptyOutput {
		t.Log("[FAIL]:: printList() failed to process empty .tempestcf\n\t-> Expected:\n\t", wantedOutput, "\n\t-> Got:\n\t", emptyOutput)
		t.Fail()
	}

	// Slice we will use for tests
	slTest := []string{
		conf.Gobin,
		conf.Gopath,
	}

	// Add slTest data to Tempestcf
	if err := addLine(slTest); err != nil {
		t.Log("[ERROR]:: Can't add lines to", Tempestcf, ":\n\t->", err)
		t.Fail()
	}

	// Try to use the func and captures the output of it
	actualOutput := captureStdout(func() {
		if err := printList(); err != nil {
			t.Log("[ERROR]:: Can't printList()")
			t.Fail()
		}
	})

	// Verify output of printList()
	wantedOutput = "Current targets currently having \"fun\" with TEMPest:\n\nIndex\t| Target\n0\t| " + conf.Gobin + "\n1\t| " + conf.Gopath + "\n"
	if actualOutput != wantedOutput {
		t.Log("[FAIL]:: The output of printList() was quite unexpected! Wow!\n\t-> ActualOutput:\n\t", actualOutput, "\n\t-> Wanted:\n\t", wantedOutput)
		t.Fail()
	}

	// Fallback
	fbTestTempestcf(t, tempestcfbup)
}
