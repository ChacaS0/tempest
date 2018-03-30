// Copyright © 2018 Sebastien Bastide
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// @DEPRECATED
// rmInt is the index to remove from the TEMPest list
// var rmInt int

// @DEPRECATED
// rmStr is the path to be removed from TEMpest
// var rmStr string

// rmOrigin defines whether the origin directory/file should be deleted too
var rmOrigin bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Used to remove a path from the tracking list of TEMPest",
	Long: `Used to remove a path from the tracking list of TEMPest.
The basic command does not delete the original directory or file.
It's just telling to TEMPest to untrack the file/directory. For example:

	tempest rm -p /tmp

or
Place yourself inside the directory and run:

	tempest rm

If it is a directory tracked by TEMPest, it will be untracked, otherwise you will get an error.

It is also possible to use the index number resulting of the tempest list command:

	tempest rm -i 1

` + color.RedString("/!\\ [WARNING] After removing a file, the index number might change!!") + `

In order to remove the original file/directory:

	tempest rm -o 1
or
	tempest rm --origin 1

=> Considering 1 is /tmp, this would remove /tmp from TEMPest AND from your device.

`,
	Run: func(cmd *cobra.Command, args []string) {
		// color.Cyan("Sorry, rm is not fully implemented yet... Coming soon don't worry!")
		if len(args) == 0 {
			args = append(args, "this")
		}

		// handles rm's args
		// slRmInt, slRmStr := processArgsRm(args)

		slicePaths, errAllP := getPaths()
		if errAllP != nil {
			color.Red(errAllP.Error())
		}
		// slicePaths = rmInSlice(rmInt, rmStr, slicePaths) // @DEPRECATED
		slRmInt, slRmStr := processArgsRm(args)
		slicePaths = rmInSlice(slRmInt, slRmStr, slicePaths)

		if errWrite := writeTempestcf(slicePaths); errWrite != nil {
			color.Red(errWrite.Error())
		}

		if errLi := printList(); errLi != nil {
			color.Red(errLi.Error())
		}
		// fmt.Println("rm called")
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rmCmd.Flags().IntVarP(&rmInt, "index", "i", -1, "Points to an index provided by TEMPest")
	// rmCmd.Flags().StringVarP(&rmStr, "path", "p", "", "The path of a target for TEMPest")

	rmCmd.Flags().BoolVarP(&rmOrigin, "origin", "o", false, "Removes the target from TEMPest, but also the original directories/files")
}

// processArgsRm takes the args as parameters and process them
func processArgsRm(args []string) ([]int, []string) {
	slRmInt := make([]int, 0)
	slRmStr := make([]string, 0)

	// // regexes
	// // RegexManyInt is the regex to see if many int were passed
	// // i.e. ``0 1 4`` or ``0``
	// // RegexManyInt := regexp.MustCompile(`([0-9]+\s|[0-9]|^[0-9]+$)`)
	// RegexIntToInt is the regex to see if we want a range of int
	// i.e. ``0-2`` (would be 0 1 2)
	RegexIntToInt := regexp.MustCompile(`(\d+-\d+)`)
	// RegexManyStr is the regex to see if many strings (targets) were passed
	// i.e. ``/tmp`` /path1/subpath1``
	RegexManyStr := regexp.MustCompile(`^(\/|\\)(\d|\D)+`)
	// RegexJustInt is the regex to see if only an int has been passed as an arg
	RegexJustInt := regexp.MustCompile(`^(\d+)$`)

	for _, arg := range args { // for each arg in args
		arg = strings.Trim(arg, " ")
		// Wildcard: * (ALL)
		if arg == "*" && len(args) == 1 {
			slRmStr = append(slRmStr, arg)
		}
		// Empty arg
		if arg == "this" {
			if !IsStringInSlice("this", slRmStr) {
				slRmStr = append(slRmStr, "this")
			}
		}
		// IntToInt
		if allItoI := RegexIntToInt.FindString(arg); len(allItoI) > 0 {
			// RegexIntToInt.ReplaceAllString(arg, ``)
			// explode with ``-`` and get the left and right value
			values := strings.Split(arg, "-")
			begin, errBegin := strconv.Atoi(values[0])
			end, errEnd := strconv.Atoi(values[1])
			if errBegin != nil || errEnd != nil {
				fmt.Println(redB(":: [ERROR]"), color.HiRedString("Sorry, could not understand those arguments", arg), "\n\t[0]:", values[0], "\n\t->", errBegin, "\n\t[1]:", values[1], "\n\t->", errEnd)
				return nil, nil
			}
			if begin > end {
				temp := end
				end = begin
				begin = temp
			}
			for i := begin; i <= end; i++ {
				if !IsIntInSlice(i, slRmInt) {
					slRmInt = append(slRmInt, i)
				}
			}
		}
		// Many Strings
		if allMStr := RegexManyStr.FindString(arg); len(allMStr) > 0 {
			if !IsStringInSlice(arg, slRmStr) {
				slRmStr = append(slRmStr, arg)
			}
		}
		// Just an Int
		if allInt := RegexJustInt.FindString(arg); len(allInt) > 0 {
			val, errInt := strconv.Atoi(arg)
			if errInt != nil {
				fmt.Println(redB(":: [ERROR]"), color.HiRedString("Sorry, could not understand thise argument", arg), "\n\t->", errInt)
				return nil, nil
			}
			if !IsIntInSlice(val, slRmInt) {
				slRmInt = append(slRmInt, val)
			}
		}

	}

	return slRmInt, slRmStr
}

