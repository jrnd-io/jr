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

package mongoDB

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI   string `json:"mongo_uri"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type MongoProducer struct {
	client     mongo.Client
	database   string
	collection string
}

func (p *MongoProducer) Initialize(ctx context.Context, configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	if config.Username != "" && config.Password != "" {
		clientOptions.Auth = &options.Credential{
			Username: config.Username,
			Password: config.Password,
		}
	}

	p.collection = config.Collection
	p.database = config.Database

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to Mongo")
	}

	p.client = *client
}

func (p *MongoProducer) Produce(ctx context.Context, k []byte, v []byte, _ any) {

	collection := p.client.Database(p.database).Collection(p.collection)

	var dev map[string]interface{}
	err := json.Unmarshal(v, &dev)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal json")
	}

	if len(k) == 0 {
		dev["_id"] = string(k)
	}

	_, err = collection.InsertOne(ctx, dev)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in Mongo")
	}
}

func (p *MongoProducer) Close(ctx context.Context) error {
	err := p.client.Disconnect(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close Mongo connection")
	}
	return err
}
