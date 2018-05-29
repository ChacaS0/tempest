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

	// running tests
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

	// running tests
	for _, tst := range tests {
		got := TreatLastChar(tst.param)
		// fmt.Println(got) // DEBUG
		if got != tst.want {
			t.Log(tst.err.Error())
			t.Fail()
		}
	}
}

// TestTreatRelativePath is the test for treatRelativePath(path *string, workDir string) {}.
// Should check if it uses full path to the working dir when adding relative paths.
// It should do nothing for full paths.
func TestTreatRelativePath(t *testing.T) {
	// working directory
	workDir, errDir := os.Getwd()
	if errDir != nil {
		t.Log("[ERROR]:: Can't retrieve the current directory", errDir)
		t.FailNow()
	}
	// params
	p1 := "./Downloads/temp"
	p2 := "Documents/temp"
	p3 := conf.Gopath
	// wants
	w1 := workDir + "/Downloads/temp"
	w2 := workDir + "/" + p2

	var tests = []struct {
		path  string
		wkdir string
		want  string
		errT  error
	}{
		{p1, workDir, w1, errors.New("[FAIL]:: Can't process './' types")},
		{p2, workDir, w2, errors.New("[FAIL]:: Can't process 'Doc/temp... types")},
		{p3, workDir, p3, errors.New("[FAIL]:: Should have done nothing, yet did much")},
	}

	for _, tst := range tests {
		treatRelativePath(&tst.path, tst.wkdir)
		if tst.path != tst.want {
			t.Log(tst.errT, "\n\t[GOT]-> ", tst.path, "\n\t[WANT]->", tst.want)
			t.Fail()
		}
	}
}

// TestAddLine is the test for addLine(args []string) error {}
// Check if it does add the proper line to a .tempestcf file.
func TestAddLine(t *testing.T) {
	// args to add
	args := []string{
		conf.Gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "github.com" + string(os.PathSeparator) + "ChacaS0" + string(os.PathSeparator) + "tempest" + string(os.PathSeparator) + "vendor",
		conf.Gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "github.com" + string(os.PathSeparator) + "ChacaS0" + string(os.PathSeparator) + "tempest" + string(os.PathSeparator) + "cmd",
		conf.Gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "github.com" + string(os.PathSeparator) + "ChacaS0" + string(os.PathSeparator) + "tempest",
	}

	// Setup test presets?
	tempestcfbup := setTestTempestcf(t, args)

	// Get all paths added in Tempestcf
	allPaths, errAllP := getPaths()
	if errAllP != nil {
		errDel := os.Remove(Tempestcf)
		if errDel != nil {
			t.Log(errDel)
		}
		t.Error(errAllP)
	}

	var cpt int
	for _, onePath := range allPaths {
		for _, oneArg := range args {
			if onePath == oneArg {
				cpt++
			}
		}
	}
	if cpt != len(args) {
		errDel := os.Remove(Tempestcf)
		if errDel != nil {
			t.Log(errDel)
		}
		t.Log("[FAIL]:: Did not add all paths to ", Tempestcf)
		t.Fail()
	}

	// Clean up the mess
	fbTestTempestcf(t, tempestcfbup)
}

// TestFindDirs
func TestFindDirs(t *testing.T) {
	// init vars
	rootPath := conf.Gopath + Slash + "tests_temp"
	testDirs := []string{
		0: conf.Gopath + Slash + "tests_temp" + Slash + "dir1",
		1: conf.Gopath + Slash + "tests_temp" + Slash + "test",
		2: conf.Gopath + Slash + "tests_temp" + Slash + "dir3" + Slash + "test",
		3: conf.Gopath + Slash + "tests_temp" + Slash + "dir4" + Slash + "notTest",
	}
	// setup test envi
	if err := createTestDir(rootPath); err != nil {
		t.Log("[ERROR]:: Could not create the test dir\n\t->", err)
		t.FailNow()
	}

	param2 := [][]string{
		0: testDirs,
	}

	if err := addToTestDir(rootPath, param2); err != nil {
		t.Log("[ERROR]:: Could not create the test content\n\t->", err)
		rmTestDirs(rootPath)
		t.FailNow()
	}

	// mock args
	args := []struct {
		root    string
		pattern string
		want    []string
	}{
		{rootPath, "test", []string{
			testDirs[1],
			testDirs[2],
		}},
	}

	for _, tst := range args {
		slFound, errFound := findDirs(tst.root, tst.pattern)
		if errFound != nil {
			t.Log("[ERROR]:: Searching for matching directories resulted to an error:\n\t-> ", errFound)
			t.Fail()
		}
		if !SameSliceValuesStr(tst.want, slFound) {
			t.Log("[FAILED]:: Wanted:\n\t", tst.want, "\nGot:\n\t", slFound)
			t.Fail()
		}
	}

	rmTestDirs(rootPath)

}
