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

// testAll tells if ``start`` should be run in 'test mode'
var testAll bool

// needShutup tells if start should be run in 'shutup mode'
var needShutup bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "The purge begins ! Run !!",
	Long: `For now, 'tempest purge' only purges all the paths added to it.
It doesn't handle the purge of one speciffic path. 
If that's what you want to do, you should try to use 'tempest help purge' for more information.

Example of start:
tempest start
	--> /!\ This actually deletes stuff !!

Test mode:
tempest start -t
	--> Test mode doesn't delete anything, but gives crucial informations about what would be deleted
`,
	Run: func(cmd *cobra.Command, args []string) {
		allPaths, err := getPaths()
		// allPaths = allPaths[:len(allPaths)-1]
		if err != nil {
			fmt.Println(color.RedString("Could not find da wae !\n"), err)
		}
		// run in shutup mode
		if needShutup {
			handleShutupMode(allPaths)
			return
		}
		// call purge
		if errPurge := callPurge(allPaths); errPurge != nil {
			fmt.Println(color.RedString("Sorry something went terribly wrong... Feels brah!\n"), errPurge)
		}
		// fmt.Println("start called")
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().BoolVarP(&testAll, "test", "t", false, "test mode -- doesn't actually delete and log to stdout")
	startCmd.Flags().BoolVarP(&needShutup, "shutup", "s", false, "in 'shutup mode', all messages are redirected to logs ("+LogShutup+")") //TODO: Replace LOGPATH by the right var
}

// callPurge call the external program named "purge" on each path provided
// it deletes everything inside a folder but not directories
// returns an error if something goes wrong, otherwise returns <nil>
func callPurge(targets []string) error {
	// nbDays := viper.GetInt("duration")
	// color.HiRed(fmt.Sprintf("%d", nbDays))

	// actual purge
	for _, path := range targets {
		/* f, e := fetchAll(pPath)
		if e != nil {
			return e
		} */
		/* if errDeletion := deleteAll(path, f); errDeletion != nil {
			color.Red("FCK ALL")
			// fmt.Println(errDeletion)
			return errDeletion
		} */
		/* pPath = path
		purgeCmd.Execute() */
		fInfo, errFInfo := fetchAll(path)
		if errFInfo != nil {
			color.HiRed("Error fetch info for this temp directory")
			return errFInfo
		}

		if errPurgeAll := deleteAllStr(path, fInfo, testAll); errPurgeAll != nil {
			color.HiRed("Error while purging all the file at once!")
			return errPurgeAll
		}

	}

	return nil
}

// handleShutupMode is the handler func for the 'shutup mode'.
// We assume that ``needShutup`` is already true.
// TODO - migrate to Target instead of string
func handleShutupMode(targets []string) {
	testAll = true // DEBUG MODE
	callBck := func() {
		if err := callPurge(targets); err != nil {
			fmt.Println(err)
		}
	}
	stdString := captureStdout(callBck)
	// write stdString to the log file
	// TODO: replace that by the right path var
	pathLog := conf.Home + string(os.PathSeparator) + ".tempest" + string(os.PathSeparator) + ".log" + string(os.PathSeparator) + "shutup.log"
	WriteLog(pathLog, string(stdString))
}
