package functions

import (
	"github.com/google/uuid"
)

func uniqueId() string {
	return uuid.New().String()
}

func randomBool() string {
	b := random.Intn(2)
	if b == 0 {
		return "false"
	} else {
		return "true"
	}
}

func yesOrNo() string {
	b := random.Intn(2)
	if b == 0 {
		return "no"
	} else {
		return "yes"
	}
}
