package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
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
	sl11 := []string{0: "path2/sub2/subsub2", 1: "path3/sub3/subsub3"}
	// empty
	emptySlInt := make([]int, 0)
	emptySlStr := make([]string, 0)
	// many ints
	sl3 := []int{0, 2}
	// this
	this, errDir := os.Getwd()
	if errDir != nil {
		t.Log(errDir.Error())
	}
	slThis := []string{this, "/tmp"}
	// all (*)
	slAll := []string{"*"}

	var tests = []struct {
		i     []int
		str   []string
		slstr []string
		want  []string
		e     error
	}{
		{emptySlInt, []string{"path1/sub1/subsub1"}, sl1, sl11, errors.New("[FAIL]:: Didn't remove shit - str")},
		{emptySlInt, sl2, sl2, emptySlStr, errors.New("[FAIL]:: Can't return the nil slice, sad - str")},
		{[]int{0}, emptySlStr, sl1, sl11, errors.New("[FAIL]:: Didn't remove shit - int")},
		{[]int{0}, emptySlStr, sl2, emptySlStr, errors.New("[FAIL]:: Can't return the nil slice, sad - int")},
		{[]int{0, 1}, emptySlStr, sl1, []string{"path3/sub3/subsub3"}, errors.New("[FAIL]:: This trash sh*t can't remove with multiple indexes")},
		// many strings
		{emptySlInt, sl11, sl1, []string{"path1/sub1/subsub1"}, errors.New("[FAIL]:: Failed to remove two paths at once. Noob")},
		// many ints
		{sl3, emptySlStr, sl1, []string{"path2/sub2/subsub2"}, errors.New("[FAIL]:: Failed to remove multiple ints at once")},
		// this
		{emptySlInt, []string{"this"}, slThis, []string{"/tmp"}, errors.New("[FAIL]:: Could not remove 'this': " + this)},
		// all (*)
		{emptySlInt, slAll, sl1, emptySlStr, errors.New("[FAIL]:: Failed to use the all-wildcard")},
	}

	// running tests
	for _, tst := range tests {
		got := rmInSlice(tst.i, tst.str, tst.slstr)
		if !SameSlices(got, tst.want) {
			fmt.Println("got:", got, "\nwant:", tst.want)
			t.Log(tst.e.Error())
			t.Fail()
		}
	}
}

// TestBackupTempestcf checks if the backup is really done
// by backupTempestcf
func TestBackupTempestcf(t *testing.T) {

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
	}

	tempestcfbup := setTestTempestcf(t, sl1)

	// Try to make a backup
	if err := backupTempestcf(); err != nil {
		t.Log("[FAIL]:: Could not back up the file:", Tempestcf, "\n\t->", err)
		t.Fail()
	}

	// Check if the bup still has the previous content (sl1)
	Tempestcf += ".old"
	sl2, errPaths := getPaths()
	if errPaths != nil {
		t.Log("[ERROR]:: Could not get the paths of the .tempestcf.temp:", Tempestcf, "\n\t->", errPaths)
		t.Fail()
	}
	if !SameSlices(sl1, sl2) {
		t.Log("[FAIL]:: Not the same values for the two files, backup changed data!")
		t.Fail()
	}

	// In the end we restore the previous Tempestcf and delete the .tempestcf.temp just in case
	fbTestTempestcf(t, tempestcfbup)
}

// TestRestoreTempestcf check if it restores well a bup of .tempestcf
// with restoreTempestcf() error {}
func TestRestoreTempestcf(t *testing.T) {

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
	}

	// Testing presets
	tempestcfbup := setTestTempestcf(t, sl1)

	// Try to make a backup
	if err := backupTempestcf(); err != nil {
		t.Log("[FAIL]:: Could not back up the file:", Tempestcf, "\n\t->", err)
		t.Fail()
	}

	// Now try to restore the previous file
	if err := restoreTempestcf(); err != nil {
		t.Log("[FAIL]:: Failed while restoring .tempestcf. Shame!\n\t->", err)
		Tempestcf += ".old"
		t.Fail()
	}

	// Then we restore the previous setup!?
	fbTestTempestcf(t, tempestcfbup)
}

// TestWriteTempestcf checks if new data is well written to .tempestcf
// with ``writeTempestcf(targets []string) error {}``.
// It is supposed to override the .tempestcf targets with a new slice of targets.
func TestWriteTempestcf(t *testing.T) {

	this, errDir := os.Getwd()
	if errDir != nil {
		log.Fatal(errDir)
	}

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
		this,
	}

	// Presets for testing
	tempestcfbup := setTestTempestcf(t, sl1)

	newSl := rmInSlice([]int{0}, []string{}, sl1)
	newSl = rmInSlice([]int{}, []string{"this"}, sl1)

	// Try to use writeTempestcf
	if err := writeTempestcf(newSl); err != nil {
		t.Log("[ERROR]:: Error while using writeTempestcf(newSl)\n\t->", err)
		t.Fail()
	}

	// Verifications
	actualSl, errPaths := getPaths()
	if errPaths != nil {
		t.Log("[ERROR]:: Could not read the new paths!!\n\t->", errPaths)
		t.Fail()
	}
	if !SameSlices(actualSl, newSl) {
		t.Log("[FAIL]:: After writing .tempestcf, data do not match you piece of ****!\n\t->")
		t.Fail()
	}

	// Then we restore the previous setup!?
	if errDel := os.Remove(Tempestcf); errDel != nil {
		t.Log("[ERROR]:: Did not delete the .tempestcf")
	}
	Tempestcf += ".old"
	if errDel := os.Remove(Tempestcf); errDel != nil {
		t.Log("[ERROR]:: Did not delete the .tempestcf.old")
	}
	Tempestcf = tempestcfbup
}

