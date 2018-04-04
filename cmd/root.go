// Package cmd contains all the commands for TEMPest.
//
// Copyright © 2018 Sebastien Bastide
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
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vrischmann/envconfig"
)

var cfgFile string

var conf struct {
	Home   string
	Gobin  string
	Gopath string
}

// type vConfig struct {
// 	duration int `yaml:"duration"`
// }

// TempestConfigDir points to the config directory of TEMPest
// it holds pretty much all configuration for TEMPest
var TempestConfigDir string

// pathProg is the path to the git root
var pathProg string

// pathTempest is the path to the tempest folder
var pathTempest string

// LogShutup is the path to the log of the 'shutup mode'
var LogShutup string

// Tempestcf is the path to the .tempestcf file
// this file holds all the paths (targets) of TEMPest
var Tempestcf string

// Tempestyml is the path to the .tempest.yaml file
// which olds all the the config for the TEMPest tool
var Tempestyml string

// TempestymlDef is the default path to the .tempest.yaml file
// use it wisely.
var TempestymlDef string

// isVersion is the flag variable that indicates whether we want to see the version
var isVersion bool

//* Bold Colors
// blueB is a func used to print in bold blue
var blueB func(...interface{}) string

// yellowB is a func used to print in bold yellow
var yellowB func(...interface{}) string

// whiteB is a func used to print in bold yellow
var whiteB func(...interface{}) string

// redB is a func used to print in bold red
var redB func(...interface{}) string

// greenB is a func used to print in bold red
var greenB func(...interface{}) string

// magB is a func used to print in bold magenta
var magB func(...interface{}) string

// Target is represented by an index and a path
// Later this will hold the type(directory or file)
type Target struct {
	Index int
	Path  string
}

// RootCmd represents the base command when called without any subcommands
//TODO Make full description with full help on how to use the CLI
var RootCmd = &cobra.Command{
	Use:   "tempest",
	Short: "TEMPest is a simple CLI to manage temporary directories.",
	Long: `TEMPest is a simple CLI to manage temporary directories.
It is still under development, so it's normal if it's not perfect .. YET!
You can start by checking if the config file exists at:
	~/.tempest.yaml
	It contains the files' contraint of age (duration in days).

Then you can initialize the list of directories handled by TEMPest. For example:
	tempest init
Then change directory (cd) to a directory you desire to add, and run:
	tempest add
Or just specify the path to the directory (you can add multiple). For example:
	tempest add /tmp /temp

# Note that, by convention, the tempory directories will be called 'temp'

To start cleaning temp directories just run:
	tempest start
Or if you want to see what files/folders would get deleted:
	tempest start -t
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// According to the flag
		switch {
		case isVersion:
			// If isVersion is true we can display the current version
			if version, errVersion := getVersion(); errVersion != nil {
				color.Red(errVersion.Error())
			} else {
				fmt.Println(color.HiYellowString(version))
			}
		default:
			// By default we print help
			if errHelp := cmd.Help(); errHelp != nil {
				color.Red(errHelp.Error())
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize the environment variables

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf.Home = home

	cobra.OnInitialize(initConfig)

	if err := envconfig.Init(&conf); err != nil {
		// log.Println(err)
	}

	TempestConfigDir = conf.Home + string(os.PathSeparator) + ".tempest"

	pathProg = conf.Gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "github.com" + string(os.PathSeparator) + "ChacaS0" + string(os.PathSeparator)
	pathTempest = pathProg + "tempest" + string(os.PathSeparator)

	Tempestcf = TempestConfigDir + string(os.PathSeparator) + ".tempestcf"
	Tempestyml = viper.ConfigFileUsed()
	TempestymlDef = TempestConfigDir + string(os.PathSeparator) + ".tempest.yaml"
	LogShutup = TempestConfigDir + string(os.PathSeparator) + ".log" + "shutup.log"

	//* Bold Colors
	yellowB = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	blueB = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	whiteB = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	redB = color.New(color.FgHiRed, color.Bold).SprintFunc()
	greenB = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	magB = color.New(color.FgHiMagenta, color.Bold).SprintFunc()

	// //* Man ?
	// header := &doc.GenManHeader{
	// 	Title:   "List",
	// 	Section: "3",
	// }
	// out := new(bytes.Buffer)

	// doc.GenMan(RootCmd, header, out)
	// fmt.Println(out.String())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", viper.ConfigFileUsed(), "config file (default is $HOME/.tempest/.tempest.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().BoolVarP(&isVersion, "version", "v", false, "Display the current version v[VERSION_NUMBER]-X-Y[REVISION_NUMBER]")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(TempestymlDef)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conf.Home = home

		// Search config in home directory with name ".tempest" (without extension).
		viper.AddConfigPath(TempestConfigDir)
		viper.SetConfigName(".tempest")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println(blueB(":: Using config file:"), Tempestyml)
	}

	Tempestyml = viper.ConfigFileUsed()
	viper.SetDefault("duration", 5)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// fmt.Println("Config file changed:", e.Name)
		fmt.Println(blueB("::"), color.HiBlueString("Config file changed:"), e.Name)
		Tempestyml = viper.ConfigFileUsed()
	})
}

// Round just does what it says it does
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// openBrowser Opens an url inside a browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// IsDirectory returns true if this path points to a directory
// If there is an error, the func will return it
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// IsIntInSlice returns true if the int is in the slice
func IsIntInSlice(index int, sl []int) bool {
	if len(sl) <= 0 || sl == nil {
		return false
	}

	for _, val := range sl {
		if val == index {
			return true
		}
	}
	return false
}

// IsStringInSlice returns true if the string is in the slice.
// This one is more simple and faster than checkRedondance(slice, sliceArgs[]string) bool
func IsStringInSlice(str string, sl []string) bool {
	if len(sl) <= 0 || sl == nil {
		return false
	}

	for _, val := range sl {
		if val == str {
			return true
		}
	}

	return false
}

// PathsToTargets is a converter, takes paths (strings) and convert them into targets (Target)
func PathsToTargets(paths []string) []Target {
	sliceTgt := make([]Target, 0)

	for i, p := range paths {
		sliceTgt = append(sliceTgt, Target{i, p})
	}

	return sliceTgt
}

// captureStdout returns the output of a function
// not thread safe
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// WriteLog write the strings passed in param into the log file pointed
func WriteLog(pathLog string, strs ...string) {
	// open file first - if does not exist, create it biatch
	var f *os.File
	// TODO replace this path by the var once merged with the non-temp branch
	f, errF := os.OpenFile(pathLog, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		f2, errF2 := os.OpenFile(pathLog, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if errF2 != nil {
			fmt.Println(redB("::"), color.HiRedString("Could not open the file\n\t->"), errF2)
		}
		defer f2.Close()
		f = f2
	}
	defer f.Close()

	// styling
	currTime := time.Now().String()
	header := "=========================  - [" + currTime + "] -  ========================="
	footer := "=================================================================================================================="

	// writing logs
	for _, str := range strs {
		// Write it for each str passed in param
		toWrite := header + "\n" + str + "\n" + footer + "\n"
		if _, err := f.WriteString(toWrite); err != nil {
			fmt.Println(redB(":: [ERROR]"), color.HiRedString("Sorry could not write logs\n\t->"), err)
		}
	}
}
