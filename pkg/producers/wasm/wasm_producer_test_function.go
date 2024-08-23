//go:build tinygo.wasm

// tinygo build -o pkg/producers/wasm/wasm_producer_test_function.wasm -target=wasi pkg/producers/wasm/wasm_producer_test_function.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

const NoError = 0

//export produce
func _produce(size uint32) uint64 {
	b := make([]byte, size)

	_, err := io.ReadAtLeast(os.Stdin, b, int(size))
	if err != nil {
		return e(err)
	}

	in := make(map[string][]byte)

	err = json.Unmarshal(b, &in)
	if err != nil {
		return e(err)
	}

	out := bytes.ToUpper(in["v"])

	_, err = os.Stdout.Write(out)
	if err != nil {
		return e(err)
	}

	return NoError
}

func e(err error) uint64 {
	if err == nil {
		return NoError
	}

	_, _ = os.Stderr.WriteString(err.Error())
	return uint64(len(err.Error()))
}

func main() {}
