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

func (p *RedisProducer) Close(_ context.Context) error {
	err := p.client.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close Redis connection")
	}
	return err
}

func (p *RedisProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	err := p.client.Set(ctx, string(k), string(v), p.Ttl).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in Redis")
	}
}
