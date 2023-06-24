//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

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

		var Green = "\033[32m"
		var Reset = "\033[0m"

		fmt.Println()
		for _, v := range emitters {
			if v.Name == args[0] {
				fmt.Printf("%sName:%s%s\n", Green, Reset, v.Name)
				fmt.Printf("%sLocale: %s%s\n", Green, Reset, v.Locale)
				fmt.Printf("%sNum: %s%d\n", Green, Reset, v.Num)
				fmt.Printf("%sFrequency: %s%s\n", Green, Reset, v.Frequency)
				fmt.Printf("%sDuration: %s%s\n", Green, Reset, v.Duration)
				fmt.Printf("%sPreload: %s%d\n", Green, Reset, v.Preload)
				fmt.Printf("%sOutput: %s%s\n", Green, Reset, v.Output)
				fmt.Printf("%sTopic: %s%s\n", Green, Reset, v.Topic)
				fmt.Printf("%sKcat: %s%v\n", Green, Reset, v.Kcat)
				fmt.Printf("%sOneline: %s%v\n", Green, Reset, v.Oneline)
				fmt.Printf("%sKey Template: %s%s\n", Green, Reset, v.KeyTemplate)
				fmt.Printf("%sValue Template: %s%s\n", Green, Reset, v.ValueTemplate)
				fmt.Printf("%sOutput Template: %s%s\n", Green, Reset, v.OutputTemplate)
			}
		}
		fmt.Println()

	},
}

func init() {
	emitterCmd.AddCommand(emitterShowCmd)
}
