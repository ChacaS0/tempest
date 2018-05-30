// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// slTargets tells if we are creating new targets from scratch
var slTargets bool

// autoGen tells if the directory name should be generated
var autoGen bool

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if slTargets {
			newTarget(args...)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	RootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().BoolVarP(&slTargets, "target", "t", false, "Create the directory and registers the fresh created directory as a target in TEMPest")
	newCmd.Flags().BoolVarP(&autoGen, "autoGen", "a", false, "Combine this flag with --target ``<path>...`` in order to generate a default name for the target")
}

// newTargets add a list of targets, create and register
func newTarget(toAdd ...string) error {
	// Fetch current directory
	workingDir, errDir := os.Getwd()
	if errDir != nil {
		// log.Fatal(errDir)
		fmt.Println(redB(":: [ERROR]"), color.HiRedString("This is not possible at the moment. Resulted with error\n\t-> ", errDir.Error()))
	}

	if len(toAdd) > 0 {
		// In case it is a relative path
		for i, path := range toAdd {
			treatRelativePath(&path, workingDir)
			// `--auto` flag - setup the dir name at default
			if autoGen {
				toAdd[i] = TreatLastChar(path) + Slash + "temp.est"
			}
		}
		// strip existing targets
		toAdd, errStrip := stripExistingTargets(toAdd)
		if errStrip != nil {
			fmt.Println(redB(":: [ERROR]"), color.HiRedString("Something went terribly wrong WTF !?"), toAdd)
			return errStrip
		}

	} else {
		// Create a new default one in the current directory
		toAdd = append(toAdd, workingDir+Slash+"temp.est")
	}

	// create first
	for _, tgt := range toAdd {
		// Check if already exists
		if _, errCheckF := os.Stat(tgt); errCheckF != nil {
			// If doesn't exist, create
			if err := os.MkdirAll(tgt, 0777); err != nil {
				fmt.Println(redB(":: [ERROR]"), color.HiRedString("Failed to create the directory for this target:"), tgt)
				return err
			}
			fmt.Println(cyanB(":: [INFO]"), color.HiCyanString("DIR CREATED::>\t"), tgt)
		}
	}

	// then register
	if err := addLine(toAdd); err != nil {
		fmt.Println(redB(":: [ERROR]"), color.HiRedString("Can't register those as \"Targets\":\n\t"), toAdd)
		return err
	}

	return nil
}
