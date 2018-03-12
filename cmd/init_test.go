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
	// Save current Tempestcf
	tempestymlOld := Tempestyml
	// Set the temporary .tempestcf used for the test
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

	// Restore previous Tempestcf
	*tempest = tempestOld

	return
}
