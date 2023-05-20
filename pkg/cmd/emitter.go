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
	"github.com/spf13/cobra"
	"time"
)

type Emitter struct {
	Name          string        `mapstructure:"name"`
	Locale        string        `mapstructure:"locale"`
	Num           int           `mapstructure:"num"`
	Frequency     time.Duration `mapstructure:"frequency"`
	Duration      time.Duration `mapstructure:"duration"`
	Preload       int           `mapstructure:"preload"`
	ValueTemplate string        `mapstructure:"valueTemplate"`
	KeyTemplate   string        `mapstructure:"keyTemplate"`
	Topic         string        `mapstructure:"topic"`
	//OutputTemplate string
	//Kcat           bool
	//Output         string
	//Oneline        bool
}

var emitterCmd = &cobra.Command{
	Use:     "emitter",
	Short:   "jr Emitter resource",
	Long:    `jr Emitter resource`,
	GroupID: "resource",
}

// var emitters = make(map[string][]Emitter)

var emitters []Emitter

func init() {
	rootCmd.AddCommand(emitterCmd)
}