// rmInSlice give it a slice with an index of path to remove.
// It will return a new slice with the item removed from it!
func rmInSlice(indexes []int, slRmStr []string, list []string) []string {

	// NOTE: alling this func, check if len(args) == 0, set record at "this"

	// for strings
	for _, record := range slRmStr {
		//* we use the the path provided
		if record == "this" {
			this, errDir := os.Getwd()
			if errDir != nil {
				color.Red(errDir.Error())
			}
			record = this
		}

		for i, v := range list {
			if v == record {
				indexes = append(indexes, i)
			}
		}

		if record == "*" {
			return []string{}
		}
	}

	// Slicing an element out of the slice.
	// Example to take of the element at position 1:
	// 	a = a[:1+copy(a[1:], a[2:])]
	// listToRet := list[:index+copy(list[index:], list[index+1:])]

	//* if we have to, remove the directories/files
	if rmOrigin {
		// need to find another way QQ
		trashSlice := make([]string, 0)
		for _, index := range indexes {
			trashSlice = append(trashSlice, list[index])
		}

		if err := simpleDelAllString(trashSlice...); err != nil {
			log.Println(redB(":: [ERROR]"), color.HiRedString("There was an error while delete original files:"), "\n\t->", err)
		}
	}

	//* Remove from slice if has stuff to remove
	listToRet := make([]string, 0)
	if len(indexes) > 0 {
		for i, item := range list {
			if !IsIntInSlice(i, indexes) {
				listToRet = append(listToRet, item)
			}
		}
	}

	// DEBUG
	// color.HiYellow(fmt.Sprintf("%v", list))

	// just for linter ...
	return listToRet
}

// simpleDelAllString is a simple function to delete files or directories.
func simpleDelAllString(paths ...string) error {
	if paths != nil {
		for _, path := range paths {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		}
	} else {
		return errors.New("Parameter is nil at rm.go:simpleDelAllString()")
	}
	return nil
}

// writeTempestcf saves the new list of paths meant to be targets for TEMPest
func writeTempestcf(targets []string) error {
	// Backup of the old one
	if errBp := backupTempestcf(); errBp != nil {
		return errBp
	}

	// init a new .tempestcf file
	if errInit := initializeTP(); errInit != nil {
		return errInit
	}

	// Deletes the current .tempestcf
	/* if errDel := os.Remove(conf.Home + "/.tempestcf"); errDel != nil {
		return errDel
	} */

	// Open the file to write to
	tmpcf, err := os.OpenFile(Tempestcf, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(redB("::"), color.RedString("Sorry! Could not find "+Tempestcf+" ! :("))
		return err
	}
	defer tmpcf.Close()

	// Add line to .tempestcf
	for _, path := range targets {
		// write each line
		// Add this directory to the file
		_, errWrite := tmpcf.WriteString(path + "\n")
		if errWrite != nil {
			if errRest := restoreTempestcf(); errRest != nil {
				return errRest
			}
			fmt.Println(redB("::"), color.RedString("Could not write to the file. Fail bitch!"))
			return errWrite
		}
	}

	info := greenB(":: ")
	info += color.HiGreenString("Removed with success!!")
	fmt.Println(info)

	fmt.Println(yellowB("::"), color.HiYellowString("Check out the new INDEXES, MUCH FUNNNNN ! !\n"))

	return nil
}

// backupTempestcf save the current .tempestcf as .tempestcf.old
func backupTempestcf() error {
	errBup := os.Rename(Tempestcf, Tempestcf+".old")
	if errBup != nil {
		return errBup
	}
	return nil
}

// retoreTempestcf bring back the previous .tempestcf
func restoreTempestcf() error {
	errRestore := os.Rename(Tempestcf+".old", Tempestcf)
	if errRestore != nil {
		return errRestore
	}
	return nil
}
