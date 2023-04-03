/*
Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"jr/jr"
	"strings"

	"github.com/spf13/cobra"
)

// manCmd represents the man command
var manCmd = &cobra.Command{
	Use:   "man [function]",
	Short: "describes an available function",
	Long: "describes an available function. Example usage:\n" +
		"jr man lorem",
	Run: func(cmd *cobra.Command, args []string) {

		list, _ := cmd.Flags().GetBool("list")
		category, _ := cmd.Flags().GetBool("category")
		find, _ := cmd.Flags().GetBool("find")

		if list {
			for k, _ := range jr.FunctionsMap() {
				printFunction(k)
			}
			return
		} else if category {
			for k, v := range jr.DescriptionMap() {
				if strings.Contains(v.Category, args[0]) {
					printFunction(k)
				}
			}
		} else if find {
			for k, v := range jr.DescriptionMap() {
				if strings.Contains(v.Description, args[0]) || strings.Contains(v.Name, args[0]) {
					printFunction(k)
				}
			}
		} else if len(args) == 1 {
			printFunction(args[0])
		} else {
			fmt.Println("Missing parameter and/or flags")
		}
		fmt.Println()
	},
}

func printFunction(name string) {

	var Reset = "\033[0m"
	var Cyan = "\033[36m"

	f, found := jr.Description(name)
	if found {
		fmt.Println()
		fmt.Printf("%sName: %s%s\n", Cyan, Reset, f.Name)
		fmt.Printf("%sCategory: %s%s\n", Cyan, Reset, f.Category)
		fmt.Printf("%sDescription: %s%s\n", Cyan, Reset, f.Description)
		fmt.Printf("%sParameters: %s%s\n", Cyan, Reset, f.Parameters)
		fmt.Printf("%sLocalizable: %s%v\n", Cyan, Reset, f.Localizable)
		fmt.Printf("%sReturn: %s%s\n", Cyan, Reset, f.Return)
		fmt.Printf("%sExample: %s%s\n", Cyan, Reset, f.Example)
		fmt.Printf("%sOutput: %s%s\n", Cyan, Reset, f.Output)
	}

}

func init() {
	rootCmd.AddCommand(manCmd)
	manCmd.Flags().BoolP("list", "l", false, "Show all functions")
	manCmd.Flags().BoolP("category", "c", false, "Search in category")
	manCmd.Flags().BoolP("find", "f", false, "Search in description and name")
}
