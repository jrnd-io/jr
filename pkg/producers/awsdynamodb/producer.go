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

package awsdynamodb

import (
	"context"
	"encoding/json"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	configuration Config

	client *dynamodb.Client
}

func (p *Producer) Initialize(ctx context.Context, configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ReadFile")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	if config.Table == "" {
		log.Fatal().Msg("Table is mandatory")
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load default AWS config")
	}
	client := dynamodb.NewFromConfig(awsConfig)

	p.client = client
	p.configuration = config
}

func (p *Producer) Produce(ctx context.Context, _ []byte, val []byte, _ any) {

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(val, &jsonMap); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal json")
	}

	item, err := attributevalue.MarshalMap(jsonMap)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal map")
	}

	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(p.configuration.Table),
		Item:      item,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to put item")
	}

}

func (p *Producer) Close(_ context.Context) error {
	return nil
}
