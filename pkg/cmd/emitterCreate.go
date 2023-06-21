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
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"time"
)

var name string
var locale string
var num int
var frequency time.Duration
var duration time.Duration
var preload int
var valueTemplate string
var keyTemplate string
var outputTemplate string
var output string
var topic string
var kcat bool
var oneline bool

var emitterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an emitter",
	Long:  `Create an emitter`,
	Run: func(cmd *cobra.Command, args []string) {

		emitter := emitter.Emitter{
			Name:           name,
			Locale:         locale,
			Num:            num,
			Frequency:      frequency,
			Duration:       duration,
			Preload:        preload,
			ValueTemplate:  valueTemplate,
			KeyTemplate:    keyTemplate,
			OutputTemplate: outputTemplate, //@TODO FIX
			Output:         output,
			Topic:          topic,
			Kcat:           kcat,
			Oneline:        oneline,
		}

		emitters = append(emitters, emitter)

		fmt.Println()
		fmt.Printf("Emitter '%s' created.\n", emitter.Name)
		fmt.Println()
		for _, v := range emitters {
			fmt.Println(v)
		}

		// @TODO add the emitter back to the config file
		// the following doesn't work OOTB
		// viper.SafeWriteConfigAs("jrconfig.json")
		// the array must marshaled back to json and added to the config file

	},
}

func init() {
	//@TODO enable when it works
	//emitterCmd.AddCommand(emitterCreateCmd)
	emitterCreateCmd.Flags().IntVarP(&num, "num", "n", constants.NUM, "Number of elements to create for each pass")
	emitterCreateCmd.Flags().StringVar(&name, "name", constants.DEFAULT_EMITTER_NAME, "Emitter name")
	emitterCreateCmd.Flags().DurationVarP(&frequency, "frequency", "f", constants.FREQUENCY, "how much time to wait for next generation pass")
	emitterCreateCmd.Flags().DurationVarP(&duration, "duration", "d", constants.DURATION, "If frequency is enabled, with Duration you can set a finite amount of time")
	emitterCreateCmd.Flags().IntVarP(&preload, "preload", "p", constants.DEFAULT_PRELOAD_SIZE, "number of elements to create in preload phase")
	emitterCreateCmd.Flags().StringVar(&locale, "locale", constants.LOCALE, "Locale")
	emitterCreateCmd.Flags().StringVar(&valueTemplate, "valueTemplate", constants.DEFAULT_VALUE_TEMPLATE, "template name to use for the value")
	emitterCreateCmd.Flags().StringVar(&keyTemplate, "keyTemplate", constants.DEFAULT_KEY, "template to use for the key")
	emitterCreateCmd.Flags().StringVar(&outputTemplate, "outputTemplate", constants.DEFAULT_OUTPUT_TEMPLATE, "Formatting of K,V on standard output")
	emitterCreateCmd.Flags().StringVarP(&output, "output", "o", constants.DEFAULT_OUTPUT, "can be one of stdout, kafka, redis, mongo, elastic, s3")
	emitterCreateCmd.Flags().StringVar(&topic, "topic", constants.DEFAULT_TOPIC, "Default topic to write to if using output='kafka'")
	emitterCreateCmd.Flags().BoolVar(&kcat, "kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	emitterCreateCmd.Flags().BoolVarP(&oneline, "oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
}
