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

package cmd

import (
	"context"
	"github.com/jrnd-io/jr/pkg/emitter"
	"github.com/spf13/cobra"
)

var emitterRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all or selected configured emitters",
	Long:  `Run all or selected configured emitters`,
	Run: func(cmd *cobra.Command, args []string) {

		dryrun, _ := cmd.Flags().GetBool("dryrun")
		RunEmitters(cmd.Context(), args, emitters2, dryrun)

	},
}

func RunEmitters(ctx context.Context, emitterNames []string, ems map[string][]emitter.Emitter, dryrun bool) {
	defer emitter.WriteStats()
	defer emitter.CloseProducers(ctx, ems)
	emittersToRun := emitter.Initialize(ctx, emitterNames, ems, dryrun)
	emitter.DoLoop(ctx, emittersToRun)
}

func init() {
	emitterCmd.AddCommand(emitterRunCmd)
	emitterRunCmd.Flags().BoolP("dryrun", "d", false, "dryrun: output of the emitters to stdout")
}
