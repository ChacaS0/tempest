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
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Call this if it is the first time you use TEMPest or to reset the temps",
	Long: `Call this if it is the first time you use TEMPest or to reset the temps:

It will create a file named .tempestcf in your /home/$USER.
This file will contain the list of all the directories you wish tempest to handle as temporary directory.

`,
	Run: func(cmd *cobra.Command, args []string) {
		// ~/.tempest DIR
		if err := initTempestDir(); err != nil {
			fmt.Println(redB("::"), color.HiRedString("Sorry could not create the .tempest dir in "), conf.Home, "\n\t[ERROR]::", err)
			return
		}

		// .tempest.yaml
		if err := initializeCfFile(); err != nil {
			fmt.Println(redB("::"), color.HiRedString("Could not initialize .tempest.yaml"))
			fmt.Println(redB("::"), color.HiRedString("If the error persists, try to create the file manually : touch $HOME/.tempest/.tempest.yaml"))
			// fmt.Println(err) //DEBUG
			return
		}

		// .tempestcf
		if err := initializeTP(); err != nil {
			// log.Fatalln(err)
			fmt.Println(redB("::"), color.HiRedString("Could not initialize .tempestcf"))
			fmt.Println(redB(err.Error()))
			return
		}

		// SUCCESS:
		fmt.Println(greenB("::"), color.HiGreenString("You are now ready to use TEMPest."))
		fmt.Println(greenB("::"), color.HiGreenString("Suggestions:"))
		fmt.Println(color.HiGreenString(`	Start using TEMPest right away by adding a temporary file :
		tempest add <DIRECTORY_PATH>
	Or get help to add new paths:
		tempest help add
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

// initTempestDir init func for .tempest directory
func initTempestDir() error {
	// remove if already exists
	if is, err := IsDirectory(TempestConfigDir); is && err == nil {
		if errRm := os.RemoveAll(TempestConfigDir); errRm != nil {
			return errRm
		}
	}
	// fmt.Println("conf dir deleted")
	// Create the conf dir
	if err := os.Mkdir(TempestConfigDir, 0755); err != nil {
		return err
	}
	// fmt.Println("conf dir created")
	// create the log dir
	err := os.Mkdir(TempestConfigDir+string(os.PathSeparator)+".log", 0755)
	// fmt.Println("log dir created if nil =>", err)
	return err
}

// initializeTP creates the ~/.tempestcf file or
// empty it if already exists
func initializeTP() error {
	f, err := os.OpenFile(Tempestcf, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errDel := os.Remove(Tempestcf)
		if errDel != nil {
			fmt.Println(redB("::"), color.RedString("Error while replacing existing file (.tempestcf)"))
			return errDel
		}
		f2, err2 := os.OpenFile(Tempestcf, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			fmt.Println(redB("::"), color.HiRedString("Huge error! All data lost, could not recreate file! Your life officially sucks!"))
			return err2
		}
		defer f2.Close()

		return nil
	}

	if err := f.Close(); err != nil {
		fmt.Println(redB("::"), color.HiRedString("Weird.. could not close .tempestcf"))
	}

	return nil
}

// initializeCfFile creates the file ``$HOME/.tempest.yaml``(Tempestyml)
// if it doesn't already exist with ``duration: 5``
func initializeCfFile() error {
	defConf := `duration: 5
auto-mode: false
`
	_, errDir := IsDirectory(Tempestyml)
	if errDir == nil {
		// if already exists, we delete

		errDel := os.RemoveAll(Tempestyml)
		if errDel != nil {
			fmt.Println(redB("::"), color.RedString("Error while replacing existing file ("+Tempestyml+")"))
			return errDel
		}
	}
	// Doesn't exist so create it!
	testYaml := conf.Gopath + string(os.PathSeparator) + ".tempest.yaml"
	if Tempestyml != testYaml {
		Tempestyml = TempestymlDef
	}
	f, err := os.OpenFile(Tempestyml, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(redB("::"), color.HiRedString("Huge error! Could not recreate file! God lost faith in you!"), Tempestyml)
		return err
	}
	defer f.Close()

	_, errWrite := f.WriteString(defConf)
	if errWrite != nil {
		fmt.Println(redB("::"), color.HiRedString("Could not write the default config to"), Tempestyml)
		fmt.Println(redB("::"), color.HiRedString(`If the problem persists, try add this to it:
	duration: 5
	auto-mode: false
`))
		return errWrite
	}
	// viper.WriteConfigAs(viper.ConfigFileUsed())
	// cfgFile = conf.Home + "/.tempest.yaml"
	viper.SetConfigFile(Tempestyml)
	// Tempestyml = TempestymlDef
	// initConfig()

	return nil
}
