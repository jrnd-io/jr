package producers

import "io"

type Producer interface {
	Produce(k []byte, v []byte, o any)
	io.Closer
}

type FactoryFunc func([]byte) Producer

var Factories = make(map[string]FactoryFunc)
