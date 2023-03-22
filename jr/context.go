package jr

import "time"

type Context struct {
	StartTime        time.Time
	TemplateDir      string
	GeneratedObjects int64
	GeneratedBytes   int64
	HowMany          int
	Range            []int
	Frequency        time.Duration
	Locales          []string
	Seed             int64
}

func NewContext(startTime time.Time, howMany int, frequency time.Duration, locales []string, seed int64) *Context {
	context := Context{StartTime: startTime, HowMany: howMany, Range: make([]int, howMany), Frequency: frequency, Locales: locales, Seed: seed}
	return &context
}
