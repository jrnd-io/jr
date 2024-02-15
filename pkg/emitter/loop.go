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

package emitter

import (
	"context"
	"fmt"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/ctx"
	"github.com/ugol/jr/pkg/functions"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Producer interface {
	Produce(k []byte, v []byte, o any)
	io.Closer
}

func Initialize(emitterNames []string, es []Emitter, dryrun bool) []Emitter {

	howManyEmitters := len(emitterNames)
	if howManyEmitters == 0 {
		for i := 0; i < len(es); i++ {
			if dryrun {
				es[i].Output = "stdout"
			}
			es[i].Initialize(configuration.GlobalCfg)
			es[i].Run(es[i].Preload, nil)
		}
		return es
	} else {
		emittersToRun := make([]Emitter, 0, howManyEmitters)
		for i := 0; i < len(es); i++ {
			if functions.Contains(emitterNames, es[i].Name) {
				if dryrun {
					es[i].Output = "stdout"
				}
				es[i].Initialize(configuration.GlobalCfg)
				emittersToRun = append(emittersToRun, es[i])
				es[i].Run(es[i].Preload, nil)
			}
		}
		return emittersToRun
	}
}

func DoLoop(es []Emitter) {
	numTimers := len(es)
	timers := make([]*time.Timer, numTimers)
	stopChannels := make([]chan struct{}, numTimers)

	var wg sync.WaitGroup
	wg.Add(numTimers)

	controlC, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i := 0; i < numTimers; i++ {

		index := i

		stopChannels[i] = make(chan struct{})

		go func(timerIndex int) {
			defer wg.Done()

			frequency := es[timerIndex].Frequency
			if frequency > 0 {
				ticker := time.NewTicker(es[timerIndex].Frequency)
				defer ticker.Stop()
				for {
					select {
					case <-controlC.Done():
						stop()
						return
					case <-ticker.C:
						doTemplate(es[index])
					case <-stopChannels[timerIndex]:
						return
					}

				}
			} else {
				doTemplate(es[index])
			}
		}(index)

		timers[i] = time.AfterFunc(es[index].Duration, func() {
			stopChannels[index] <- struct{}{}
		})
	}

	wg.Wait()
}

func doTemplate(emitter Emitter) {
	ctx.JrContext.Locale = emitter.Locale
	ctx.JrContext.CountryIndex = functions.IndexOf(strings.ToUpper(emitter.Locale), "country")

	for i := 0; i < emitter.Num; i++ {
		ctx.JrContext.CurrentIterationLoopIndex++

		k := emitter.KTpl.Execute()
		v := emitter.VTpl.Execute()
		if emitter.Oneline {
			v = strings.ReplaceAll(v, "\n", "")
		}
		emitter.Producer.Produce([]byte(k), []byte(v), nil)

		ctx.JrContext.GeneratedObjects++
		ctx.JrContext.GeneratedBytes += int64(len(v))
	}
}

func CloseProducers(es []Emitter) {
	for i := 0; i < len(es); i++ {
		p := es[i].Producer
		if p != nil {
			if err := p.Close(); err != nil {
				fmt.Printf("Error in closing producers: %v\n", err)
			}
		}
	}
	time.Sleep(100 * time.Millisecond)
}

func WriteStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(ctx.JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", ctx.JrContext.GeneratedObjects)
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", ctx.JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(ctx.JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}
