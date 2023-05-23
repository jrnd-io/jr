package elastic

import (
    "context"
    "encoding/json"
    "log"
    "strings"
    "io/ioutil"

    "github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
)

type Config struct {
    ElasticURI  string `json:"es_uri"`
    ElasticIndex  string `json:"index"`
}

type ElasticProducer struct {
    client elasticsearch.Client
    index string
}

func (p *ElasticProducer) Initialize(configFile string) {
    var config Config
    file, err := ioutil.ReadFile(configFile)
    err = json.Unmarshal(file, &config)
    if err != nil {
    	log.Fatalf("Failed to parse configuration parameters: %s", err)
    }

    cfg := elasticsearch.Config{
    	Addresses: []string{config.ElasticURI},
    }

    client, err := elasticsearch.NewClient(cfg)

    if err != nil {
        log.Fatalf("Can't connect to Elastic: %s", err)
    }

    p.index = config.ElasticIndex
    p.client = *client
}

func (p *ElasticProducer) Produce(k []byte, v []byte, o interface{})  {

    req := esapi.IndexRequest{
    	Index:      p.index,
    	DocumentID: string(k),
    	Body:       strings.NewReader(string(v)),
    	Refresh:    "true",
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

func (p *ElasticProducer) Close() {

}
