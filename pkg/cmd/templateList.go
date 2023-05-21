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
	"bytes"
	"fmt"
	"github.com/ugol/jr/pkg/constants"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/functions"
)

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long:  `List all available templates`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println()
		fmt.Println("List of available JR templates:")
		fmt.Println()
		templateDir, _ := cmd.Flags().GetString("templateDir")
		templateDir = os.ExpandEnv(templateDir)

		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			return
		}

		var Red = "\033[31m"
		var Green = "\033[32m"
		var Reset = "\033[0m"

		err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, "tpl") {
				fullPath, _ := cmd.Flags().GetBool("fullPath")

				t, _ := os.ReadFile(path)
				valid, err := isValidTemplate(t)
				if valid {
					fmt.Print(Green)
				} else {
					fmt.Print(Red)
				}

				if fullPath {
					fmt.Println(path)
				} else {
					name, _ := strings.CutSuffix(info.Name(), ".tpl")
					fmt.Print(name)
					if err != nil {
						fmt.Println(" -> ", err)
					} else {
						fmt.Println()
					}
				}
			}
			return nil
		})
		fmt.Println(Reset)

		if err != nil {
			fmt.Printf("Error in %q: %v\n", templateDir, err)
			return
		}
	},
}

func isValidTemplate(t []byte) (bool, error) {

	tt, err := template.New("test").Funcs(functions.FunctionsMap()).Parse(string(t))
	if err != nil {
		return false, err
	}

	var buf bytes.Buffer
	if err = tt.Execute(&buf, nil); err != nil {
		return false, err
	}

	return true, err

}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateListCmd.Flags().String("templateDir", constants.TEMPLATEDIR, "directory containing templates")
	templateListCmd.Flags().BoolP("fullPath", "f", false, "Print full path")
}
