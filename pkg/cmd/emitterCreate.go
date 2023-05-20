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
	"github.com/ugol/jr/pkg/functions"
)

var emitterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an emitter",
	Long:  `Create an emitter`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println()
		fmt.Println("Emitter created:")
		fmt.Println()

	},
}

func init() {
	emitterCmd.AddCommand(emitterCreateCmd)
	emitterCreateCmd.Flags().IntP("num", "n", functions.NUM, "Number of elements to create for each pass")
	emitterCreateCmd.Flags().DurationP("frequency", "f", functions.FREQUENCY, "how much time to wait for next generation pass")
	emitterCreateCmd.Flags().DurationP("duration", "d", functions.DURATION, "If frequency is enabled, with Duration you can set a finite amount of time")
	emitterCreateCmd.Flags().IntP("preload", "p", functions.DEFAULT_PRELOAD_SIZE, "Default number of elements to create in preload phase")
}
