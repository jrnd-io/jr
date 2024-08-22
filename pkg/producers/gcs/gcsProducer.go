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

package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Bucket string `json:"bucket_name"`
}

type GCSProducer struct {
	client storage.Client
	bucket string
}

func (p *GCSProducer) Initialize(ctx context.Context, configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

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

func (p *GCSProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	bucket := p.bucket
	var key string

	if len(k) == 0 || strings.ToLower(string(k)) == "null" {
		// generate a UUID as index
		key = uuid.New().String()
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

func (p *GCSProducer) Close(_ context.Context) error {
	p.client.Close()
	return nil
}
