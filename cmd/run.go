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
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template. 
  Without any other flag, [template] is just the name of a template in the templates directory, which by default is in '$HOME/.jr/templates' Example: 
jr run net-device
  With the --template flag, [template] is a string containing a full template. Example:
jr run --template "{{name}}"
 With the -templateFileName flag [template] is a file name with a template. Example:
jr run --templateFileName ~/.jr/templates/net-device.json
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		keyTemplate, _ := cmd.Flags().GetString("key")
		embeddedTemplate, _ := cmd.Flags().GetBool("template")
		templateFileName, _ := cmd.Flags().GetBool("templateFileName")
		silent, _ := cmd.Flags().GetBool("silent")
		kcat, _ := cmd.Flags().GetBool("kcat")
		topic, _ := cmd.Flags().GetString("topic")
		kafkaConfig, _ := cmd.Flags().GetString("kafkaConfig")
		var producer *kafka.Producer
		var valueTemplate []byte
		var err error

		if topic != "" {
			producer, err = jr.Initialize(kafkaConfig)
			if err != nil {
				log.Fatal(err)
			}
		}

		if embeddedTemplate {
			valueTemplate = []byte(args[0])
		} else if templateFileName {
			valueTemplate, err = os.ReadFile(os.ExpandEnv(args[0]))
		} else {
			templateDir, _ := cmd.Flags().GetString("templateDir")
			templateDir = os.ExpandEnv(templateDir)
			templatePath := fmt.Sprintf("%s/%s.json", templateDir, args[0])
			valueTemplate, err = os.ReadFile(templatePath)
		}
		if err != nil {
			log.Fatal(err)
		}

		key, err := template.New("Key").Funcs(jr.FunctionsMap()).Parse(keyTemplate)
		if err != nil {
			log.Fatal(err)
		}

		value, err := template.New("Value").Funcs(jr.FunctionsMap()).Parse(string(valueTemplate))
		if err != nil {
			log.Fatal(err)
		}

		howMany, _ := cmd.Flags().GetInt("num")
		frequency, _ := cmd.Flags().GetDuration("frequency")
		duration, _ := cmd.Flags().GetDuration("duration")
		seed, _ := cmd.Flags().GetInt64("seed")
		oneline, _ := cmd.Flags().GetBool("oneline")

		if kcat {
			oneline = true
		}

		jr.Random.Seed(seed)

		c := jr.NewContext(time.Now(), howMany, frequency, []string{"IT"}, seed)
		infinite := true
		if duration > 0 {
			timer := time.NewTimer(duration)

			go func() {
				<-timer.C
				infinite = false
			}()
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if frequency != 0 {
		Infinite:
			for ok := true; ok; ok = infinite {
				select {
				case <-time.After(frequency):
					for range c.Range {
						k, v, _ := executeTemplate(key, value, c, oneline)
						output(k, v, producer, topic, silent, kcat)
					}

				case <-ctx.Done():
					stop()
					break Infinite
				}
			}
		} else {
			for range c.Range {
				k, v, _ := executeTemplate(key, value, c, oneline)
				output(k, v, producer, topic, silent, kcat)
			}
		}

		if topic != "" {
			jr.Close(producer)
		}

		writeStats(c)

	},
}

func writeStats(c *jr.Context) {
	fmt.Fprintln(os.Stderr)
	elapsed := time.Since(c.StartTime)
	fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", c.GeneratedObjects)
	fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", c.GeneratedBytes)
	fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(c.GeneratedBytes)/elapsed.Seconds())
	fmt.Fprintln(os.Stderr)
}

func executeTemplate(key *template.Template, value *template.Template, c *jr.Context, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, c); err != nil {
		log.Println(err)
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, c); err != nil {
		log.Println(err)
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	c.GeneratedObjects++
	c.GeneratedBytes += int64(len(v))

	return k, v, err
}

func output(key string, value string, p *kafka.Producer, topic string, silent bool, kcat bool) {

	if kcat {
		fmt.Printf("%s,%s\n", key, value)
		fmt.Println()
	} else if !silent {
		fmt.Println(value)
	}
	if topic != "" {
		jr.Produce(p, []byte(key), []byte(value), topic)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("num", "n", 1, "Number of elements to create for each pass")
	runCmd.Flags().DurationP("frequency", "f", 0, "how much time to wait for next generation pass")
	runCmd.Flags().DurationP("duration", "d", 0, "If frequency is enabled, with Duration you can set a finite amount of time")
	runCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")
	runCmd.Flags().BoolP("oneline", "o", false, "strips /n from output, for example to be pipelined to tools like kcat")
	runCmd.Flags().String("templateDir", "$HOME/.jr/templates", "directory containing templates")
	runCmd.Flags().StringP("kafkaConfig", "F", "./kafka/config.properties", "Kafka configuration")
	runCmd.Flags().Bool("templateFileName", false, "If enabled, [template] must be a template file")
	runCmd.Flags().Bool("template", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")
	runCmd.Flags().StringP("key", "k", "key", "A template to generate a key")
	runCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it writes key (--key), value (output of the template) on stdout")
	runCmd.Flags().StringP("topic", "t", "", "Kafka topic name")
	runCmd.Flags().BoolP("silent", "s", false, "No standard output (makes sense only when writing directly to Kafka)")
}
