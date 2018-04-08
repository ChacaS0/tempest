package cmd

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// TestGetAge tests if a value is returned by this func
// and if the value is a string a the right one maybe?
func TestGetAge(t *testing.T) {
	intAge := viper.GetInt("duration")
	strAge := fmt.Sprintf("%d", intAge)

	tAge := getAge()

	// Same type?
	if "string" != fmt.Sprintf("%T", tAge) {
		t.Log()
		t.Fail()
	}

	// Same value?
	if strAge != tAge {
		t.Log("[Error - getAge()]:: ")
		t.Fail()
	}

}

// TestGetAllLogs is the test for getAllLogs(args []string){}.
// It should display all the logs available
func TestGetAllLogs(t *testing.T) {
	//setup
	testAll = true
	logshutupbup, _ := setTestLogShutup(t)

	var want string

	// styling
	headerShutup := magB("===========================================  - [ShutupLogs] -  ===================================================")
	footerShutup := magB("========================================  - [EOF - ShutupLogs] -  ================================================")

	fileCtnt, errRF := ioutil.ReadFile(LogShutup)
	if errRF != nil {
		fmt.Println(redB(":: [ERROR]"), color.HiRedString("Could not read the file -_-\n\t->", LogShutup, "\n\t->", errRF))
	}

	// what we want
	want += fmt.Sprintf("%v\n", headerShutup)
	want += fmt.Sprintln(string(fileCtnt))
	want += fmt.Sprintf("%v\n", footerShutup)

	got := captureStdout(func() {
		getAllLogs([]string{})
	})

	if got != want {
		t.Log("[FAIL]:: Result is different than expected:\n\t[GOT] \n", got, "\n\t[WANT]\n", want)
		t.Fail()
	}

	// fall back to previous conf
	fbTestLogShutup(t, logshutupbup)
}
