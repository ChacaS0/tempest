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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Call this if it is the first time you use TEMPest or to reset the temps",
	Long: `Call this if it is the first time you use TEMPest or to reset the temps:

It will create a file named .tempestcf in your /home/$USER.
This file will contain the list of all the directories you wish tempest to handle as temporay directory.

`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initializeTP(); err != nil {
			log.Fatalln(err)
		}
		fmt.Println(color.HiGreenString(`You are now ready to use TEMPest.

Suggestions:
	Start using TEMPest right away by adding a temporay file :
		tempest add <FILE_PATH>
`))
	},
}

func init() {

	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initializeTP creates the ~/.tempestcf file or
// empty it if already exists
func initializeTP() error {
	f, err := os.OpenFile(conf.Home+"/.tempestcf", os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errDel := os.Remove(conf.Home + "/.tempestcf")
		if errDel != nil {
			fmt.Println(color.RedString("::Error while replacing existing file"))
			return errDel
		}
		f2, err2 := os.OpenFile(conf.Home+"/.tempestcf", os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			fmt.Println(color.HiRedString("::Huge error! All data lost, could not recreate file! Your life officially sucks!"))
			return err2
		}
		defer f2.Close()
	}
	if err := f.Close(); err != nil {
		// fmt.Println("here dude")
		log.Fatal(err)
	}

	return nil
}
