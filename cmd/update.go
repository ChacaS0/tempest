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
	"log"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// var repo string = ""

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the TEMPest tool",
	Long: `Simply update the TEMPest tool. Exmaple:

	tempest update

	It checks first if there are updates available. If there are, start updating.
	
	/!\ Careful, sometimes it takes some time to check/update depending on your internet speed.
`,
	Run: func(cmd *cobra.Command, args []string) {
		switch isUpToDate, errUpt := checkUpdate(); {
		case !isUpToDate && errUpt == nil:
			errPull := gitPull()
			if errPull != nil {
				color.Red("Error while pulling the new content!")
				fmt.Println(errPull)
			}
		case isUpToDate && errUpt == nil:
			color.HiBlue("You are already up to date! Someone call TV, we have Captain Obvious on duty here!")
		default:
			color.Red("Error while checking for updates !")
			fmt.Println(errUpt)
		}

		// color.HiYellow("Does nothing for now, sad life for sad people hunh ?")
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// checkUpdate returns true if the current code is different from the one in
// the bitbucket repository
func checkUpdate() (bool, error) {
	// Fetch data about the remote TEMPest
	fmt.Println(yellowB("::"), color.HiYellowString("Fetching data about the incoming TEMPest . . ."))
	commFetch := exec.Command("git", "fetch")
	commFetch.Dir = pathTempest
	fetchResult, errFetch := commFetch.Output()
	if errFetch != nil {
		commFetch = exec.Command("git", "fetch", "origin", "master")
		color.Red("Error while fetching stretching legs and branches!!")
		return false, errFetch
	}
	fmt.Println(string(fetchResult))

	// Check if up to date
	fmt.Println(yellowB("::"), color.HiYellowString("Checking if the TEMPest is coming this way . . ."))
	log.Println(":::::")
	commCheck := exec.Command("git", "log", "HEAD..origin/master", "--oneline")
	commCheck.Dir = pathTempest
	resultComm, errCheck := commCheck.Output()
	if errCheck != nil {
		color.Red("Error while checking for updates, duuuuuuude!")
		return false, errCheck
	}
	// color.HiGreen(fmt.Sprintf("%d", len(resultComm)))
	if len(resultComm) > 0 {
		// color.Cyan("false")
		return false, nil
	}
	return true, nil
}

// gitPull does a simple git pull on origin/master
func gitPull() error {
	// set the working dir as the one containing the code
	fmt.Println(yellowB("::"), color.HiYellowString("The TEMPest is here ! Hid & run for your lives ! ! !"))
	cmdx := exec.Command(pathTempest + "make.sh")
	output, errRuning := cmdx.Output()
	fmt.Println(string(output))
	if errRuning == nil {
		fmt.Println(yellowB("::"), color.HiYellowString("TEMPest now up to date ! :p"))
	}
	return errRuning
}
