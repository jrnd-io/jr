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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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

		err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, "tpl") {
				fmt.Println(path)
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

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("templateDir", "$HOME/.jr/templates", "directory containing templates")
}
