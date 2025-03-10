//go:build exclude

package wamp

import (
	"context"
	"testing"

	"github.com/jrnd-io/jr/pkg/producers/wamp"
)

func TestProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer, err := wamp.ProducerFactory("wamp")
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

	producer, err := wamp.ProducerFactory("wamp")
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

	producer, err := wamp.ProducerFactory("wamp")
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
