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
	"log"
	"os"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a template",
	Long:  `Show a template. Templates must be in templates directory, which by default is in '$HOME/.jr/templates'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("Template missing. Try the list command to see available templates")
			os.Exit(1)
		}

		templateDir, _ := cmd.Flags().GetString("templateDir")
		templateDir = os.ExpandEnv(templateDir)
		templateName := fmt.Sprintf("%s/%s.json", templateDir, args[0])
		templateScript, err := os.ReadFile(templateName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(templateScript))

	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().String("templateDir", "$HOME/.jr/templates", "directory containing templates")

}
