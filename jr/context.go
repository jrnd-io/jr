package jr

import "time"

type Context struct {
	StartTime        time.Time
	GeneratedObjects int64
	GeneratedBytes   int64
	HowMany          int
	Range            []int
	Frequency        int
	Localization     []string
	Seed             int64
}

func NewContext(startTime time.Time, howMany int, frequency int, locales []string, seed int64) *Context {
	context := Context{StartTime: startTime, HowMany: howMany, Range: make([]int, howMany), Frequency: frequency, Localization: locales, Seed: seed}
	return &context
}
