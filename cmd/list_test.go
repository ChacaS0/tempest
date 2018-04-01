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

// TestFixTargets is the func that checks if it deletes all the broken targets
func TestFixTargets(t *testing.T) {
	// Create a temp dir for the test
	testPath := conf.Gopath + string(os.PathSeparator) + "test"
	if err := os.Mkdir(testPath, 0777); err != nil {
		t.Log("[ERROR]:: Could not create the test directory: ", testPath, "\n\t->", err)
		t.Fail()
	}

	// Add targets to TEMPest
	tempestcfbup := setTestTempestcf(t, []string{conf.Gobin, conf.Gopath, testPath})

	// Get paths added
	allPaths, errGP := getPaths()
	if errGP != nil {
		t.Log("[ERROR]:: Could not retrieve the paths added:\n\t", errGP)
	}

	// We can test here if nothing changes when it's all good
	if err := fixTargets(); err != nil {
		t.Log("[ERROR]:: Sorry it spotted broken targets and removed them when there was supposed to be none.\n\t->ERROR:: ", err)
		t.Fail()
	}
	// Real check
	pathsAfter, errGPAfter := getPaths()
	if errGP != nil {
		t.Log("[ERROR]:: Could not retrieve the paths added:\n\t", errGPAfter)
		t.Fail()
	}
	if !SameSlices(allPaths, pathsAfter) {
		t.Log("[FAIL]:: Removed stuff it should have to.\n\tGOT: ", pathsAfter, "\n\tWANT: ", allPaths)
		t.Fail()
	}

	// Typical broken path is on who got deleted on the system but not in TEMPest
	// So that's we are going to simulate
	// Deletion of testPath
	if err := os.Remove(testPath); err != nil {
		t.Log("[ERROR]:: Fuck could not delete the f*cking testPath, can't test like this!\n\t", err)
		t.Fail()
	}

	// Now we can test the func to see if it does report a broken Target (testPath)
	if err := fixTargets(); err != nil {
		t.Log("[ERROR]:: Sorry the fix was wrong (did not take care of the broken target)!\n\t->ERROR: ", err)
		t.Fail()
	}
	// Real check
	pathsAfter, errGPAfter = getPaths()
	if errGP != nil {
		t.Log("[ERROR]:: Could not retrieve the paths added:\n\t", errGPAfter)
		t.Fail()
	}
	// preparing results
	cpt := 0 // we want it to get to 3
	for _, pa := range pathsAfter {
		switch pa {
		case conf.Gobin:
			cpt++
		case conf.Gopath:
			cpt++
		}
	}
	if !IsStringInSlice(testPath, pathsAfter) {
		cpt++
	}
	// Handle results
	if cpt != 3 {
		t.Log("[FAIL]:: Damn, something went terribly wrong!!\n\tGOT: ", pathsAfter, "\n\tWANT: ", []string{conf.Gobin, conf.Gopath})
		t.Fail()
	}

	// Fall back to repvious .tempestcf
	fbTestTempestcf(t, tempestcfbup)

}

// TestGetState is the func that checks if getState give the right state.
// Meaning => { True : "good", False: "Not good"}
func TestGetState(t *testing.T) {
	// Create a temp dir for the test
	testPath := conf.Gopath + string(os.PathSeparator) + "test"
	if err := os.Mkdir(testPath, 0777); err != nil {
		t.Log("[ERROR]:: Could not create the test directory: ", testPath, "\n\t->", err)
		t.Fail()
	}

	// Add targets to TEMPest
	tempestcfbup := setTestTempestcf(t, []string{conf.Gobin, conf.Gopath, testPath})

	// Get paths added
	allPaths, errGP := getPaths()
	if errGP != nil {
		t.Log("[ERROR]:: Could not retrieve the paths added:\n\t", errGP)
	}
	// String paths to targets
	allTargets := PathsToTargets(allPaths)

	// retrieve the individual Targets
	var testTarget Target
	var testGopath Target
	var testGobin Target
	for _, tgt := range allTargets {
		switch tgt.Path {
		case testPath:
			testTarget = tgt
		case conf.Gobin:
			testGobin = tgt
		case conf.Gopath:
			testGopath = tgt
		}
	}

	// We can test here if all targets are "good"
	stateTgts := getState(allTargets)
	if !stateTgts[testTarget] || !stateTgts[testGobin] || !stateTgts[testGopath] {
		t.Log("[FAIL]:: Sorry it spotted broken targets when there was supposed to be none.\n\t->GOT: ", stateTgts)
		t.Fail()
	}

	// Typical broken path is on who got deleted on the system but not in TEMPest
	// So that's we are going to simulate
	// Deletion of testPath
	if err := os.Remove(testPath); err != nil {
		t.Log("[ERROR]:: Fuck could not delete the f*cking testPath, can't test like this!\n\t", err)
		t.Fail()
	}

	// Now we can test the func to see if it does report a broken Target (testPath)
	stateTgts = getState(allTargets)
	if stateTgts[testTarget] || !stateTgts[testGobin] || !stateTgts[testGopath] {
		t.Log("[FAIL]:: Sorry the states were wrong (did not spot the broken target)!\n\t->GOT: ", stateTgts)
		t.Fail()
	}

	// Fall back to repvious .tempestcf
	fbTestTempestcf(t, tempestcfbup)

}
