package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Options redis.Options `json:"options"`
	TTL     string        `json:"ttl"`
}
type RedisProducer struct {
	client redis.Client
	Ttl    time.Duration
}

func (p *RedisProducer) Initialize(configBytes []byte) {
	//var options redis.Options

	config := &Config{}

	err := json.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	p.client = *redis.NewClient(&config.Options)

	if config.TTL == "" {
		p.Ttl = -1 * time.Nanosecond
	}

	ttl, err := time.ParseDuration(config.TTL)
	if err != nil {
		p.Ttl = -1 * time.Nanosecond
	}

	p.Ttl = ttl
}

func (p *RedisProducer) Close() error {
	err := p.client.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close Redis connection")
	}
	return err
}

func (p *RedisProducer) Produce(k []byte, v []byte, o any) {
	ctx := context.Background()
	err := p.client.Set(ctx, string(k), string(v), p.Ttl).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in Redis")
	}
}
