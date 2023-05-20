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
	"github.com/spf13/viper"
	"github.com/ugol/jr/pkg/functions"
	"log"
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
var topic string

var emitterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an emitter",
	Long:  `Create an emitter`,
	Run: func(cmd *cobra.Command, args []string) {

		emitter := Emitter{
			Name:          name,
			Locale:        locale,
			Num:           num,
			Frequency:     frequency,
			Duration:      duration,
			Preload:       preload,
			ValueTemplate: valueTemplate,
			KeyTemplate:   keyTemplate,
			Topic:         topic,
		}

		err := viper.UnmarshalKey("emitters", &emitters)
		if err != nil {
			log.Println(err)
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
	emitterCmd.AddCommand(emitterCreateCmd)
	emitterCreateCmd.Flags().IntVarP(&num, "num", "n", functions.NUM, "Number of elements to create for each pass")
	emitterCreateCmd.Flags().StringVar(&name, "name", functions.DEFAULT_EMITTER_NAME, "Emitter name")
	emitterCreateCmd.Flags().DurationVarP(&frequency, "frequency", "f", functions.FREQUENCY, "how much time to wait for next generation pass")
	emitterCreateCmd.Flags().DurationVarP(&duration, "duration", "d", functions.DURATION, "If frequency is enabled, with Duration you can set a finite amount of time")
	emitterCreateCmd.Flags().IntVarP(&preload, "preload", "p", functions.DEFAULT_PRELOAD_SIZE, "Default number of elements to create in preload phase")
	emitterCreateCmd.Flags().StringVar(&locale, "locale", functions.LOCALE, "Locale")
	emitterCreateCmd.Flags().StringVar(&valueTemplate, "valueTemplate", functions.DEFAULT_VALUE_TEMPLATE, "Locale")
	emitterCreateCmd.Flags().StringVar(&keyTemplate, "keyTemplate", functions.DEFAULT_KEY, "Locale")
	emitterCreateCmd.Flags().StringVar(&topic, "topic", functions.DEFAULT_TOPIC, "Locale")
}
