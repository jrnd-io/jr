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
			for _, v := range emitters {
				v.RunPreload(GlobalCfg)
			}
		} else {
			for _, v := range emitters {
				if functions.Contains(args, v.Name) {
					v.RunPreload(GlobalCfg)
				}
			}
		}

		// Create an array of timers and stop channels
		numTimers := len(emitters)
		timers := make([]*time.Timer, numTimers)
		stopChannels := make([]chan struct{}, numTimers)

		// Create a wait group to ensure all timers are processed
		var wg sync.WaitGroup
		wg.Add(numTimers)

		// Start the timers
		for i := 0; i < numTimers; i++ {
			index := i // Capture the current index for the goroutine

			// Create a stop channel for the current timer
			stopChannels[i] = make(chan struct{})

			// Start the timer's goroutine
			go func(timerIndex int) {
				defer wg.Done()

				// Create a ticker for the current timer's frequency
				frequency := emitters[timerIndex].Frequency
				if frequency > 0 {
					ticker := time.NewTicker(emitters[timerIndex].Frequency)

					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:

							fmt.Printf("%s every %d \n", emitters[index].Name, emitters[index].Frequency)

						case <-stopChannels[timerIndex]:
							return
						}
					}
				}
			}(index)

			// Create a timer to stop the goroutine after the duration elapses
			timers[i] = time.AfterFunc(emitters[index].Duration, func() {
				stopChannels[index] <- struct{}{}
			})
		}

		// Wait for all timers to complete
		wg.Wait()

	},
}

func init() {
	emitterCmd.AddCommand(emitterRunCmd)
}
