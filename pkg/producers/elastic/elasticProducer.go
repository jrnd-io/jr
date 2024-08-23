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

package elastic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ElasticURI      string `json:"es_uri"`
	ElasticIndex    string `json:"index"`
	ElasticUsername string `json:"username"`
	ElasticPassword string `json:"password"`
}

type ElasticProducer struct {
	client *elasticsearch.Client
	index  string
}

func (p *ElasticProducer) Initialize(configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ReadFile")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{config.ElasticURI},
		Username:  config.ElasticUsername,
		Password:  config.ElasticPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Elastic")
	}

	p.index = config.ElasticIndex
	p.client = client
}

func (p *ElasticProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {

	var req esapi.IndexRequest

	if k == nil || len(k) == 0 {
		// generate a UUID as index
		id := uuid.New()

		req = esapi.IndexRequest{
			Index:      p.index,
			DocumentID: id.String(),
			Body:       strings.NewReader(string(v)),
			Refresh:    "true",
		}
	} else {
		req = esapi.IndexRequest{
			Index:      p.index,
			DocumentID: string(k),
			Body:       strings.NewReader(string(v)),
			Refresh:    "true",
		}
	}

	res, err := req.Do(ctx, p.client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in Elastic")
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatal().Str("response", res.String()).Msg("failed to index document")
	}
}

func (p *ElasticProducer) Close(_ context.Context) error {
	log.Warn().Msg("elasticsearch Client doesn't provide a close method!")
	return nil
}
