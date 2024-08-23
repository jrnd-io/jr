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
	jrctx "github.com/jrnd-io/jr/pkg/ctx"
	"github.com/jrnd-io/jr/pkg/functions"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type Producer interface {
	Produce(ctx context.Context, key []byte, val []byte, o any)
	Close(ctx context.Context) error
}

func Initialize(ctx context.Context, emitterNames []string, es map[string][]Emitter, dryrun bool) []Emitter {

	runAll := len(emitterNames) == 0
	emittersToRun := make([]Emitter, 0, len(es))

	if runAll {
		for _, emitters := range es {
			emittersToRun = InitializeEmitters(ctx, emitters, dryrun, emittersToRun)
		}
	} else {
		for _, name := range emitterNames {
			emitters, enabled := es[name]
			if enabled {
				emittersToRun = InitializeEmitters(ctx, emitters, dryrun, emittersToRun)
			}
		}
	}

	return emittersToRun
}

func InitializeEmitters(ctx context.Context, emitters []Emitter, dryrun bool, emittersToRun []Emitter) []Emitter {
	for i := 0; i < len(emitters); i++ {
		if dryrun {
			emitters[i].Output = "stdout"
		}
		emitters[i].Initialize(ctx, configuration.GlobalCfg)
		emittersToRun = append(emittersToRun, emitters[i])
		emitters[i].Run(ctx, emitters[i].Preload, nil)
	}
	return emittersToRun
}

func DoLoop(ctx context.Context, es []Emitter) {

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
						doTemplate(ctx, es[index])
					case <-stopChannels[timerIndex]:
						return
					}

				}
			} else {
				doTemplate(ctx, es[index])
			}
		}(index)

		timers[i] = time.AfterFunc(es[index].Duration, func() {
			stopChannels[index] <- struct{}{}
		})
	}

	wg.Wait()
}

func doTemplate(ctx context.Context, emitter Emitter) {
	jrctx.JrContext.Locale = emitter.Locale
	jrctx.JrContext.CountryIndex = functions.IndexOf(strings.ToUpper(emitter.Locale), "country")

	for i := 0; i < emitter.Num; i++ {
		jrctx.JrContext.CurrentIterationLoopIndex++

		k := emitter.KTpl.Execute()
		v := emitter.VTpl.Execute()
		if emitter.Oneline {
			v = strings.ReplaceAll(v, "\n", "")
		}
		kInValue := functions.GetV("KEY")

		if (kInValue) != "" {
			emitter.Producer.Produce(ctx, []byte(kInValue), []byte(v), nil)
		} else {
			emitter.Producer.Produce(ctx, []byte(k), []byte(v), nil)
		}

		jrctx.JrContext.GeneratedObjects++
		jrctx.JrContext.GeneratedBytes += int64(len(v))
	}

}

func CloseProducers(ctx context.Context, es map[string][]Emitter) {
	for _, v := range es {
		for i := 0; i < len(v); i++ {
			p := v[i].Producer
			if p != nil {
				if err := p.Close(ctx); err != nil {
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
	// fmt.Printf("%d %d %d\n", d, f, n)

	if d > 0 && f > 0 && n > 0 {
		expected := (d / f) * int64(n)
		jrctx.JrContext.ExpectedObjects += expected
	}
}

func WriteStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(jrctx.JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", jrctx.JrContext.GeneratedObjects)

	ungenerated := jrctx.JrContext.ExpectedObjects - jrctx.JrContext.GeneratedObjects
	if ungenerated > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Data NOT Generated (Objects): %d\n", ungenerated)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", jrctx.JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(jrctx.JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}
