package elastic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
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
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to ReadFile: %s", err)
	}
	if err != nil {
		log.Fatalf("Failed to parse configuration parameters: %s", err)
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
		log.Fatalf("Can't connect to Elastic: %s", err)
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
		log.Fatalf("Failed to write data in Elastic:\n%s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("failed to index document: %s", res.String())
	}
}

func (p *ElasticProducer) Close() error {
	log.Println("elasticsearch Client doesn't provide a close method!")
	return nil
}
