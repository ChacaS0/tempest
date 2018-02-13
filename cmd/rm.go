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
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rmInt is the index to remove from the TEMPest list
var rmInt int

// rmStr is the path to be removed from TEMpest
var rmStr string

// rmOrigin defines whether the origin directory/file should be deleted too
var rmOrigin bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Used to remove a path from the tracking list of TEMPest",
	Long: `Used to remove a path from the tracking list of TEMPest.
The basic command does not delete the original directory or file.
It's just telling to TEMPest to untrack the file/directory. For example:

	tempest rm /tmp

or
Place yourself inside the directory and run:

	tempest rm

If it is a directory tracked by TEMPest, it will be untracked, otherwise you will get an error.

[SOON] It is also possible to use the index number resulting of the tempest list command:

	tempest rm 1

` + color.RedString("/!\\ [WARNING] After removing a file, the index number might change!!") + `

In order to remove the original file/directory:

	tempest rm -o 1
or
	tempest rm --origin 1

=> Considering 1 is /tmp, this would remove /tmp from TEMPest AND from your device.

`,
	Run: func(cmd *cobra.Command, args []string) {
		// color.Cyan("Sorry, rm is not fully implemented yet... Coming soon don't worry!")
		if len(args) == 0 && rmInt == -1 && rmStr == "" {
			rmStr = "this"
		}
		slicePaths, errAllP := getPaths()
		if errAllP != nil {
			color.Red(errAllP.Error())
		}
		slicePaths = rmInSlice(rmInt, rmStr, slicePaths)

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
	rmCmd.Flags().IntVarP(&rmInt, "index", "i", -1, "Points to an index provided by TEMPest")
	rmCmd.Flags().StringVarP(&rmStr, "path", "p", "", "The path of a target for TEMPest")

	rmCmd.Flags().BoolVarP(&rmOrigin, "origin", "o", false, "Removes the target from TEMPest, but also the original directories/files")
}

// rmInSlice give it a slice with an index of path to remove.
// It will return a new slice with the item removed from it!
func rmInSlice(index int, record string, list []string) []string {
	// alling this func, check if len(args) == 0, set record at "this"

	if index == -1 {
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
				index = i
			}
		}
	}

	// Slicing an element out of the slice.
	// Exmaple to take of the element at position 1:
	// 	a = a[:1+copy(a[1:], a[2:])]
	list = list[:index+copy(list[index:], list[index+1:])]

	// DEBUG
	// color.HiYellow(fmt.Sprintf("%v", list))

	return list
}

// xxx saves the new list of paths meant to be targets for TEMPest
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
	tmpcf, err := os.OpenFile(conf.Home+"/.tempestcf", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(redB("::"), color.RedString("Sorry! Could not find ~/.tempestcf ! :("))
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
	info += color.HiGreenString("Removed with success!!\n")
	fmt.Println(info)

	return nil
}

// backupTempestcf save the current .tempestcf as .tempestcf.old
func backupTempestcf() error {
	errBup := os.Rename(conf.Home+"/.tempestcf", conf.Home+"/.tempestcf.old")
	if errBup != nil {
		return errBup
	}
	return nil
}

// retoreTempestcf bring back the previous .tempestcf
func restoreTempestcf() error {
	errRestore := os.Rename(conf.Home+"/.tempestcf.old", conf.Home+"/.tempestcf")
	if errRestore != nil {
		return errRestore
	}
	return nil
}
