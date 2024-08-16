package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisProducer struct {
	client redis.Client
	Ttl    time.Duration
}

func (p *RedisProducer) Initialize(configFile string) {
	var options redis.Options

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load Redis configFile")
	}

	err = json.Unmarshal(data, &options)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	p.client = *redis.NewClient(&options)
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
