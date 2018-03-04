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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vrischmann/envconfig"
)

var this string

// var conf struct {
// 	Home string
// }

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <DIRECTORY_PATH>",
	Short: "Adds a new temporary directory to be handled by TEMPest",
	Long: `Adds a new temporary directory to be handled by TEMPest
Or add manually to '~/.tempestcf', or use this command. 

For example, to add the current directory :
tempest add

To add another directory:
tempest add /tmp

> By convention we will name the temporary directories "temp".
  This way they will be easy to spot
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(color.HiBlueString("lol"))
		if errAddLine := addLine(args); errAddLine != nil {
			fmt.Println("::An error occured while adding path(s):\n", errAddLine)
		}

	},
}

func init() {
	if err := envconfig.Init(&conf); err != nil {
		log.Fatal(err)
	}

	RootCmd.AddCommand(addCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// addCmd.Flags().StringVarP(&this, "this", "t", "nothing", "Points to current directory")
}

// func addLine(argFlags []string, args []string) error {
func addLine(args []string) error {
	// Check if the path already exists first in tempestcf
	ctnt, errRead := ioutil.ReadFile(conf.Home + "/.tempestcf")
	if errRead != nil {
		fmt.Println(redB("::"), color.RedString("Could not read the file muthafuckkah!"))
		return errRead
	}

	// Store it as a slice of strings
	ctntStr := string(ctnt)
	ctntSlice := strings.Split(ctntStr, "\n")

	// Open the file to write to (adding new temp path)
	// tmpcf, err := os.Open(conf.Home + "/.tempestcf")
	tmpcf, err := os.OpenFile(conf.Home+"/.tempestcf", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(color.RedString("::Sorry! Could not find ~/.tempestcf ! :("))
		return err
	}
	defer tmpcf.Close()

	// if no args
	if len(args) == 0 {
		// Fetch current directory
		this, errDir := os.Getwd()
		if errDir != nil {
			log.Fatal(errDir)
		}
		// Check if already exists
		sliceThis := make([]string, 1)
		sliceThis[0] = this
		if checkRedondance(ctntSlice, sliceThis) {
			return errors.New(blueB(":: ") + color.HiBlueString("A path provided is already in TEMPest! Later b*tch"))
		}

		// Check if path points to valid target
		// if dir, errExist := IsDirectory()

		// Add this directory to the file
		nbBytes, errWrite := tmpcf.WriteString(this + "\n")
		if errWrite != nil {
			fmt.Println(redB("::"), color.RedString("Could not write to the file. Fail bitch!"))
			return errWrite
		}
		fmt.Println(greenB("[NEW TEMP]::"), color.GreenString(fmt.Sprintf("\t%d", nbBytes)+"::>"), this)
	} else {
		// Check if already exists
		if checkRedondance(ctntSlice, args) {
			return errors.New(blueB(":: ") + color.HiBlueString("A path provided is already in TEMPest! Later b*tch"))
		}

		for _, path := range args {
			if _, errIsDir := IsDirectory(path); errIsDir != nil {
				fmt.Println(redB("::"), color.RedString("Can't add this target, check it does exist >>"), path)
				return errors.New("")
			}
			// Add all the paths passed in the file
			nbBytes, errWS := tmpcf.WriteString(path + "\n")
			if errWS != nil {
				fmt.Println(color.RedString("::Are you sure you can handle this much? Without askin your mom first!?"))
			}
			fmt.Println(color.GreenString("[NEW TEMP]::"+fmt.Sprintf("\t%d", nbBytes)+"::>"), path)
		}
		fmt.Println(color.YellowString("::"), "All paths were added to TEMPest !")
	}

	return nil
}

// checkRedondance returns true if a string is contained in both slices
// Eventually it is used to check if the path has already been set in ~/.tempestcf
// It also checks ``sliceArgs`` with ``sliceArgs`` to see if there were multiple same paths
// provided as arguments
func checkRedondance(slice, sliceArgs []string) (doesit bool) {
	for _, pathExisting := range slice {
		// fmt.Println(pathExisting)
		for _, anArg := range sliceArgs {
			if pathExisting == anArg {
				fmt.Println(redB("::"), color.RedString("This path is already taken care of by sweet TEMPest:"), anArg)
				// return errors.New("Won't override the path in ~/.tempestcf")
				doesit = true
			}
		}
	}

	// now check in the slice passed as args with itself
	sliceArgs2 := sliceArgs
	var count int
	for _, pathArg1 := range sliceArgs {
		for _, pathArg2 := range sliceArgs2 {
			if pathArg1 == pathArg2 {
				count++
			}
		}
	}
	if count > len(sliceArgs) {
		fmt.Println(redB("::"), color.RedString("A path has been provided twice to TEMPest"))
		doesit = true
	}
	return
}

// func isValidTarget(slice []string) (dir bool, inexistant error) {
// 	for _, singlePath := range slice {
// 		if dir, errExist := IsDirectory(singlePath); errExist != nil {
// 			return dir, errExist
// 		}
// 	}
// 	return dir, nil
// }
