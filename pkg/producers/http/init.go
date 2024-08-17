package http

import (
	"github.com/ugol/jr/pkg/producers"
)

func init() {
	producers.Factories[Name] = Create
}
