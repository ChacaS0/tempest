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
	"github.com/spf13/cobra/doc"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all the paths submitted to TEMPest.",
	Long: `All the paths set for TEMPest.

They follow this pattern:
<IndexPath> <Path>
For example:
0 /tmp
1 /temp
2 /tempo
3 /tempor

The IndexPath can then be used to select the path 
`,
	Run: func(cmd *cobra.Command, args []string) {
		if errLi := printList(); errLi != nil {
			fmt.Println(color.HiRedString("Could not list paths, sorry bra!", errLi))
		}
	},
}

var listHeader = &doc.GenManHeader{
	Title:   "Start",
	Section: "3",
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
			fmt.Println(color.HiMagentaString(":: No path set yet\n:: Suggestion - Run: \n\ttempest help add\nFor more information about adding paths!"))
			return nil
		}
		return errSliced
	}

	// color.Red(fmt.Sprintf("%d", len(ctntSlice)))
	if len(ctntSlice) >= 1 {
		fmt.Println(color.HiYellowString("Current paths currently having \"fun\" with TEMPest:\n"))
		fmt.Println(color.HiYellowString("Index\t| Path"))
		// fmt.Println(color.HiYellowString("--------------------------------------------------"))
	}

	for i, aPath := range ctntSlice {
		// fmt.Println(color.HiYellowString(fmt.Sprintf("%d\t|", i)), aPath)
		switch {
		case i < len(aPath):
			fmt.Println(color.HiYellowString(fmt.Sprintf("%d\t|", i)), aPath)
		case i == 0:
			fmt.Println(color.HiMagentaString("No path set yet\nSuggestion - Run: \n\ttempest help add\nFor more information about adding paths!"))
		}
	}
	return nil
}

// getPaths returns a slice of string being all the paths for TEMPest to purge
// and an error if something goes wrong
func getPaths() (returnSlice []string, pathsError error) {
	ctnt, pathsError := ioutil.ReadFile(conf.Home + "/.tempestcf")
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
