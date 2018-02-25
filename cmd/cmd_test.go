package cmd

import "testing"

func getPathsTest(t *testing.T) {
	_, err := getPaths()
	if err != nil {
		t.Error("getPaths() could not read from the file .tempestcf and returned with this error:\n", err.Error())
	}
}
