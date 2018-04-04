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
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gAge is true if the user wants to know the age set in config
var gAge bool

// allLogs tells if the user wants to see logs
var allLogs bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrives information, mainly from config",
	Long: `Retrives information, mainly from config

Get is pretty simple to use. For example:
	tempest get --age

`,
	Run: func(cmd *cobra.Command, args []string) {
		printAnyIfSet(args)
		// color.HiCyan("\nNot implemented yet!\n")
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	getCmd.Flags().BoolVarP(&gAge, "age", "a", false, "Use this if you want to know the age set in .tempest.yaml")
	getCmd.Flags().BoolVarP(&allLogs, "logs", "l", false, "Use this flag to view all logs")
}

// getAllLogs handle the use of ``logs`` flag.
// We assume that ``allLogs`` is sets to true when calling this func.
func getAllLogs(args []string) error {
	// var declare
	var showShutup bool

	//* Setup rights for shutup
	if len(args) == 0 {
		showShutup = true
	}

	//* Display Log Shutup
	if is, err := IsDirectory(LogShutup); !is && err == nil && showShutup {
		// Fetch the content of LogShutup
		fileCtnt, errRF := ioutil.ReadFile(LogShutup)
		if errRF != nil {
			fmt.Println(redB(":: [ERROR]"), color.HiRedString("Could not read the file -_-\n\t->", LogShutup, "\n\t->", errRF))
		}
		fmt.Println(magB("===========================================  - [ShutupLogs] -  ==================================================="))
		fmt.Println(string(fileCtnt))
		fmt.Println(magB("========================================  - [EOF - ShutupLogs] -  ================================================"))
	}
	return nil
}

// printAnyIfSet displays the config set for the ones asked.
// If none is asked, it shows everything
func printAnyIfSet(args []string) {
	switch {
	case gAge:
		// Age A.K.A. the duration
		fmt.Println(blueB("::"), whiteB("Age:"), getAge())
	case allLogs:
		getAllLogs(args)
	default:
		// help ?
		getHelp()
	}
}

// getAge returns the age as a string
func getAge() string {
	return fmt.Sprintf("%d", viper.GetInt("duration"))
}

// getHelp() calls the regualr helpCommand
func getHelp() {
	fmt.Println(RootCmd.UsageString())
}
