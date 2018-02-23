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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// flags vars
var age int

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
		switch {
		case age != -1:
			if errAge := setAge(); errAge != nil {
				fmt.Println(errAge)
			}
			// color.HiCyan("Sorry! No working at the moment! -\\('o')/-")
		default:
			color.HiRed("You must provide some flags or whatever! Just do something!")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().IntVarP(&age, "age", "a", -1, "Set the maximum age (in days) of the temp directories' content")
}

// setAge sets the "duration" value in the config file (~/.tempest.yaml)
// duration being the maximum age of the temp's content
func setAge() error {
	if age > 0 {
		// ageStr := fmt.Sprintf("%d", age)
		viper.Set("duration", age)

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
