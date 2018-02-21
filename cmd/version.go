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
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the tool",
	Long: `Show the current version of the tool. For example:

	tempest version
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("version called")
		if version, errV := getVersion(); errV == nil {
			fmt.Println(color.HiYellowString(version))
		} else {
			color.Red("error")
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// getVersion returns a string which is the current version
// Fetches the version from git command
// if something goes wrong, an error is returned and the version is an empty string
func getVersion() (string, error) {
	// cmd to use git to get the version
	commVersion := exec.Command("git", "describe")
	commVersion.Dir = pathTempest
	tempOutput, errVersion := commVersion.Output()
	if errVersion != nil {
		fmt.Println(redB("::"), color.HiRedString("Could not read the current version"))
		fmt.Println(redB("::"), color.HiRedString("Still no version out yet?"))

		return "", errVersion
	}
	version := string(tempOutput)
	return version, errVersion
}
