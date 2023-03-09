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
	"jg/jg"
	"log"
	"os"
	"text/template"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "execute a template",
	Long: `Execute a template
								.`,
	Run: func(cmd *cobra.Command, args []string) {

		templateName := fmt.Sprintf("templates/%s.json", args[0])
		templateScript, err := os.ReadFile(templateName)
		if err != nil {
			fmt.Print(err)
		}

		report, err := template.New("json").Funcs(jg.FunctionsMap()).Parse(string(templateScript))
		if err != nil {
			log.Fatal(err)
		}

		howMany, _ := cmd.Flags().GetInt("n")
		frequency, _ := cmd.Flags().GetInt("f")
		seed, _ := cmd.Flags().GetInt64("seed")
		jg.Random.Seed(seed)

		c := jg.NewContext(howMany, frequency, []string{"IT"}, seed)
		if err := report.Execute(os.Stdout, c); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().Int("n", 1, "Number of elements to create for each pass")
	runCmd.Flags().Int("f", 1, "Frequency: elements generated per second")
	runCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")

}
