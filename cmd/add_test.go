package cmd

import (
	"errors"
	"fmt"
	"testing"
)

// TestCheckRedondance check if it does check right redundance !
func TestCheckRedondance(t *testing.T) {
	// parameter variables declarations
	s1 := make([]string, 0)
	s2 := make([]string, 0)
	s3 := make([]string, 0)
	s4 := make([]string, 0)
	s5 := make([]string, 0)

	// Feed the parameter variables
	s1 = append(s1, "/path/1")
	s1 = append(s1, "/path/2")
	s1 = append(s1, "/path/3")
	s1 = append(s1, "/path/4")

	s2 = append(s2, "/path/8")
	s2 = append(s2, "/path/9")
	s2 = append(s2, "/path/5")

	s3 = append(s3, "/path/3")
	s3 = append(s3, "/path/4")

	s4 = append(s4, "/rob/bor")
	s4 = append(s4, "/maestro/ortseam")
	s4 = append(s4, "/rob/bor")

	s5 = append(s5, "/path/paf/paph")
	s5 = append(s5, "/paf/path/paph")
	s5 = append(s5, "/paph/path/paph")
	s5 = append(s5, "/paf/paf/paf")

	// tests holds the tests we want to do and the result expected
	var tests = []struct {
		p1   []string
		p2   []string
		want bool
		err  error
	}{
		{s1, s2, false, errors.New("[CONFUSION]:: Thinks there is an existing string but there is not (or inverse)")},
		{s1, s3, true, errors.New("[FAIL]:: Failed to see the existing string")},
		{s2, s4, true, errors.New("[FAIL]:: There are two or more same paths in the same array")},
		{s1, s5, false, errors.New("[CONFUSION]:: Thinking there is an error, but the paths name are just looking alike")},
	}

	// runing tests
	for _, tst := range tests {
		got := checkRedondance(tst.p1, tst.p2)
		fmt.Println(got)
		if got != tst.want {
			t.Log(tst.err.Error())
			t.Fail()
		}
	}
}
