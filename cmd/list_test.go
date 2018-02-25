package cmd

import (
	"testing"
)

// TestGetPaths() test function for list/getPaths()
func TestGetPaths(t *testing.T) {
	_, err := getPaths()
	if err != nil {
		// if err.Error() != "empty" {
		// 	t.Error("getPaths() could not read from the file .tempestcf and returned with this error:\n", err.Error())
		// }
	}
}
