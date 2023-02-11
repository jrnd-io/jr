package functions

import (
	"github.com/google/uuid"
	"math/rand"
)

func uniqueId() string {
	return uuid.New().String()
}

func randomBool() string {
	b := rand.Intn(2)
	if b == 0 {
		return "false"
	} else {
		return "true"
	}
}

func yesOrNo() string {
	b := rand.Intn(2)
	if b == 0 {
		return "no"
	} else {
		return "yes"
	}
}
