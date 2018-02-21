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
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge deletes everything within a directory older than the number of days provided.",
	Long: `Purge deletes everything within a directory older than the number of days provided.

The number of days is defined in the conf file default location: ~/.tempest.yaml.

We recommand to use the default system, meaning:
	#1 Add the path to tempest
	#2 Call 'tempest start' in order to start the deletion
		If you mess up on purge, it's done, your life might end !!
		tempest add ...<path>
		tempest start

Test mode, very convenient to see what would get deleted. For example:
	tempest purge -p /tmp -t

For real purge, use the same command but without the '-t' flag:
	tempest purge -p /tmp

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		/* for _, v := range f {
			fmt.Println(v.Name())
		} */
		switch {
		case pPath != "":
			f, e := fetchAll(pPath)
			if e != nil {
				fmt.Println("-----", e)
			}
			//call str
			if errDeletion := deleteAllStr(pPath, f, tTest); errDeletion != nil {
				color.Red("FCK ALL")
				fmt.Println(errDeletion)
			}
		case pInt > -1:
			//call int
			if errInt := deleteAllInt(pInt, tTest); errInt != nil {
				color.Red("FCK ALL")
				fmt.Println(errInt)
			}
		}
		// DEBUG:
		// fmt.Println("purge called")
	},
}

// var days int

// pPath is the path to use in purge
var pPath string

// tTest is defining if it should run test mode or not
var tTest bool

// pInt default is -1, if not, we use the index number for deletion
var pInt int

func init() {
	RootCmd.AddCommand(purgeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// purgeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// purgeCmd.Flags().IntVarP(&days, "days", "d", -1, "age of the files to be deleted in days")
	purgeCmd.Flags().StringVarP(&pPath, "path", "p", "", "path to purge")
	purgeCmd.Flags().BoolVarP(&tTest, "test", "t", false, "test mode -- doesn't actually delete and log to stdout")
	purgeCmd.Flags().IntVarP(&pInt, "index", "i", -1, "use the index number pointed by tempest list for purge.")
	// purgeCmd.Flags().Parse()

}

// fetchAll retrieves all files AND directories from the root path provided
// returns a slice of os.FileInfo and an error
func fetchAll(root string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		color.HiRed("Not a directory, de yoo no da wae !?")
		return nil, err
	}

	return files, err
}

// deleteAllStr is the basic deletion func in here. It deletes everything
// within the directory pointed by the path provided
// Doesn't delete anything if testMode is true.
// It just displays what would be deleted
func deleteAllStr(path string, targets []os.FileInfo, testMode bool) error {
	path += string(os.PathSeparator)
	days := viper.GetInt("duration")
	smthToDel := false
	if days > 1 {
		color.HiMagenta("List of item to be removed:\n\n")
		color.HiMagenta("Size\tUnit\t\t Item")
		if len(targets) == 0 {
			smthToDel = true
		}
		for _, target := range targets {
			// Check if time is right
			timeDiff := time.Now().Sub(target.ModTime()).Hours()
			if timeDiff >= float64(days*24) {
				var size = float64(target.Size())
				size *= 0.001
				size = Round(size, .5, 2)
				// if target.ModTime() <
				switch {
				case testMode:
					// It is in test mode
					fmt.Println(color.HiCyanString(fmt.Sprintf("%v\tKBytes", size)), "\t\t", path+target.Name())
					smthToDel = true
				default:
					// Then it is an actual deletion
					/* color.HiMagenta("List of item removed:\n")
					color.HiMagenta("Size\t\tItem") */
					fmt.Println(color.HiCyanString(fmt.Sprintf("%v\tKBytes", size)), "\t\t", path+target.Name())
					errRemove := os.RemoveAll(path + target.Name())
					if errRemove != nil {
						return errRemove
					}
					smthToDel = true
				}
			}
		}
		// Comment on action
		if !smthToDel {
			fmt.Println("0\tKB\t\t Forever Alone ?")
		}
		if testMode {
			fmt.Println(greenB("::"), color.HiGreenString("Nothing got removed, it is just a recap of what would get deleted <_<"))
			fmt.Println(greenB("::"), color.HiGreenString("CHAMPAGNE !"))
		} else {
			fmt.Println(greenB("::"), color.HiGreenString("All done for"), path)
			fmt.Println(greenB("::"), color.HiGreenString("CHAMPAGNE !"))
		}
	} else {
		fmt.Println(redB("::"), color.HiRedString("Cannot delete files younger than 1 day!"))
		return errors.New("")
	}

	return nil
}

// deleteAllInt deletes everything within the directory pointed by the path
// provided, using the index from ``.tempestcf``
// Doesn't delete anything if testMode is true.
// It just displays what would get deleted
func deleteAllInt(index int, testMode bool) error {
	allPaths, errPaths := getPaths()
	if errPaths != nil {
		color.Red("::Error while reading .tempestcf")
		return errPaths
	}
	if index >= 0 && index < len(allPaths) {
		for indx, indPath := range allPaths {
			if index == indx {
				// color.Cyan(indPath)
				fInt, eInt := fetchAll(indPath)
				if eInt != nil {
					fmt.Println("-----", eInt)
				}
				return deleteAllStr(indPath, fInt, testMode)
			}
		}
	}
	return errors.New("Nothing to purge")
}
