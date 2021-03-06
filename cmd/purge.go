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

// var days int

// pPath is the path to use in purge
var pPath string

// tTest is defining if it should run test mode or not
var tTest bool

// pInt default is -1, if not, we use the index number for deletion
var pInt int

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
	tempest purge -p /tmp/temp.est -t

For real purge, use the same command but without the '-t' flag:
	tempest purge -p /tmp/temp.est

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		/* for _, v := range f {
			fmt.Println(v.Name())
		} */
		switch {
		case pPath != "":
			f, fInfo, e := fetchAll(pPath)
			if e != nil {
				fmt.Println("-----", e)
			}
			if fInfo == nil {
				//call str
				if errDeletion := deleteAllStr(pPath, f, tTest); errDeletion != nil {
					color.Red("FCK ALL")
					fmt.Println(errDeletion)
				}
			} else {
				if errEmptyF := emptyFile(pPath, fInfo, tTest); errEmptyF != nil {
					color.Red("FCK ME")
					fmt.Println(errEmptyF)
				}
			}
		case pInt > -1:
			//call int
			if errInt := deleteAllInt(pInt, tTest); errInt != nil {
				color.Red("FCK ALL")
				fmt.Println(errInt)
			}
		default:
			cmd.Help()
		}
	},
}

func init() {
	RootCmd.AddCommand(purgeCmd)

	// purgeCmd.PersistentFlags().String("foo", "", "A help for foo")
	// purgeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// purgeCmd.Flags().IntVarP(&days, "days", "d", -1, "age of the files to be deleted in days")

	purgeCmd.Flags().StringVarP(&pPath, "path", "p", "", "path to purge")
	purgeCmd.Flags().BoolVarP(&tTest, "test", "t", false, "test mode -- doesn't actually delete and log to stdout")
	purgeCmd.Flags().IntVarP(&pInt, "index", "i", -1, "use the index number pointed by tempest list for purge.")
	// purgeCmd.Flags().Parse()

}

// fetchAll retrieves all files AND directories from the root path provided if the path points to a directory.
// If it points to a file, it returns the os.FileInfo of that file.
// If there is an error, it will be returned
func fetchAll(root string) (sliceDir []os.FileInfo, targetInfo os.FileInfo, errFunc error) {
	sliceDir, errFunc = ioutil.ReadDir(root)
	if errFunc != nil {
		// Not a directory so probably a file ? Check for it
		targetInfo, errFunc = os.Stat(root)
		if errFunc != nil {
			fmt.Println(redB(":: [ERROR]"), color.HiRedString("Could not read the file(s), sad story!"))
			targetInfo = nil
		}
		sliceDir = nil
	}
	return
}

// emptyFile is one of the main func in here. It deletes everything within
// the file pointed by the path provided.
// Doesn't delete anything if testMode is true.
// It just displays
func emptyFile(path string, target os.FileInfo, testMode bool) error {

	// header message
	fmt.Println(magB("\n:: File:"), path)
	fmt.Println(magB("Targeted   Size\tUnit\t Item"))

	//? Maybe add a new age speciffic to files?
	// get the age in config
	days := viper.GetInt("duration")

	// check the age
	timeDiff := time.Now().Sub(target.ModTime()).Hours()

	// infos
	size, unit := FormatSize(float64(target.Size()))

	var msg string

	if timeDiff >= float64(days*24) {

		if testMode {
			// Test mode, don't delete
			// TODO: Improve
			msg = fmt.Sprintln(redB("YES\t   "), color.HiCyanString(fmt.Sprintf("%v\t%s", size, unit)), "\t", path+target.Name())
		} else {
			// Actual deletion
			if err := os.Truncate(path, 0); err != nil {
				return err
			}
			msg = fmt.Sprintln(redB("DONE\t   "), color.HiCyanString(fmt.Sprintf("%v\t%s", size, unit)), "\t", path+target.Name())
		}
	} else {
		msg = fmt.Sprintln(greenB("NOPE\t   "), color.HiCyanString(fmt.Sprintf("%v\t%s", size, unit)), "\t", path+target.Name())
	}

	fmt.Println(msg)

	return nil
}

// deleteAllStr is one of the main func in here. It deletes everything
// within the directory pointed by the path provided.
// Doesn't delete anything if testMode is true.
// It just displays what would be deleted
func deleteAllStr(path string, targets []os.FileInfo, testMode bool) error {
	path += string(os.PathSeparator)
	days := viper.GetInt("duration")
	smthToDel := false
	if days > 1 {
		// color.HiMagenta("List of item to be removed:\n\n")
		fmt.Println(magB("\n:: List of items to be removed in:"), path)
		// color.HiMagenta("Size\tUnit\t\t Item")
		fmt.Println(magB("Size\tUnit\t\t Item"))
		if len(targets) == 0 {
			smthToDel = true
		}
		for _, target := range targets {
			// Check if time is right
			timeDiff := time.Now().Sub(target.ModTime()).Hours()
			if timeDiff >= float64(days*24) {
				size, unit := FormatSize(float64(target.Size()))

				// if target.ModTime() <
				switch {
				case testMode:
					// It is in test mode
					fmt.Println(color.HiCyanString(fmt.Sprintf("%v\t%s", size, unit)), "\t\t", path+target.Name())
					smthToDel = true
				default:
					// Then it is an actual deletion
					/* color.HiMagenta("List of item removed:\n")
					color.HiMagenta("Size\t\tItem") */
					fmt.Println(color.HiCyanString(fmt.Sprintf("%v\t%s", size, unit)), "\t\t", path+target.Name())
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
			fmt.Println("0\tKB\t\t Forever Alone ? Nothing to remove here !")
		}
		if testMode {
			fmt.Println(greenB("::"), color.HiGreenString("Nothing got removed, it is just a recap of what would get deleted <_<"))
			fmt.Println(greenB("::"), color.HiGreenString("CHAMPAGNE !"))
		} else {
			fmt.Println(greenB("\n::"), color.HiGreenString("All done for"), path)
			fmt.Println(greenB("::"), color.HiGreenString("CHAMPAGNE !\n"))
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
		color.Red(":: Error while reading .tempestcf")
		return errPaths
	}
	if index >= 0 && index < len(allPaths) {
		for indx, indPath := range allPaths {
			if index == indx {
				// color.Cyan(indPath)
				fInt, fInfo, eInt := fetchAll(indPath)
				if eInt != nil {
					fmt.Println("-----", eInt)
				}
				if fInfo != nil {
					return emptyFile(indPath, fInfo, testMode)
				}
				return deleteAllStr(indPath, fInt, testMode)
			}
		}
	}
	return errors.New("Nothing to purge")
}
