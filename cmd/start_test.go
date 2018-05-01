package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
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

	//! call it without test mode on
	tTest = false
	if err := callPurge(slTest); err != nil {
		t.Log("[FAIL]:REGMODE: An error occurred when trying to use ``callPurge(", slTest, ")\n\t->", err)
		t.Fail()
	}

	// =============
	// Test file to be added to tempestcf for testing purpose
	f, errCF := os.Create(conf.Gobin + string(os.PathSeparator) + "test_file.temp")
	if errCF != nil {
		t.Log("[ERROR]:CREATE_F: Can't create a file for testing")
		t.Fail()
	}
	f.Close()

	// Setup the test .tempestcf
	tempestcfbup := setTestTempestcf(t, []string{
		conf.Gobin,
		conf.Gobin + string(os.PathSeparator) + "test_file.temp",
	})

	// call it with test mode on
	tTest = true
	if err := callPurge(slTest); err != nil {
		t.Log("[FAIL]:REGMODE: An error occurred when trying to use ``callPurge(", slTest, ")\n\t->", err)
		t.Fail()
	}

	//! don't call it without test mode on, it might delete stuff !

	// fallback to the previous tempestcf config
	fbTestTempestcf(t, tempestcfbup)
}

// TestHandleShutupMode checks if the func handles well the shutup mode.
// - use a test log file
// - check the creatation of log file if doesn't already exist, don't replace the existing, just append.
// - capture some stdout with test mode on and then check that the right values got written
func TestHandleShutupMode(t *testing.T) {
	// presets for test
	logShutupbup, _ := setTestLogShutup(t)
	tempestcfbup := setTestTempestcf(t, []string{})

	// activate the test mode for TEMPest
	testAll = true

	// capture what we should get for a test mode shutup in the log file, yeah much explicit
	// First one - empty
	want1 := captureStdout(func() {
		if err := callPurge([]string{}); err != nil {
			fmt.Println(err)
		}
	})
	want1 = HeaderLog + "\n" + want1 + "\n" + FooterLog + "\n"

	// try to use shutup mode
	handleShutupMode([]string{})

	// fetch what was written
	fileCtnt, errRF := ioutil.ReadFile(LogShutup)
	if errRF != nil {
		t.Log("[ERROR]:: Could not read the file -_-\n\t->", LogShutup, "\n\t->", errRF)
		t.Fail()
	}
	fileCtntStr := string(fileCtnt)
	// compare
	if fileCtntStr != want1 {
		t.Log("[FAIL]:1: Test failed cause we didn't get what we wanted:\n\t[WANT]\n", want1, "\n\t[GOT]\n", fileCtntStr)
		t.Fail()
	} else {
		fmt.Println("[SUCCESS]:: First test succeed!")
	}

	//============================
	// Set up second test
	fbTestTempestcf(t, tempestcfbup)
	tempestcfbup = setTestTempestcf(t, []string{conf.Gopath, conf.Gobin, conf.Home})
	fbTestLogShutup(t, logShutupbup)
	logShutupbup, _ = setTestLogShutup(t)

	// get some stuff
	targets, errTgts := getPaths()
	if errTgts != nil {
		t.Log("[ERROR]:: Error while trying to get some f*cking paths")
		t.Fail()
	}

	// capture what we should get for a test mode shutup in the log file, yeah much explicit
	// Second one - not empty
	want2 := captureStdout(func() {
		if err := callPurge(targets); err != nil {
			fmt.Println(err)
		}
	})
	want2 = HeaderLog + "\n" + want2 + "\n" + FooterLog + "\n"

	// try to use shutup mode
	handleShutupMode(targets)

	// fetch what was written
	fileCtnt, errRF = ioutil.ReadFile(LogShutup)
	if errRF != nil {
		t.Log("[ERROR]:: Could not read the file -_-\n\t->", LogShutup, "\n\t->", errRF)
		t.Fail()
	}
	fileCtntStr = string(fileCtnt)
	// compare
	if fileCtntStr != want2 {
		t.Log("[FAIL]:2: Test failed cause we didn't get what we wanted:\n\t[WANT]\n", want2, "\n\t[GOT]\n", fileCtntStr)
		t.Fail()
	}

	// clean up the mess
	testAll = false
	fbTestLogShutup(t, logShutupbup)
	fbTestTempestcf(t, tempestcfbup)
}