// TestProcessArgsRm is the test for processArgsRm.
// It verifies that it returns the right slices
func TestProcessArgsRm(t *testing.T) {
	// parameter variables
	// this arg
	s3 := []string{"this"}
	// index ranges
	s1 := []string{"0-1"}
	s2 := []string{"0-2", "5-7"}
	s12 := []string{"1-2", "5-3", "7-8"}
	// Just targets
	s4 := []string{"/temp", "/tmp/user/aur"}
	s5 := []string{"/tmp"}
	// Just Ints
	s6 := []string{"0", "2", "5"}
	s7 := []string{"1"}
	// mixed
	s8 := []string{"0", "2-5", "7", "/tmp"}
	s9 := []string{"/tmp", "7", "2-5"}
	s10 := []string{"2-5", "/tmp", "7"}
	s11 := []string{"1-2", "3-4", "/temp", "5", "/tmp/user/aur", "6"}
	s14 := []string{"0", "this", "3-3", "/tmp"}
	// Overlaping arguments
	s13 := []string{"0-2", "1-3", "/tmp", "10", "/tmp"}

	emptySliceStr := make([]string, 0)
	emptySliceInt := make([]int, 0)

	// tests holds the tests we want to do and the result expected
	var tests = []struct {
		param     []string
		wantSlInt []int
		wantSlStr []string
		err       error
	}{
		// Empty
		{s3, emptySliceInt, []string{"this"}, errors.New("[FAIL]:: Failed to prcess empty args")},
		// IntToInt
		{s1, []int{0, 1}, emptySliceStr, errors.New("[FAIL]:: Failed to return the slice of ints")},
		{s2, []int{0, 1, 2, 5, 6, 7}, emptySliceStr, errors.New("[FAIL]:: Failed to process 2 ranges of ints")},
		{s12, []int{1, 2, 3, 4, 5, 7, 8}, emptySliceStr, errors.New("[FAIL]:: Could not treat some stuff of the for 10-4")},
		// Many strings
		{s4, emptySliceInt, s4, errors.New("[FAIL]:: Cannot process many targets")},
		// One string
		{s5, emptySliceInt, s5, errors.New("[FAIL]:: Cannot process a single target")},
		// Many ints
		{s6, []int{0, 2, 5}, emptySliceStr, errors.New("[FAIL]:: Failed to process a bunch of ints")},
		{s7, []int{1}, emptySliceStr, errors.New("[FAIL]:: Could not process just one fcking int ffs")},
		// Overlaping arguments
		{s13, []int{0, 1, 2, 3, 10}, s5, errors.New("[FAIL]:: Failed to handle redundant information")},
		// Wildcard * (ALL)
		{[]string{"*"}, emptySliceInt, []string{"*"}, errors.New("[FAIL]:: Failed to handle wildcard: *")},
		// Mixed
		{s8, []int{0, 2, 3, 4, 5, 7}, s5, errors.New("[FAIL]:: Failed to process many mixed args (1)")},
		{s9, []int{7, 2, 3, 4, 5}, s5, errors.New("[FAIL]:: Failed to process many mixed args (2)")},
		{s10, []int{2, 3, 4, 5, 7}, s5, errors.New("[FAIL]:: Failed to process many mixed args (3)")},
		{s11, []int{1, 2, 3, 4, 5, 6}, s4, errors.New("[FAIL]:: Failed to process many mixed args (4)")},
		{s14, []int{0, 3}, []string{"this", "/tmp"}, errors.New("[FAIL]:: Failed to process many mixed args (5)")},
	}

	// running tests
	for _, tst := range tests {
		gotSlInt, gotSlStr := processArgsRm(tst.param)
		if !SameSlicesInt(gotSlInt, tst.wantSlInt) || !SameSlices(gotSlStr, tst.wantSlStr) {
			fmt.Println("[GOT] -->\t", gotSlInt, gotSlStr)
			fmt.Println("[WANT]-->\t", tst.wantSlInt, tst.wantSlStr)
			fmt.Printf(":GOT:\t\t %T %T\n", gotSlInt, gotSlStr)
			fmt.Printf(":WANT:\t\t %T %T\n", tst.wantSlInt, tst.wantSlStr)
			fmt.Println("-------------------------------------------------")
			t.Log(tst.err.Error())
			t.Fail()
		}
	}
}

// TestSimpleDelAllString is the test for the func simpleDelAllString which is supposed
// to delete from the system all paths provided in args
func TestSimpleDelAllString(t *testing.T) {
	// create a test directories
	dir1 := conf.Gopath + string(os.PathSeparator) + "dir1"
	dir2 := conf.Gopath + string(os.PathSeparator) + "dir2"
	dir3 := conf.Gopath + string(os.PathSeparator) + "dir3"

	dirs := []string{dir1, dir2, dir3}

	for _, d := range dirs {
		if err := os.Mkdir(d, 0777); err != nil {
			t.Log("[ERROR]:: Failed to create test directeory:", d)
			t.Fail()
		}
	}

	// apply func on them in order to delete it
	if err := simpleDelAllString(dirs...); err != nil {
		t.Log("[FAIL]:: Failed to remove those directories")
		t.Fail()
	}

	stillExisting := make([]string, 0)

	// check if they got deleted for real
	for _, d := range dirs {
		if isDir, err := IsDirectory(d); isDir || err == nil {
			stillExisting = append(stillExisting, d)
			t.Log("[FAIL]:: Think it deletes but it doesn't, this is so mean, imma cry now. Bye!")
			t.Fail()
		}
	}

	// clean up the mess if one of the test dirs still exist
	for _, se := range stillExisting {
		if err := os.RemoveAll(se); err != nil {
			t.Log("[ERROR]:: F*ck can't do sh*t.")
			t.Fail()
		}
	}
}
