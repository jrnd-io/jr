//go:build exclude

// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package redis

import (
	"context"
	"github.com/jrnd-io/jr/pkg/storage"
	"testing"
)

func TestProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer, err := storage.ProducerFactory("redis")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
}

func TestProducer_Close(t *testing.T) {
	configFile := "config.json.example"

	producer, err := storage.ProducerFactory("redis")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	producer.Close()
}

func TestProducer_Produce(t *testing.T) {
	configFile := "config.json.example"

	producer, err := storage.ProducerFactory("redis")
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}
	err = producer.Initialize(configFile)
	if err != nil {
		t.Fatalf("Error initializing producer: %v", err)
	}

	ctx := context.Background()
	key := "loo"
	val := "foo"
	exp := 0
	producer.Produce(ctx, key, val, exp)

}
