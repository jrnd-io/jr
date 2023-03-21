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
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"jr/jr"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var listCmd = &cobra.Command{
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

		err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, "tpl") {
				fullPath, _ := cmd.Flags().GetBool("fullPath")

				t, _ := os.ReadFile(path)
				if isValidTemplate(t) {
					fmt.Print(Green)
				} else {
					fmt.Print(Red)
				}

				if fullPath {
					fmt.Println(path)
				} else {
					name, _ := strings.CutSuffix(info.Name(), ".tpl")
					fmt.Println(name)
				}
			}
			return nil
		})
		fmt.Println()

		if err != nil {
			fmt.Printf("Error in %q: %v\n", templateDir, err)
			return
		}
	},
}

func isValidTemplate(t []byte) bool {
	jr.JrContext = jr.NewContext(time.Now(), 1, -1, []string{"us"}, 0)

	tt, err := template.New("test").Funcs(jr.FunctionsMap()).Parse(string(t))
	if err != nil {
		return false
	}

	var buf bytes.Buffer
	if err = tt.Execute(&buf, nil); err != nil {
		return false
	}

	return true

}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("templateDir", "$HOME/.jr/templates", "directory containing templates")
	listCmd.Flags().BoolP("fullPath", "f", false, "Print full path")
}
