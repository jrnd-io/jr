package jr

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
