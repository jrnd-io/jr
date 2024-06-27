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

var emitterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available emitters",
	Long:  `List all available emitters`,
	Run: func(cmd *cobra.Command, args []string) {

		noColor, _ := cmd.Flags().GetBool("nocolor")
		all, _ := cmd.Flags().GetBool("full")

		var Green = ""
		var Reset = ""
		if !noColor {
			Green = "\033[32m"
			Reset = "\033[0m"
		}

		fmt.Println()
		fmt.Println("List of JR emitters:")
		fmt.Println()

		for k, v := range emitters2 {
			if all {
				fmt.Printf("%s%s%s", Green, k, Reset)
				fmt.Print(" -> (")
				for _, e := range v {
					fmt.Printf(" %s ", e.Name)
				}
				fmt.Println(")")

			} else {
				fmt.Printf("%s%s%s\n", Green, k, Reset)
			}
		}
		fmt.Println()

	},
}

func init() {
	emitterCmd.AddCommand(emitterListCmd)
	emitterListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
	emitterListCmd.Flags().BoolP("full", "f", false, "Do not color output")

}
