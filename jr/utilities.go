package jr

import (
	"github.com/google/uuid"
)

func counter(c string, start, step int) int {
	val, exists := JrContext.Counters[c]
	if exists {
		JrContext.Counters[c] = val + step
		return JrContext.Counters[c]
	} else {
		JrContext.Counters[c] = start
		return start
	}
}

func uniqueId() string {
	return uuid.New().String()
}

func randomBool() string {
	b := Random.Intn(2)
	if b == 0 {
		return "false"
	} else {
		return "true"
	}
}

func yesOrNo() string {
	b := Random.Intn(2)
	if b == 0 {
		return "no"
	} else {
		return "yes"
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
