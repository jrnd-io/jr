// Copyright © 2022 Vincenzo Marchese <vincenzo.marchese@gmail.com>
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
package azcosmosdb

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	configuration Config
	client        *azcosmos.Client
}

func (p *Producer) Initialize(configFile string) {
	cfgBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	config := Config{}
	if err := json.Unmarshal(cfgBytes, &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	if config.Endpoint == "" {
		log.Fatal().Msg("Endpoint is mandatory")
	}

	if config.PrimaryAccountKey == "" {
		log.Fatal().Msg("PrimaryAccountKey is mandatory")
	}

	if config.PartitionKey == "" {
		log.Fatal().Msg("PartitionKey is mandatory")
	}

	cred, err := azcosmos.NewKeyCredential(config.PrimaryAccountKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create key credential")
	}

	client, err := azcosmos.NewClientWithKey(config.Endpoint, cred, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create  azure cosmosdb client")
	}

	p.configuration = config
	p.client = client

}

func (p *Producer) Produce(_ []byte, v []byte, _ any) {

	// This is ugly but it works
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(v, &jsonMap); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal json")
	}

	// getting partition key value
	pkValue := jsonMap[p.configuration.PartitionKey]
	if pkValue == nil {
		log.Fatal().Str("partition_key", p.configuration.PartitionKey).Msg("Partition key not found in value")
	}
	log.Debug().Str("pkValue", pkValue.(string)).Msg("Partition key value")

	container, err := p.client.NewContainer(p.configuration.Database, p.configuration.Container)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create container")
	}

	pk := azcosmos.NewPartitionKeyString(pkValue.(string))
	resp, err := container.CreateItem(context.Background(), pk, v, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create item")
	}

	log.Debug().Interface("resp", resp).Msg("Item created")

}

func (p *Producer) Close() error {
	return nil
}
