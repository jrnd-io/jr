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
	client elasticsearch.Client
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
	p.client = *client
}

func (p *ElasticProducer) Produce(k []byte, v []byte, o any) {

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

	res, err := req.Do(context.Background(), &p.client)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in Elastic")
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatal().Str("response", res.String()).Msg("failed to index document")
	}
}

func (p *ElasticProducer) Close() error {
	log.Warn().Msg("elasticsearch Client doesn't provide a close method!")
	return nil
}
