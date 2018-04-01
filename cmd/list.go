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
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// doFix is the var receiving the value of the --fix flag
// true if needs to fix the targets
var doFix bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the targets submitted to TEMPest.",
	Long: `All the targets set for TEMPest.

They follow this pattern:
<IndexPath> <Path>
For example:
0 /tmp
1 /temp
2 /tempo
3 /tempor

The IndexPath can then be used to select the path 

To fix broken targets (for example targets that points to non-existing paths):
	tempest list --fix
`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case doFix && len(args) == 0:
			if errFix := fixTargets(); errFix != nil {
				fmt.Println(redB(":: [ERROR]"), color.HiRedString("Could not fix broken paths, feels bra!"), errFix)
			}
		default:
			if errLi := printList(); errLi != nil {
				fmt.Println(redB(":: [ERROR]"), color.HiRedString("Could not list targets, sorry bra!", errLi))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().BoolVarP(&doFix, "fix", "f", false, "Check if TEMPest is watching some broken paths and act on them in this case.")

	/* //* TEST
	out := new(bytes.Buffer)
	doc.GenMan(listCmd, listHeader, out)
	fmt.Print(out.String()) */
}

func printList() error {
	// Check if the path already exists first
	ctntSlice, errSliced := getPaths()
	if errSliced != nil {
		// Just a small enhancement of the "no paths set yet" display
		if errSliced.Error() == "empty" {
			fmt.Println(color.HiMagentaString(":: No target set yet\n:: Suggestion - Run: \n\ttempest help add\nFor more information about adding targets!"))
			return nil
		}
		return errSliced
	}

	// color.Red(fmt.Sprintf("%d", len(ctntSlice)))
	if len(ctntSlice) >= 1 {
		fmt.Println(color.HiYellowString("Current targets currently having \"fun\" with TEMPest:\n"))
		fmt.Println(color.HiYellowString("Index\t| Target"))
		// fmt.Println(color.HiYellowString("--------------------------------------------------"))
	}

	for i, aPath := range ctntSlice {
		// fmt.Println(color.HiYellowString(fmt.Sprintf("%d\t|", i)), aPath)
		switch {
		case i < len(aPath):
			fmt.Println(color.HiYellowString(fmt.Sprintf("%d\t|", i)), aPath)
		case i == 0:
			fmt.Println(color.HiMagentaString("No target set yet\nSuggestion - Run: \n\ttempest help add\nFor more information about adding targets!"))
		}
	}
	return nil
}

// getPaths returns a slice of string being all the paths for TEMPest to purge
// and an error if something goes wrong
func getPaths() (returnSlice []string, pathsError error) {
	ctnt, pathsError := ioutil.ReadFile(Tempestcf)
	if pathsError != nil {
		fmt.Println(color.RedString("::Could not read the file muthafuckkah!"))
		return nil, pathsError
	}

	if len(ctnt) == 0 {
		return returnSlice, errors.New("empty")
	}

	// Store it as a slice of strings
	ctntStr := string(ctnt)
	returnSlice = strings.Split(ctntStr, "\n")
	return returnSlice[:len(returnSlice)-1], nil
}

// fixTargets is a func that handle the fix paths process
func fixTargets() error {

	// first we fetch all current paths
	allPaths, errGP := getPaths()
	if errGP != nil {
		return errGP
	}

	// Convert those as targets
	allTargets := PathsToTargets(allPaths)

	// Get the state of each target
	states := getState(allTargets)

	// slRmInt is the slice that will hold all the targets to remove (subject to changes to use Targets instead)
	slRmInt := make([]int, 0)
	for tgt, state := range states {
		if !state {
			slRmInt = append(slRmInt, tgt.Index)
		}
	}

	// Then just process those trashes
	// TODO - Later maybe we could ask the user what to do ?
	fmt.Println(yellowB("::"), yellowB("Status of targets:"))
	fmt.Println(yellowB("\nStatus\t| Index\t | Target"))
	for tgt, ste := range states {
		if !ste {
			fmt.Print(color.HiRedString("BROKEN"), yellowB("\t|"))
		} else {
			fmt.Print(color.HiGreenString("GOOD"), yellowB("\t|"))
		}
		fmt.Println(" ", tgt.Index, "\t", yellowB("|"), tgt.Path)
	}
	if len(slRmInt) > 0 {
		fmt.Println("")
		slicePaths := rmInSlice(slRmInt, []string{}, allPaths)
		if err := writeTempestcf(slicePaths); err != nil {
			fmt.Println(redB(":: [ERROR]"), color.HiRedString("Did not succeed to write the new targets: \n\t", err.Error()))
		}
		// INFO
		fmt.Println(color.CyanString("::"), color.HiCyanString("Broken targets deleted with success ! ;)"))
	} else {
		// INFO
		fmt.Println("\n", color.CyanString("::"), color.HiCyanString("No broken targets ! Much WOW ! !"))
	}

	return nil
}

// getState returns a map[Target]bool named ``states``
// states relates the state of the target's path.
// {
// 	True  : Still fine, path is right,
// 	False : Nope, broken path or doesn't exists
// }
func getState(targets []Target) map[Target]bool {

	states := make(map[Target]bool, 0)

	for _, tgt := range targets {
		if _, err := IsDirectory(tgt.Path); err == nil {
			states[tgt] = true
		} else {
			states[tgt] = false
		}
	}

	return states
}
