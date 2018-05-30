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
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vrischmann/envconfig"
)

var this string

// autoAdd indicates if TEMPest should look for all ``temp.est`` dirs and
// add them as targets
var autoAdd bool

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
		// Auto add flag used
		if autoAdd && len(args) == 0 {
			toAdd, errDirs := findDirs("/home/chacanterg/", "temp.est")
			if errDirs != nil {
				log.Fatal(errDirs)
			}
			toAdd, errStrip := stripExistingTargets(toAdd)
			lenAdd := len(toAdd)
			if lenAdd > 0 {
				if errStrip != nil {
					fmt.Println(errStrip.Error())
					return
				}
			} else {
				fmt.Println(cyanB("[INFO]::"), color.HiCyanString("No paths were added"))
			}
			// add these lines
			if errAddLine := addLine(toAdd); errAddLine != nil {
				fmt.Println("::An error occurred while adding path(s):\n", errAddLine)
			}
		} else if !autoAdd {
			// FALLBACK CASE
			if errAddLine := addLine(args); errAddLine != nil {
				fmt.Println("::An error occurred while adding path(s):\n", errAddLine)
			}
		} else {
			cmd.Help()
		}

	},
}

func init() {
	if err := envconfig.Init(&conf); err != nil {
		log.Fatal(err)
	}

	RootCmd.AddCommand(addCmd)

	// addCmd.PersistentFlags().String("foo", "", "A help for foo")
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// addCmd.Flags().StringVarP(&this, "this", "t", "nothing", "Points to current directory")
	addCmd.Flags().BoolVarP(&autoAdd, "auto", "a", false, "Look for all ``temp.est`` directories of the system and add them to the targets list of TEMPest if they are not already there.")
}

// addLine add each string as target into TEMPest (~/.tempestcf)
func addLine(args []string) error {
	// Check if the path already exists first in tempestcf
	// ctnt, errRead := ioutil.ReadFile(conf.Home + "/.tempestcf")
	ctnt, errRead := ioutil.ReadFile(Tempestcf)
	if errRead != nil {
		fmt.Println(redB("::"), color.RedString("Could not read the file muthafuckkah!"))
		return errRead
	}

	// Store it as a slice of strings
	ctntStr := string(ctnt)
	ctntSlice := strings.Split(ctntStr, "\n")

	// Open the file to write to (adding new temp path)
	// tmpcf, err := os.Open(conf.Home + "/.tempestcf")
	tmpcf, err := os.OpenFile(Tempestcf, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(color.RedString("::Sorry! Could not find "), viper.ConfigFileUsed())
		return err
	}
	defer tmpcf.Close()

	// Fetch current directory
	this, errDir := os.Getwd()
	if errDir != nil {
		log.Fatal(errDir)
	}

	// if no args
	if len(args) == 0 {
		// Check if already exists
		sliceThis := make([]string, 1)
		sliceThis[0] = this
		if checkRedondance(ctntSlice, sliceThis) {
			return errors.New(blueB(":: ") + color.HiBlueString("A path provided is already in TEMPest! Later b*tch"))
		}

		// Add this directory to the file
		nbBytes, errWrite := tmpcf.WriteString(this + "\n")
		if errWrite != nil {
			fmt.Println(redB("::"), color.RedString("Could not write to the file. Fail bitch!"))
			return errWrite
		}
		fmt.Println(greenB(":: [SUCCESS] ")+color.GreenString(fmt.Sprintf("NEW TARGET::%d", nbBytes)+"::>\t"), this)
	} else {
		// Treat the last character
		for ind, onePath := range args {
			args[ind] = TreatLastChar(onePath)
		}
		// Check if already exists
		if checkRedondance(ctntSlice, args) {
			return errors.New(blueB(":: ") + color.HiBlueString("A path provided is already in TEMPest! Later b*tch"))
		}

		for _, path := range args {
			if _, errIsDir := IsDirectory(path); errIsDir != nil {
				fmt.Println(redB("::"), color.RedString("Can't add this target, check it does exist >>"), path)
				return errors.New("")
			}
			// In case it is a relative path
			treatRelativePath(&path, this)

			// Add all the paths passed in the file
			nbBytes, errWS := tmpcf.WriteString(path + "\n")
			if errWS != nil {
				fmt.Println(color.RedString(":: Are you sure you can handle this much? Without askin your mom first!?"))
			}
			// fmt.Println(color.GreenString("[NEW TEMP]::"+fmt.Sprintf("%d", nbBytes)+"::>"), path)
			fmt.Println(greenB(":: [SUCCESS] ")+color.GreenString(fmt.Sprintf("NEW TARGET::%d", nbBytes)+"::>\t"), path)
		}
		fmt.Println(color.YellowString("::"), "All paths were added to TEMPest !")
	}

	return nil
}

// treatRelativePath treats a string in case it is a relative path.
// The func does check if it should act and does the action.
// if the first character is not a pathSeparator or the first two chars are ``./``
// It shall be replaced by the full path of the working directory followed by a path separator.
func treatRelativePath(path *string, workingDir string) {
	// setup
	var matchSimpleRel string
	var matchDotRel string
	workingDir += string(os.PathSeparator)

	// regexes
	// for match
	regSimpleRel := regexp.MustCompile(`^([[:alpha:]]|[\d])`)
	regDotRel := regexp.MustCompile(`^(\.(\/|\\))`)

	// Matching
	// Simple Relative path (with the form ``Documents/temp``)
	matchSimpleRel = regSimpleRel.FindString(*path)
	// Dot Relative path (with the form ``./Documents/temp``)
	matchDotRel = regDotRel.FindString(*path)

	switch {
	case matchSimpleRel != "":
		// replace
		*path = regSimpleRel.ReplaceAllString(*path, workingDir+matchSimpleRel)
	case matchDotRel != "":
		// replace
		*path = regDotRel.ReplaceAllString(*path, workingDir)
	}
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

// TreatLastChar takes a string as parameter.
// It analyzes the last character of this string,
// if it is a path separator character, it gets removed.
// Returns the new path, wether there was change or not.
func TreatLastChar(str string) string {
	if str[len(str)-1:] == string(os.PathSeparator) {
		str = str[:len(str)-1]
	}
	return str
}

// findDirs returns all paths of the directories matching the pattern from
// the root path provided
func findDirs(root, pattern string) ([]string, error) {
	dirs := make([]string, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == pattern {
			dirs = append(dirs, path)
		}
		return nil
	})

	return dirs, err
}

// stripExistingTargets takes the slice of targets to add as parameter
// returns a new slice of targets to register but without the existing ones
func stripExistingTargets(wantAdd []string) ([]string, error) {
	strippedList := make([]string, 0)
	existingTgts, errAllTgt := getPaths()

	if errAllTgt != nil {
		if errAllTgt.Error() == "empty" {
			for _, tgt := range wantAdd {
				strippedList = append(strippedList, tgt)
			}
			return strippedList, nil
		}
		fmt.Println(redB("[ERROR]::"), color.HiRedString("Failed to fetch existing targets. Check if the config is right\n\t"), errAllTgt)
		return nil, errAllTgt
	}

	for _, tgt := range wantAdd {
		if !IsStringInSlice(tgt, existingTgts) {
			strippedList = append(strippedList, tgt)
		}
	}

	return strippedList, nil
}
