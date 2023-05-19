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
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/functions"
)

var showCmd = &cobra.Command{
	Use:   "show [template]",
	Short: "Show a template",
	Long:  `Show a template. Templates must be in templates directory, which by default is in '$HOME/.jr/templates'`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("Template missing. Try the list command to see available templates")
			os.Exit(1)
		}

		templateDir, _ := cmd.Flags().GetString("templateDir")
		nocolor, _ := cmd.Flags().GetBool("nocolor")
		templateDir = os.ExpandEnv(templateDir)
		templatePath := fmt.Sprintf("%s/%s.tpl", templateDir, args[0])
		templateScript, err := os.ReadFile(templatePath)
		valid, err := isValidTemplate([]byte(templateScript))
		templateString := string(templateScript)

		var Reset = "\033[0m"
		if runtime.GOOS != "windows" && !nocolor {
			var Cyan = "\033[36m"
			coloredOpeningBracket := fmt.Sprintf("%s%s", Cyan, "{{")
			coloredClosingBracket := fmt.Sprintf("%s%s", "}}", Reset)
			templateString = strings.ReplaceAll(templateString, "{{", coloredOpeningBracket)
			templateString = strings.ReplaceAll(templateString, "}}", coloredClosingBracket)
		}
		fmt.Println(templateString)
		fmt.Print(Reset)
		if !valid {
			log.Println(err)
		}
	},
}

func init() {
	templateCmd.AddCommand(showCmd)
	showCmd.Flags().String("templateDir", functions.TEMPLATEDIR, "directory containing templates")
	showCmd.Flags().BoolP("nocolor", "n", false, "disable colorized output")
}
