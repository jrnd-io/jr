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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"jr/jr"
	"log"
	"os"
	"os/signal"
	"regexp"
	"text/template"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a template",
	Long:  `Execute a template. Templates must be in templates directory, which by default is in '$HOME/.jr/templates'`,
	Run: func(cmd *cobra.Command, args []string) {

		t, _ := cmd.Flags().GetString("t")
		templatePath, _ := cmd.Flags().GetString("templatePath")

		if len(args) == 0 && len(t) == 0 && len(templatePath) == 0 {
			log.Println("Template missing. Try the list command to see available templates")
			os.Exit(1)
		}

		var templateScript []byte
		var err error

		if len(t) > 0 {
			templateScript = []byte(t)
		} else if len(templatePath) > 0 {
			templatePath = os.ExpandEnv(templatePath)
			templateScript, err = os.ReadFile(templatePath)
		} else {
			templateDir, _ := cmd.Flags().GetString("templateDir")
			templateDir = os.ExpandEnv(templateDir)
			templatePath := fmt.Sprintf("%s/%s.json", templateDir, args[0])
			templateScript, err = os.ReadFile(templatePath)
		}
		if err != nil {
			log.Fatal(err)
		}

		report, err := template.New("json").Funcs(jr.FunctionsMap()).Parse(string(templateScript))
		if err != nil {
			log.Fatal(err)
		}

		howMany, _ := cmd.Flags().GetInt("n")
		frequency, _ := cmd.Flags().GetInt("f")
		seed, _ := cmd.Flags().GetInt64("seed")
		oneline, _ := cmd.Flags().GetBool("oneline")

		jr.Random.Seed(seed)

		c := jr.NewContext(time.Now(), howMany, frequency, []string{"IT"}, seed)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if frequency != -1 {
		Infinite:
			for {
				select {
				case <-time.After(time.Millisecond * time.Duration(frequency)):
					for range c.Range {
						executeTemplate(report, c, oneline)
					}
				case <-ctx.Done():
					stop()
					break Infinite
				}
			}
		} else {
			for range c.Range {
				executeTemplate(report, c, oneline)
			}
		}

		writeStats(c)

	},
}

func writeStats(c *jr.Context) {
	fmt.Fprintln(os.Stderr)
	elapsed := time.Since(c.StartTime)
	fmt.Fprintf(os.Stderr, "Elapsed time: %s\n", elapsed)
	fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", c.GeneratedObjects)
	fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", c.GeneratedBytes)
	//fmt.Fprintf(os.Stderr, "Data Generated (bytes per second): %d\n", c.GeneratedBytes / elapsed )
}

func executeTemplate(report *template.Template, c *jr.Context, oneline bool) {
	var bt bytes.Buffer
	if err := report.Execute(&bt, c); err != nil {
		log.Fatal(err)
	}
	output := bt.String()
	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		output = re.ReplaceAllString(output, "")
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}
	c.GeneratedObjects++
	c.GeneratedBytes += int64(len(output))
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Int("n", 1, "Number of elements to create for each pass")
	runCmd.Flags().Int("f", -1, "Frequency: number of milliseconds to wait for next generation pass")
	runCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")
	runCmd.Flags().Bool("oneline", false, "strips /n from output, for example to be pipelined to tools like kcat")
	runCmd.Flags().String("templateDir", "$HOME/.jr/templates", "directory containing templates")
	runCmd.Flags().String("templatePath", "", "Path to the template file")
	runCmd.Flags().String("t", "", "use a template on the fly")

}
