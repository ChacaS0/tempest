package cmd

import (
	"fmt"
	"testing"

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
