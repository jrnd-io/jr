// Copyright © 2024 JR team
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
)

var emitterShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show a configured emitter",
	Long: `show a configured emitter: 
jr emitter show net_device`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		noColor, _ := cmd.Flags().GetBool("nocolor")

		var Green = ""
		var Reset = ""
		if !noColor {
			Green = "\033[32m"
			Reset = "\033[0m"
		}

		fmt.Println()
		for k, v := range emitters2 {
			if k == args[0] {
				for _, e := range v {
					fmt.Printf("%sName:%s%s\n", Green, Reset, e.Name)
					fmt.Printf("%sLocale: %s%s\n", Green, Reset, e.Locale)
					fmt.Printf("%sNum: %s%d\n", Green, Reset, e.Num)
					fmt.Printf("%sFrequency: %s%s\n", Green, Reset, e.Frequency)
					fmt.Printf("%sDuration: %s%s\n", Green, Reset, e.Duration)
					fmt.Printf("%sPreload: %s%d\n", Green, Reset, e.Preload)
					fmt.Printf("%sOutput: %s%s\n", Green, Reset, e.Output)
					fmt.Printf("%sTopic: %s%s\n", Green, Reset, e.Topic)
					fmt.Printf("%sKcat: %s%v\n", Green, Reset, e.Kcat)
					fmt.Printf("%sOneline: %s%v\n", Green, Reset, e.Oneline)
					fmt.Printf("%sKey Template: %s%s\n", Green, Reset, e.KeyTemplate)
					fmt.Printf("%sValue Template: %s%s\n", Green, Reset, e.ValueTemplate)
					fmt.Printf("%sOutput Template: %s%s\n", Green, Reset, e.OutputTemplate)
				}
			}
		}
		fmt.Println()

	},
}

func init() {
	emitterCmd.AddCommand(emitterShowCmd)
	emitterShowCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")

}
