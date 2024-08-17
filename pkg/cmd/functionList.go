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
	"github.com/ugol/jr/pkg/functions"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var functionListCmd = &cobra.Command{
	Use:   "list",
	Short: "describes available functions",
	Long: "describes available functions. Example usage:\n" +
		"jr function list lorem",
	Run: func(cmd *cobra.Command, args []string) {
		doList(cmd, args)
	},
}

var functionManCmd = &cobra.Command{
	Use:   "man",
	Short: "describes available functions",
	Long: "describes available functions. Example usage:\n" +
		"jr man lorem",
	Run: func(cmd *cobra.Command, args []string) {
		doList(cmd, args)
	},
}

func doList(cmd *cobra.Command, args []string) {
	category, _ := cmd.Flags().GetBool("category")
	find, _ := cmd.Flags().GetBool("find")
	run, _ := cmd.Flags().GetBool("run")
	isMarkdown, _ := cmd.Flags().GetBool("markdown")
	noColor, _ := cmd.Flags().GetBool("nocolor")

	if category && len(args) > 0 {
		var functionNames []string
		for k, v := range functions.DescriptionMap() {
			if strings.Contains(v.Category, args[0]) {
				functionNames = append(functionNames, k)
			}
		}
		sortAndPrint(functionNames, isMarkdown, noColor)
	} else if find && len(args) > 0 {
		var functionNames []string
		for k, v := range functions.DescriptionMap() {
			if strings.Contains(v.Description, args[0]) || strings.Contains(v.Name, args[0]) {
				functionNames = append(functionNames, k)
			}
		}
		sortAndPrint(functionNames, isMarkdown, noColor)
	} else if len(args) == 1 {

		if run {
			f, found := printFunction(args[0], isMarkdown, noColor)
			if found {
				fmt.Println()
				cmd := exec.Command("/bin/sh", "-c", f.Example)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		} else {
			printFunction(args[0], isMarkdown, noColor)
		}
	} else {

		//l := len(functions.FunctionsMap())
		l := len(functions.DescriptionMap())
		functionNames := make([]string, l)

		i := 0
		for k := range functions.DescriptionMap() {
			functionNames[i] = k
			i++
		}
		sortAndPrint(functionNames, isMarkdown, noColor)

	}
	fmt.Println()
}

func sortAndPrint(functionNames []string, isMarkdown bool, noColor bool) {
	slices.Sort(functionNames)
	for _, k := range functionNames {
		printFunction(k, isMarkdown, noColor)
		//fmt.Println(k)
	}
	fmt.Println()
	fmt.Printf("Total functions: %d\n", len(functionNames))
}

func printFunction(name string, isMarkdown bool, noColor bool) (functions.FunctionDescription, bool) {
	f, found := functions.Description(name)

	var Cyan = ""
	var Reset = ""
	if !noColor {
		Cyan = "\033[36m"
		Reset = "\033[0m"
	}

	if found {
		if isMarkdown {
			fmt.Println()
			fmt.Printf("## Name: %s \n", f.Name)
			fmt.Printf("**Category:** %s\\\n", f.Category)
			fmt.Printf("**Description:** %s\\\n", f.Description)

			if len(f.Parameters) > 0 {
				fmt.Printf("**Parameters:** `%s`\\\n", f.Parameters)
			} else {
				fmt.Printf("**Parameters:** %s \\\n", f.Parameters)
			}
			fmt.Printf("**Localizable:** `%v`\\\n", f.Localizable)
			fmt.Printf("**Return:** `%s`\\\n", f.Return)
			fmt.Printf("**Example:** `%s`\\\n", f.Example)
			fmt.Printf("**Output:** `%s`\n", f.Output)
		} else {
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
	return f, found

}

func init() {
	rootCmd.AddCommand(functionManCmd)
	functionCmd.AddCommand(functionListCmd)
	functionListCmd.Flags().BoolP("category", "c", false, "IndexOf in category")
	functionListCmd.Flags().BoolP("find", "f", false, "IndexOf in description and name")
	functionListCmd.Flags().BoolP("markdown", "m", false, "Output the list as markdown")
	functionListCmd.Flags().BoolP("run", "r", false, "Run the example")
	functionListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")

	functionManCmd.Flags().BoolP("category", "c", false, "IndexOf in category")
	functionManCmd.Flags().BoolP("markdown", "m", false, "Output the list as markdown")
	functionManCmd.Flags().BoolP("find", "f", false, "IndexOf in description and name")
	functionManCmd.Flags().BoolP("run", "r", false, "Run the example")
	functionManCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")

}
