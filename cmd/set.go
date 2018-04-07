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
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags vars
var age int

// autoStart defines if TEMPest should start automatically or not.
//    -> Default is nope.
var autoStart string

// LinuxAutoStartP is the path to $HOME/.config/autostart/tempest.desktop
// which is initiaalized in init()
var LinuxAutoStartP string

// WindowsAutoStartSL is the path to "%AppData%\Microsoft\Windows\Start Menu\Programs\Startup\tempest".
// which is a Symlink to the scripts/startup.bat file, initialized in init()
var WindowsAutoStartSL string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set is used to set variables, content, config on TEMPest",
	Long: `set is used to set variables, content, config on TEMPest. For example:

	tempest set --age 4
	# This sets the maximum "age" of the temp directories content

More features might come next. Stay tuned!
`,
	Run: func(cmd *cobra.Command, args []string) {
		if age != -1 {
			if errAge := setAge(); errAge != nil {
				fmt.Println(errAge)
			}
			// color.HiCyan("Sorry! No working at the moment! -\\('o')/-")
		}
		if autoStart != "" {
			if err := setAutoStart(); err != nil {
				fmt.Println(redB("::"), color.HiRedString(err.Error()))
			}
		}
		if age == -1 && autoStart == "" {
			if err := cmd.Help(); err != nil {
				color.HiRed("You must provide some flags or whatever! Just do something!")
			}
		}
		// fmt.Println("set called")
		/* errWC := viper.MergeInConfig()
		if errWC != nil {
			log.Println(errWC)
		} */
	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	// $HOME/.config/autostart/tempest.desktop
	LinuxAutoStartP = conf.Home + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "autostart" + string(os.PathSeparator) + "tempest.desktop"

	WindowsAutoStartSL = `%AppData%\Microsoft\Windows\"Start Menu"\Programs\Startup\TEMPest-startup`

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().IntVarP(&age, "age", "a", -1, "Set the maximum age (in days) of the temp directories' content")
	setCmd.Flags().StringVarP(&autoStart, "auto-start", "s", "", "Set to 'on' to activate or 'off' to deactivate. (not activated by default)")
}

// setAge sets the "duration" value in the config file (~/.tempest/.tempest.yaml)
// - ``duration`` being the maximum age of the temp's content
func setAge() error {
	if age > 0 {
		// ageStr := fmt.Sprintf("%d", age)
		viper.Set("duration", age)
		viper.Set("auto-mode", viper.GetBool("auto-mode"))

		//* Save config
		errS := viper.WriteConfigAs(viper.ConfigFileUsed())
		if errS != nil {
			color.Red(errS.Error())
		}

		fmt.Println(greenB("::"), color.HiGreenString("New age set at:"), viper.GetInt("duration"))
		return nil
	}
	color.Red("::Age must be greater than 0, genius . . .")
	return errors.New("Error while setting the age")
}

// setAutoStart sets the "auto-mode" value in the config file (~/.tempest/.tempest.yaml).
// - ``auto-mode`` being called sometimes in the code "auto-start"
func setAutoStart() error {
	if autoStart == "" {
		return errors.New(fmt.Sprint("Empty autoStart", autoStart, ";"))
	}

	// set up
	var boolAutoMode bool
	switch autoStart {
	case "on":
		boolAutoMode = true
	case "off":
		boolAutoMode = false
	default:
		return errors.New("Wrong parameter, only 'on' or 'off' are accepted for this flag (quote excluded)")
	}

	// Define config
	viper.Set("duration", viper.GetInt("duration"))
	viper.Set("auto-mode", boolAutoMode)

	//* actual shit
	switch runtime.GOOS {
	case "linux":
		if err := autoStartLinux(boolAutoMode); err != nil {
			return err
		}
	case "windows":
		fmt.Println(cyanB(":: [NOTE]"), color.HiCyanString("This functionality hasn't been fully tested for Windows."))
	// case "darwin":
	// 	err = exec.Command("open", url).Start()
	default:
		return errors.New(fmt.Sprint(redB(":: [ERROR]"), color.HiRedString("Unsupported platform")))
	}

	// Save config
	if err := viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		return err
	}

	fmt.Println(greenB("::"), color.HiGreenString("Auto start is now"), greenB(autoStart))

	return nil
}

// autoStartWindows set the auto start for TEMPest (for Windows - only tested on W10).
// - ``{ should = true }``: activate
// - ``{ should = false }``: deactivate
//
// This func creates a link of the bash file ``$HOME/.config/autostart/tempest.desktop``.
// If the file does not exist, it gets created when autoStart is set to "on".
func autoStartWindows(should bool) error {
	// var ctntAutoF string

	// if exists,and should not exist, we delete
	if is, err := IsDirectory(WindowsAutoStartSL); (!is || err == nil) && !should {
		if errRm := os.Remove(WindowsAutoStartSL); errRm != nil {
			return errRm
		}
	} else if should {
		// We create a symlink to the startup folder
		os.Symlink(conf.Gopath+string(os.PathSeparator)+"scripts"+string(os.PathSeparator)+"startup.bat", WindowsAutoStartSL)
	}

	return nil
}

// autoStartLinux set the auto start for TEMPest.
// - ``{ should = true }``: activate
// - ``{ should = false }``: deactivate
//
// This func writes in the file ``$HOME/.config/autostart/tempest.desktop``.
// If the file does not exist, it gets created when autoStart is set to "on".
func autoStartLinux(should bool) error {
	var ctntAutoF string

	// Remove if already existed so we create a fresh one
	_ = os.Remove(LinuxAutoStartP)
	autoF, errCreate := os.OpenFile(LinuxAutoStartP, os.O_RDWR|os.O_CREATE, 0644)
	if errCreate != nil {
		fmt.Println(redB(":: [ERROR]"), color.HiRedString("No file create, but we tried !! #sadface:\n\t->", LinuxAutoStartP, "\n\t->"))
		return errCreate
	}
	defer autoF.Close()

	// Adapt the content to be written according to the mode selected
	if should {
		ctntAutoF = `[Desktop Entry]
Encoding=UTF-8
Type=Application
Name=TEMPest
Comment=Autorun of TEMPest
Exec=` + conf.Gopath + `/src/github.com/ChacaS0/tempest/scripts/autostart.sh
`
	} else {
		// If exists
		// Keep the file but set ``Hidden`` to ``true``
		ctntAutoF = `[Desktop Entry]
Encoding=UTF-8
Type=Application
Name=TEMPest
Comment=Autorun of TEMPest
Exec=` + conf.Gopath + `/src/github.com/ChacaS0/tempest/scripts/autostart.sh
Hidden=true
`
	}

	// then write to it the config
	if _, errW := autoF.WriteString(ctntAutoF); errW != nil {
		fmt.Println(redB(":: [ERROR]"), color.HiRedString("Can't write to the file, WTF!?\n\t->"))
		return errW
	}

	return nil
}
