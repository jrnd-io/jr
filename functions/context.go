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

package functions

import (
	"os"
	"time"
)

var JrContext Context

const NUM = 1
const FREQUENCY = 0
const DURATION = 0
const TEMPLATEDIR = "$HOME/.jr/templates"

type Context struct {
	StartTime        time.Time
	TemplateDir      string
	TemplateType     string
	GeneratedObjects int64
	GeneratedBytes   int64
	Num              int
	Range            []int
	Frequency        time.Duration
	Duration         time.Duration
	Locales          []string
	Seed             int64
	Counters         map[string]int
}

func init() {

	JrContext = Context{
		StartTime:        time.Now(),
		TemplateDir:      os.ExpandEnv(TEMPLATEDIR),
		TemplateType:     "",
		GeneratedBytes:   0,
		GeneratedObjects: 0,
		Num:              NUM,
		Range:            make([]int, NUM),
		Frequency:        FREQUENCY,
		Duration:         DURATION,
		Locales:          []string{"us"},
		Seed:             time.Now().UTC().UnixNano(),
		Counters:         make(map[string]int),
	}
}
