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
const NUM_TEMPLATES = 5
const FREQUENCY = 0
const DURATION = 0
const TEMPLATEDIR = "$HOME/.jr/templates"

type Configuration struct {
	TemplateNames    []string
	KeyTemplate      string
	OutputTemplate   string
	EmbeddedTemplate bool
	TemplateFileName bool
	Kcat             bool
	Output           string
	Oneline          bool
	Locale           string
	Num              int
	Frequency        time.Duration
	Duration         time.Duration
	Seed             int64
	KafkaConfig      string
	RegistryConfig   string
	Topic            []string
	Preload          []string
	PreloadSize      []int
	TemplateDir      string
	Autocreate       bool
	SchemaRegistry   bool
	Serializer       string
	RedisTtl         time.Duration
	RedisConfig      string
}

// Context is the object passed on the templates which contains all the needed details.
type Context struct {
	StartTime        time.Time
	TemplateDir      string
	TemplateType     []string
	GeneratedObjects int64
	GeneratedBytes   int64
	Num              int
	NumTemplates     int
	Range            []int
	Frequency        time.Duration
	Duration         time.Duration
	Locale           string
	Seed             int64
	CtxCounters      map[string]int
	Ctx              map[string]string
	CtxList          map[string][]string
	LastIndex        int
	CountryIndex     int
	CityIndex        int
}

func init() {

	JrContext = Context{
		StartTime:        time.Now(),
		TemplateDir:      os.ExpandEnv(TEMPLATEDIR),
		TemplateType:     make([]string, NUM_TEMPLATES),
		GeneratedBytes:   0,
		GeneratedObjects: 0,
		Num:              NUM,
		NumTemplates:     NUM_TEMPLATES,
		Range:            make([]int, NUM),
		Frequency:        FREQUENCY,
		Duration:         DURATION,
		Locale:           "us",
		Seed:             time.Now().UTC().UnixNano(),
		CtxCounters:      make(map[string]int),
		Ctx:              make(map[string]string),
		CtxList:          make(map[string][]string),
		LastIndex:        -1,
		CountryIndex:     232,
		CityIndex:        -1,
	}
}
