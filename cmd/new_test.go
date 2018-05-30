package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestNewTarget(t *testing.T) {
	workingDir, errDir := os.Getwd()
	if errDir != nil {
		// log.Fatal(errDir)
		t.Log("[ERROR]::", "Failed to fetch working directory\n\t", errDir)
	}

	//=================  - part 1 -  =================//
	// set tempestcfbup
	tempestcfbup := setTestTempestcf(t, []string{})

	// case no arg
	newTarget([]string{}...)

	// Get all paths added in Tempestcf
	allPaths, errAllP := getPaths()
	if errAllP != nil {
		errDel := os.Remove(Tempestcf)
		if errDel != nil {
			t.Log(errDel)
		}
		t.Error(errAllP)
	}

	// compare
	if !IsStringInSlice(workingDir+Slash+"temp.est", allPaths) {
		t.Log("[FAIL]:: Using current dir failed, result is different than expected:\n\tWanted:", workingDir+Slash+"temp.est", "\n\tGOT:\n\t\t", allPaths)
		t.Fail()
	}
	fmt.Println("::[PASS] NO_ARGS")
	// fallback tempestcfbup
	fbTestTempestcf(t, tempestcfbup)
	// speciffic clean
	if err := os.Remove(workingDir + Slash + "temp.est"); err != nil {
		t.Log("[ERROR]::", "Cannot remove", workingDir+Slash+"temp.est")
	}
	//=================  - part 2 -  =================//
	// case args
	testArgs := []string{
		conf.Gopath,
		conf.Home,
	}

	// set tempestcfbup
	tempestcfbup = setTestTempestcf(t, []string{})

	// Use the tested func
	autoGen = true
	newTarget(testArgs...)

	// Get all paths added in Tempestcf once again
	allPaths, errAllP = getPaths()
	if errAllP != nil {
		errDel := os.Remove(Tempestcf)
		if errDel != nil {
			t.Log(errDel)
		}
		t.Error(errAllP)
	}

	// compare
	for _, tstA := range testArgs {
		if !IsStringInSlice(tstA, allPaths) {
			t.Log("[FAIL]:: No", tstA, "in", allPaths)
			t.Fail()
		}
	}

	// speciffic clean
	if err := os.Remove(conf.Gopath + Slash + "temp.est"); err != nil {
		t.Log("[ERROR]::", "Cannot remove", conf.Gopath+Slash+"temp.est")
	}
	if err := os.Remove(conf.Home + Slash + "temp.est"); err != nil {
		t.Log("[ERROR]::", "Cannot remove", conf.Home+Slash+"temp.est")
	}
	// fallback tempestcfbup
	fbTestTempestcf(t, tempestcfbup)
}
