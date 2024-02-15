package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
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
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to parse configuration parameters: %s", err)
	}

	ctx := context.Background()
	// Use Google Application Default Credentials to authorize and authenticate the client.
	// More information about Application Default Credentials and how to enable is at
	// https://developers.google.com/identity/protocols/application-default-credentials.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
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
		key = id.String() + "/.json"
	} else {
		key = string(k) + "/.json"
	}

	objectHandle := p.client.Bucket(bucket).Object(key)
	writer := objectHandle.NewWriter(ctx)
	kvPair := fmt.Sprintf("%s=%s\n", key, v)

	_, err := writer.Write([]byte(kvPair))
	if err != nil {
		log.Fatalf("Failed to write to GCS: %v", err)
	}

	writer.Close()

}

func (p *GCSProducer) Close() error {
	p.client.Close()
	return nil
}
