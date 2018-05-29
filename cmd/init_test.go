package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

// TestInitializeTP tests the initialization of
func TestInitializeTP(t *testing.T) {
	// Save current Tempestcf
	tempestcfOld := Tempestcf
	// Set the temporary .tempestcf used for the test
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	//* Case doesn't exist yet, should create one
	err := initializeTP()
	if err != nil {
		t.Log("[ERROR]:: Error while initializing on non-existing .tempestcf\n\t", err)
		t.Fail()
	}

	verifyTempestcf(t, tempestcfOld)

	//* If we get here we assume a file already exists and we will replace it
	// so we run initializeTP again
	err = initializeTP()
	if err != nil {
		t.Log("[ERROR]:: Error while initializing on existing .tempestcf\n\t", err)
		t.Fail()
	}

	verifyTempestcf(t, tempestcfOld)

	if err = cleanTempest(t, &Tempestcf, tempestcfOld); err != nil {
		t.Log("[ERROR]:: An error occurred:", err)
	}
}

// TestInitializeCfFile test the initializeCfFile
// meaning the reset of it if it already exists
// and its creation if it doesn't exist yet.
func TestInitializeCfFile(t *testing.T) {
	// Save current Tempestyml
	tempestymlOld := Tempestyml
	// Set the temporary .tempestyml used for the test
	Tempestyml = conf.Gopath + string(os.PathSeparator) + ".tempest.yaml"
	// viper.SetConfigFile(Tempestyml)

	//* Case doesn't exist yet, should create one
	err := initializeCfFile()
	if err != nil {
		t.Log("[ERROR]:: Error while initializing on non-existing .tempest.yml\n\t", err)
		t.Fail()
	}

	verifyTempestyml(t, tempestymlOld)

	//* If we get here we assume a file already exists and we will replace it
	// so we run initializeCfFile again
	err = initializeCfFile()
	if err != nil {
		t.Log("[ERROR]:: Error while initializing on existing .tempest.yml\n\t", err)
		t.Fail()
	}

	verifyTempestyml(t, tempestymlOld)

	// Set back to default conf and clean temp test files
	if err = cleanTempest(t, &Tempestyml, tempestymlOld); err != nil {
		t.Log("[ERROR]:: An error occurred:", err)
		t.Fail()
	}
	viper.SetConfigFile(Tempestyml)
}

// verifyTempestyml checks if the initializeCfFile() worked
func verifyTempestyml(t *testing.T, tempestymlOld string) {
	// Check if it really got created
	if isD, err := IsDirectory(Tempestyml); err != nil || isD {
		t.Log("[FAIL]:: Nothing got created\n\t", err)
		// Restore previous Tempestyml
		Tempestyml = tempestymlOld
		t.FailNow() // end it now
	}

	//* Try to use the file a bit
	// Check config file used
	if viper.ConfigFileUsed() != Tempestyml {
		t.Log("[NOTE]:: Config file was wrong, rectiffying ...")
		viper.SetConfigFile(Tempestyml)
	}

	age = viper.GetInt("duration") + 5 // should return 10

	setAge()

	// Check the age we just set
	if viper.GetInt("duration") != age {
		t.Log("[FAIL]:: Did not set the age correctly")
		t.Fail()
	}
}

// verifyTempestcf checks if the initializeTP() worked
func verifyTempestcf(t *testing.T, tempestcfOld string) {
	// Check if really got created
	if isD, err := IsDirectory(Tempestcf); err != nil || isD {
		t.Log("[FAIL]:: Nothing got created\n\t", err)
		// Restore previous Tempestcf
		Tempestcf = tempestcfOld
		t.FailNow()
	}

	// Try to use the file a bit
	addCmd.Run(addCmd, []string{
		conf.Gopath,
		conf.Gobin,
	})

	// if err := addCmd.Execute(); err != nil {
	// 	t.Log("[FAIL]:: Something went very wrong, please do something!\n\t", err)
	// 	t.Fail()
	// }
}

// cleanTempest is an internal function used to clean up the mess
// created while testing
func cleanTempest(t *testing.T, tempest *string, tempestOld string) (errReturn error) {
	// Delete the file created // ignore if file doesn't exist
	if err := os.Remove(*tempest); err != nil && err != os.ErrNotExist {
		errReturn = err
		t.Log("[ERROR]:: There was an error when trying to delete the freshly created file")
	}

	// Restore previous Tempest
	*tempest = tempestOld

	return
}

// setTestTempestcf set some presets for testing
func setTestTempestcf(t *testing.T, slTest []string) (tempestcfbup string) {
	// bup of current Tempestcf
	tempestcfbup = Tempestcf
	// new testing Tempestcf
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	tmpcf, errCreate := os.OpenFile(Tempestcf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errCreate != nil {
		t.Log("[ERROR]:: No file create, but we tried !! #sadface:\n\t->", Tempestcf, "\n\t->", errCreate)
		Tempestcf = tempestcfbup
		t.FailNow()
	}
	defer tmpcf.Close()

	// Add slTest data to Tempestcf
	if err := addLine(slTest); err != nil {
		t.Log("[ERROR]:: Can't add lines to", Tempestcf, ":\n\t->", err)
		t.Fail()
	}

	return
}

// setTestLogShutup set some presets for testing
func setTestLogShutup(t *testing.T) (logShutupbup string, nbBytesWritten int) {
	// bup of current Logfile
	logShutupbup = LogShutup
	// new testing Logfile
	LogShutup = conf.Gopath + string(os.PathSeparator) + ".shutup.log"

	logSUF, errCreate := os.OpenFile(LogShutup, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errCreate != nil {
		t.Log("[ERROR]:: No file create, but we tried !! #sadface:\n\t->", LogShutup, "\n\t->", errCreate)
		LogShutup = logShutupbup
		t.FailNow()
	}
	defer logSUF.Close()

	return
}

// fbTestLogShutup falls back to the previous LogShutup config
func fbTestLogShutup(t *testing.T, logShutupbup string) {
	if err := os.Remove(LogShutup); err != nil {
		t.Log("[ERROR]:: An error occurred when trying to remove the test logshutup:", LogShutup)
		t.Fail()
	}
	LogShutup = logShutupbup
}

// fbTestTempestcf falls back to the previous TEMPestcf config
func fbTestTempestcf(t *testing.T, tempestcfbup string) {
	if err := os.Remove(Tempestcf); err != nil {
		t.Log("[ERROR]:: An error occurred when trying to remove the test tempestcf:", Tempestcf)
		t.Fail()
	}
	Tempestcf = tempestcfbup
}

// createTestDir createa directory meant for testing
func createTestDir(path string) error {
	return os.MkdirAll(path, 0777)
}

// addToTestDir add items to the path (test dir)
// ex:
// ```go
// {
// 	0: {
// 		"dir1",
// 		"dir2",
// 	},
// 	1: {
// 		"file1",
// 		"file2",
// 	}
// }
// ```
func addToTestDir(pathTestDir string, toAdd [][]string) error {
	const DIR = 0
	const FILE = 1

	// Create dirs
	if len(toAdd[DIR]) > 0 {
		for _, path := range toAdd[DIR] {
			if err := os.MkdirAll(path, 0777); err != nil {
				return err
			}
		}
	}

	if len(toAdd) >= 2 {
		for _, path := range toAdd[FILE] {
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			f.Close()
		}
	}

	return nil
}

// rmTestDirs removes the test directory(ies) and everything in it
func rmTestDirs(dirs ...string) error {
	for _, pathDir := range dirs {
		if err := os.RemoveAll(pathDir); err != nil {
			return err
		}
	}
	return nil
}
