package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Bucket string `json:"bucket_name"`
}

type GCSProducer struct {
	client storage.Client
	bucket string
}

func (p *GCSProducer) Initialize(configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	ctx := context.Background()
	// Use Google Application Default Credentials to authorize and authenticate the client.
	// More information about Application Default Credentials and how to enable is at
	// https://developers.google.com/identity/protocols/application-default-credentials.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client")
	}

	p.client = *client
	p.bucket = config.Bucket
}

func (p *GCSProducer) Produce(k []byte, v []byte, o any) {
	ctx := context.Background()

	bucket := p.bucket
	var key string

	if k == nil || len(k) == 0 {
		// generate a UUID as index
		id := uuid.New()
		key = id.String()
	} else {
		key = string(k)
	}

	objectHandle := p.client.Bucket(bucket).Object(key)
	writer := objectHandle.NewWriter(ctx)
	kvPair := fmt.Sprintf("%s=%s\n", key, v)

	_, err := writer.Write([]byte(kvPair))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write to GCS")
	}

	writer.Close()

}

func (p *GCSProducer) Close() error {
	p.client.Close()
	return nil
}
