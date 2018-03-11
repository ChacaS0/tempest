package cmd

import (
	"errors"
	"fmt"
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

// TestBackupTempestcf checks if the backup is really done
// by backupTempestcf
func TestBackupTempestcf(t *testing.T) {
	// Save old Tempestcf
	tempestcfBup := Tempestcf
	// Set the new Tempestcf for testing
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	// Create first the temporary new .tempestcf as $GOPATH/.tempestcf.temp
	f, err := os.OpenFile(Tempestcf, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
	}

	// Add a few things to it
	if err := addLine(sl1); err != nil {
		t.Log("[ERROR]:: Failed to add lines to", Tempestcf, "\n\t->", err)
		if errDel := os.Remove(Tempestcf); errDel != nil {
			t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
		}
		Tempestcf = tempestcfBup
		t.FailNow()
	}

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
	if errDel := os.Remove(Tempestcf); errDel != nil {
		t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
	}
	Tempestcf = tempestcfBup
}

// TestRestoreTempestcf check if it restores well a bup of .tempestcf
// with restoreTempestcf() error {}
func TestRestoreTempestcf(t *testing.T) {
	// Save old Tempestcf
	tempestcfBup := Tempestcf
	// Set the new Tempestcf for testing
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	// Create first the temporary new .tempestcf as $GOPATH/.tempestcf.temp
	f, err := os.OpenFile(Tempestcf, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
	}

	// Add a few things to it
	if err := addLine(sl1); err != nil {
		t.Log("[ERROR]:: Failed to add lines to", Tempestcf, "\n\t->", err)
		if errDel := os.Remove(Tempestcf); errDel != nil {
			t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
		}
		Tempestcf = tempestcfBup
		t.FailNow()
	}

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
	if errDel := os.Remove(Tempestcf); errDel != nil {
		t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
	}
	Tempestcf = tempestcfBup
}

// TestWriteTempestcf checks if new data is well written to .tempestcf
// with ``writeTempestcf(targets []string) error {}``.
// It is supposed to override the .tempestcf targets with a new slice of targets.
func TestWriteTempestcf(t *testing.T) {
	// Save old Tempestcf
	tempestcfBup := Tempestcf
	// Set the new Tempestcf for testing
	Tempestcf = conf.Gopath + string(os.PathSeparator) + ".tempestcf"

	// Create first the temporary new .tempestcf as $GOPATH/.tempestcf.temp
	f, err := os.OpenFile(Tempestcf, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	sl1 := []string{
		conf.Gobin,
		conf.Gopath,
	}

	// Add a few things to it
	if err := addLine(sl1); err != nil {
		t.Log("[ERROR]:: Failed to add lines to", Tempestcf, "\n\t->", err)
		if errDel := os.Remove(Tempestcf); errDel != nil {
			t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
		}
		Tempestcf = tempestcfBup
		t.FailNow()
	}

	newSl := rmInSlice(0, "", sl1)

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
		t.Log("[ERROR]:: Could not delete the .tempestcf.temp")
	}
	Tempestcf = tempestcfBup
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
