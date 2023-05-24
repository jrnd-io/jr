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
	"github.com/ugol/jr/pkg/ctx"
	"github.com/ugol/jr/pkg/functions"
	"github.com/ugol/jr/pkg/tpl"
	"log"
	"os"
	"sync"
	"time"
)

var emitterRunCmd = &cobra.Command{
	Use:   "run",
	Short: "RunPreload all or selected configured emitters",
	Long:  `RunPreload all or selected configured emitters`,
	Run: func(cmd *cobra.Command, args []string) {

		err := viper.UnmarshalKey("emitters", &emitters)
		if err != nil {
			log.Println(err)
		}

		if len(args) == 0 {
			for i := 0; i < len(emitters); i++ {
				emitters[i].Initialize(GlobalCfg)
				emitters[i].RunPreload(GlobalCfg)
			}
		} else {
			for i := 0; i < len(emitters); i++ {
				if functions.Contains(args, emitters[i].Name) {
					emitters[i].Initialize(GlobalCfg)
					emitters[i].RunPreload(GlobalCfg)
				}
			}
		}

		doLoop()

	},
}

func doLoop() {
	numTimers := len(emitters)
	timers := make([]*time.Timer, numTimers)
	stopChannels := make([]chan struct{}, numTimers)

	var wg sync.WaitGroup
	wg.Add(numTimers)

	for i := 0; i < numTimers; i++ {
		index := i

		stopChannels[i] = make(chan struct{})

		go func(timerIndex int) {
			defer wg.Done()

			frequency := emitters[timerIndex].Frequency
			if frequency > 0 {
				ticker := time.NewTicker(emitters[timerIndex].Frequency)

				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:

						//fmt.Printf("%s every %s \n", emitters[index].Name, emitters[index].Frequency)

						keyTpl, err := tpl.NewTpl("key", emitters[index].KeyTemplate, functions.FunctionsMap(), &ctx.JrContext)
						if err != nil {
							log.Println(err)
						}
						templatePath := fmt.Sprintf("%s/%s.tpl", os.ExpandEnv(GlobalCfg.TemplateDir), emitters[index].ValueTemplate)
						vt, err := os.ReadFile(templatePath)
						valueTpl, err := tpl.NewTpl("value", string(vt), functions.FunctionsMap(), &ctx.JrContext)
						if err != nil {
							log.Println(err)
						}

						k := keyTpl.Execute()
						v := valueTpl.Execute()
						emitters[index].Producer.Produce([]byte(k), []byte(v), nil)

					case <-stopChannels[timerIndex]:
						return
					}
				}
			}
		}(index)

		timers[i] = time.AfterFunc(emitters[index].Duration, func() {
			stopChannels[index] <- struct{}{}
		})
	}

	wg.Wait()
}

func init() {
	emitterCmd.AddCommand(emitterRunCmd)
}
