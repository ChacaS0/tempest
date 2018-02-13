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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gAge is true if the user wants to know the age set in config
var gAge bool

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
}

// printAnyIfSet displays the config set for the ones asked.
// If none is asked, it shows everything
func printAnyIfSet(args []string) {
	// Age A.K.A. the duration
	if gAge {
		fmt.Println(blueB("::")+whiteB("Age:"), "\t", getAge())
	}
	// help ?
	if len(args) == 1 && args[0] == "help" {
		getHelp()
	}
}

// getAge returns the age as a string
func getAge() string {
	return fmt.Sprintf("%d", viper.GetInt("duration"))
}

func getHelp() {
	fmt.Println(RootCmd.UsageString())
}
