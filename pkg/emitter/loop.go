// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package emitter

import (
	"context"
	"fmt"
	"github.com/jrnd-io/jr/pkg/configuration"
	"github.com/jrnd-io/jr/pkg/ctx"
	"github.com/jrnd-io/jr/pkg/functions"
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

func Initialize(emitterNames []string, es map[string][]Emitter, dryrun bool) []Emitter {

	runAll := len(emitterNames) == 0
	emittersToRun := make([]Emitter, 0, len(es))

	if runAll {
		for _, emitters := range es {
			emittersToRun = InitializeEmitters(emitters, dryrun, emittersToRun)
		}
	} else {
		for _, name := range emitterNames {
			emitters, enabled := es[name]
			if enabled {
				emittersToRun = InitializeEmitters(emitters, dryrun, emittersToRun)
			}
		}
	}

	return emittersToRun
}

func InitializeEmitters(emitters []Emitter, dryrun bool, emittersToRun []Emitter) []Emitter {
	for i := 0; i < len(emitters); i++ {
		if dryrun {
			emitters[i].Output = "stdout"
		}
		emitters[i].Initialize(configuration.GlobalCfg)
		emittersToRun = append(emittersToRun, emitters[i])
		emitters[i].Run(emitters[i].Preload, nil)
	}
	return emittersToRun
}

func DoLoop(es []Emitter) {

	for _, e := range es {
		addEmitterToExpectedObjects(e)
	}
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
		kInValue := functions.GetV("KEY")

		if (kInValue) != "" {
			emitter.Producer.Produce([]byte(kInValue), []byte(v), nil)
		} else {
			emitter.Producer.Produce([]byte(k), []byte(v), nil)
		}

		ctx.JrContext.GeneratedObjects++
		ctx.JrContext.GeneratedBytes += int64(len(v))
	}

}

func CloseProducers(es map[string][]Emitter) {
	for _, v := range es {
		for i := 0; i < len(v); i++ {
			p := v[i].Producer
			if p != nil {
				if err := p.Close(); err != nil {
					fmt.Printf("Error in closing producers: %v\n", err)
				}
			}
		}
	}
	time.Sleep(100 * time.Millisecond)
}

func addEmitterToExpectedObjects(e Emitter) {
	d := e.Duration.Milliseconds()
	f := e.Frequency.Milliseconds()
	n := e.Num
	//fmt.Printf("%d %d %d\n", d, f, n)

	if d > 0 && f > 0 && n > 0 {
		expected := (d / f) * int64(n)
		ctx.JrContext.ExpectedObjects += expected
	}
}

func WriteStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(ctx.JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", ctx.JrContext.GeneratedObjects)

	ungenerated := ctx.JrContext.ExpectedObjects - ctx.JrContext.GeneratedObjects
	if ungenerated > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Data NOT Generated (Objects): %d\n", ungenerated)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", ctx.JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(ctx.JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}
