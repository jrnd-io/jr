package mongoDB

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

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

func (p *MongoProducer) Initialize(configFile string) {
	var config Config
	file, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to parse configuration parameters: %s", err)
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

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatalf("Can't connect to Mongo: %s", err)
	}

	p.client = *client
}

func (p *MongoProducer) Produce(k []byte, v []byte, o interface{}) {

	collection := p.client.Database(p.database).Collection(p.collection)

	var dev map[string]interface{}
	err := json.Unmarshal(v, &dev)
	if err != nil {
		log.Fatalf("Failed to unmarshal json:\n%s", err)
	}

	if k == nil || len(k) == 0 {
		dev["_id"] = string(k)
	}

	_, err = collection.InsertOne(context.Background(), dev)
	if err != nil {
		log.Fatalf("Failed to write data in Mongo:\n%s", err)
	}
}

func (p *MongoProducer) Close() error {
	err := p.client.Disconnect(context.Background())
	if err != nil {
		log.Println("Failed to close Mongo connection:\n%s", err)
	}
	return err
}
