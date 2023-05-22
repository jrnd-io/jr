//go:build exclude

package redis

import (
	"context"
	"jr/storage"
	"testing"
)

func TestProducer_Initialize(t *testing.T) {
	configFile := "config.json.example"

	producer, err := storage.ProducerFactory("mongo")
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

	producer, err := storage.ProducerFactory("mongo")
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

	producer, err := storage.ProducerFactory("mongo")
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
