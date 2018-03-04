package cmd

import (
	"errors"
	"os"
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
		if got != tst.want {
			t.Log(tst.err.Error())
			t.Fail()
		}
	}
}

// TestTreatLastChar is the test func for TreatLastChar.
// It checks if it does strip only the last character
// if it is a path separator character and does nothing otherwise.
func TestTreatLastChar(t *testing.T) {
	// test variables
	p1 := string(os.PathSeparator) + "path1" + string(os.PathSeparator) + "sub" + string(os.PathSeparator) + "dir"
	w1 := p1

	p2 := string(os.PathSeparator) + "path1" + string(os.PathSeparator) + "sub" + string(os.PathSeparator) + "dir" + string(os.PathSeparator)
	w2 := p1

	// tests holds the tests we want to do and the result expected
	var tests = []struct {
		param string
		want  string
		err   error
	}{
		{p1, w1, errors.New("[CONFUSION]:: The path was correct damn it")},
		{p2, w2, errors.New("[FAIL]:: Did not change when it was supposed to")},
	}

	// runing tests
	for _, tst := range tests {
		got := TreatLastChar(tst.param)
		// fmt.Println(got) // DEBUG
		if got != tst.want {
			t.Log(tst.err.Error())
			t.Fail()
		}
	}
}
