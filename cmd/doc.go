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
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var isMan bool

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Display the content of README.md",
	Long: `Display the content of README.md.

	We highly recommand to use Showdown with template doc.
		It can be found at: https://github.com/craigbarnes/showdown
		Or in the aur/showdown-git
`,
	Run: func(cmd *cobra.Command, args []string) {
		launchDoc()
		// fmt.Println("doc called")
	},
}

func init() {
	RootCmd.AddCommand(docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	docCmd.Flags().BoolVarP(&isMan, "man", "m", false, "Add this falg to have a \"man like\" view.")
}

func launchDoc() {
	pathReadme := pathTempest + "README.md"
	// if no flag "man" added, open with Showdown
	if !isMan {
		commShowdown := exec.Command("showdown", pathReadme)
		errorShowdown := commShowdown.Run()
		if errorShowdown != nil {
			// then showdown is not installed
			color.HiCyan("::We recommand you to install Showdown to view this")
			color.HiCyan("::For more information, please visit: https://github.com/craigbarnes/showdown")
			// open in browser
			openBrowser("file://" + pathReadme)
		}
	} else {
		/* // Else open in man mode
		comm := "cat " + pathReadme + " | less"
		// print result of command
		commMan := exec.Command("bash", "-c", comm)
		_ = commMan.Run()
		outputMan, errMan := commMan.Output()
		if errMan != nil {
			color.Red("Error while trying to get into the pussy")
			return
		} */
		/* cat := exec.Command("cat " + pathReadme)
		less := exec.Command("less")

		catOut, _ := cat.StdoutPipe()
		cat.Start()
		less.Stdin = catOut

		outputMan, errMan := less.Output()
		if errMan != nil {
			color.Red("Cant get into the pussy")
		} */

		/* cat := exec.Command("cat " + pathReadme + " | less")
		pipe, errPipe := cat.StderrPipe()
		if errPipe != nil {
			color.Red("Error piping")
		}
		errWait := cat.Wait() */

		// fmt.Println(string(outputMan))
		color.HiMagenta("\nSuck it up, it's not there yet lel\n")

		// // https://gist.github.com/kylelemons/1525278
		// Collecting the ouput of the commands
		/* var output bytes.Buffer
		var stderr bytes.Buffer

		cat := exec.Command("cat " + pathReadme + " | less") */
	}
	// return nil
}
