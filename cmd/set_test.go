package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// func Test(t *testing.T)

// TestSetAge checks if when changing the age,
// it is stored correctly in viper.Config()
func TestSetAge(t *testing.T) {
	// Current settings
	currAge := viper.GetInt("duration")
	currAuto := viper.GetBool("auto-mode")
	currCfFile := viper.ConfigFileUsed()

	// set env for testing
	_, errDir := IsDirectory(conf.Gopath + "/.tempest_test.yaml")
	if errDir != nil {
		// if already exists, we create
		f, err := os.OpenFile(conf.Gopath+"/.tempest_test.yaml", os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(color.HiRedString("::Huge error! Could not recreate file! God lost faith in you!"))
			t.Log("[ERROR]File creation::", err)
			t.FailNow()
		}
		defConf := `duration: 5
auto-mode: false
`
		_, errWrite := f.WriteString(defConf)
		if errWrite != nil {
			t.Log("Could not write default config to test file: ", errWrite)
			t.FailNow()
		}
		defer f.Close()
		viper.SetConfigFile(conf.Gopath + "/.tempest_test.yaml")
	}

	//* First test - upgrade age
	if currAge > 1 {
		age = currAge + 5
		fmt.Println("AGE:", age)
	} else {
		age = 6
	}
	// Set the age
	_ = setAge()
	// Check if it changed
	if age != viper.GetInt("duration") && viper.GetInt("duration") != 6 {
		t.Log("[CHANGE]:: Couldn't set the age")
		t.Fail()
	}
	// Check if auto-mode didn't get ereased or affected
	if viper.GetBool("auto-mode") != currAuto {
		t.Log("[AFFECTED]:: auto-mode changed and wasn't supposed to")
		t.Fail()
	}

	// Go back to the previous configuration
	// Age
	viper.Set("duration", currAge)
	errS := viper.WriteConfigAs(viper.ConfigFileUsed())
	if errS != nil {
		color.Red(errS.Error())
	}
	// ConfigFile
	// clean the one just created
	errRm := os.Remove(conf.Gopath + "/.tempest_test.yaml")
	if errRm != nil {
		t.Log("Could not remove the test file, you might one to remove it yourself!")
		t.Log(errRm)
	}
	viper.SetConfigFile(currCfFile)

	// All done
}

// TestSetAutoStart tests if it really activates auto start (depends on the device OS)
func TestSetAutoStart(t *testing.T) {
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

	//* Linux
	// For linux it should create the file ~/.config/autostart/tempest.desktop
	// Should be already "off", meaning ``false``.
	// Check if no args messes up
	if err := setAutoStart(); err != nil {
		t.Log("[ERROR]::  An error was encoutered while using setAutostart:\n\t->", err)
	}

	// Now set it to off
	autoStart = "off"
	if err := setAutoStart(); err != nil {
		t.Log("[ERROR]::  An error was encoutered while using setAutostart:\n\t->", err)
	}

	// And then to on
	autoStart = "on"
	if err := setAutoStart(); err != nil {
		t.Log("[ERROR]::  An error was encoutered while using setAutostart:\n\t->", err)
	}

	// TODO Need windows mode to be done too :/
	// Maybe try to change the runtime.GOOS ? -- will need to set the env and paths separators for windows !?

	// Set back to default conf and clean temp test files
	if err = cleanTempest(t, &Tempestyml, tempestymlOld); err != nil {
		t.Log("[ERROR]:: An error occurred:", err)
		t.Fail()
	}
	viper.SetConfigFile(Tempestyml)
}
