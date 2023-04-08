package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type Producer struct {
	client redis.Client
}

func (p Producer) Initialize(configFile string) error {
	var options redis.Options

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load Redis configFile: %s", err)
	}

	err = json.Unmarshal(data, &options)
	if err != nil {
		log.Fatalf("Failed to parsa configuration parameters: %s", err)
	}

	p.client = *redis.NewClient(&options)
	return err
}

func (p Producer) Close() {
	err := p.client.Close()
	if err != nil {
		log.Fatalf("Failed to close Redis connection:\n%s", err)
	}
}

func (p Producer) Produce(params ...interface{}) {
	ctx := context.Background()
	key, ok := params[0].(string)
	if !ok {
		fmt.Println("Failed to cast interface{} to string")
	}
	val, ok := params[1].([]byte)
	if !ok {
		fmt.Println("Failed to cast interface{} to []byte]")
	}
	exp, ok := params[2].(time.Duration)
	if !ok {
		fmt.Println("Failed to cast interface{} to time.Duration")
	}
	err := p.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		log.Fatalf("Failed to write data in Redis:\n%s", err)
	}
}
