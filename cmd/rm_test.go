package cmd

import (
	"errors"
	"fmt"
	"testing"
)

// TestRmInSlice checks rmInSlice.
// If the slice returned doesn't contain
// the string passed as param and contains all the other strings
// of the slice passed as param
func TestRmInSlice(t *testing.T) {
	sl1 := []string{0: "path1/sub1/subsub1", 1: "path2/sub2/subsub2", 2: "path3/sub3/subsub3"}
	sl2 := []string{0: "path1/sub1"}
	//---
	sl22 := []string{}
	sl11 := []string{0: "path2/sub2/subsub2", 1: "path3/sub3/subsub3"}

	var tests = []struct {
		i     int
		str   string
		slstr []string
		want  []string
		e     error
	}{
		{-1, "path1/sub1/subsub1", sl1, sl11, errors.New("[FAIL]:: Didn't remove shit - str")},
		{-1, "path1/sub1", sl2, sl22, errors.New("[FAIL]:: Can't return the nil slice, sad - str")},
		{0, "", sl1, sl11, errors.New("[FAIL]:: Didn't remove shit - int")},
		{0, "", sl2, sl22, errors.New("[FAIL]:: Can't return the nil slice, sad - int")},
	}

	// runing tests
	for _, tst := range tests {
		got := rmInSlice(tst.i, tst.str, tst.slstr)
		if !SameSlices(got, tst.want) {
			fmt.Println("got:", got, "\nwant:", tst.want)
			t.Log(tst.e.Error())
			t.Fail()
		}
	}
}

// SameSlices checks equality between two slices
// returns true if they are identiques
func SameSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
