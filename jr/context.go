package jr

import (
	"os"
	"time"
)

var JrContext Context

const HOWMANY = 1
const FREQUENCY = 0
const DURATION = 0
const TEMPLATEDIR = "$HOME/.jr/templates"

type Context struct {
	StartTime        time.Time
	TemplateDir      string
	GeneratedObjects int64
	GeneratedBytes   int64
	HowMany          int
	Range            []int
	Frequency        time.Duration
	Duration         time.Duration
	Locales          []string
	Seed             int64
}

func init() {

	JrContext = Context{
		StartTime:        time.Now(),
		TemplateDir:      os.ExpandEnv(TEMPLATEDIR),
		GeneratedBytes:   0,
		GeneratedObjects: 0,
		HowMany:          HOWMANY,
		Range:            make([]int, 1),
		Frequency:        FREQUENCY,
		Duration:         DURATION,
		Locales:          []string{"us"},
		Seed:             time.Now().UTC().UnixNano(),
	}
}
