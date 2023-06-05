//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package kafka

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/ugol/jr/pkg/types"
)

type KafkaManager struct {
	producer       *kafka.Producer
	admin          *kafka.AdminClient
	schema         schemaregistry.Client
	schemaRegistry bool
	Topic          string
	Serializer     string
	TemplateType   string
}

func (k *KafkaManager) Initialize(configFile string) {

	var err error
	conf := convertInKafkaConfig(readConfig(configFile))
	k.admin, err = kafka.NewAdminClient(&conf)
	if err != nil {
		log.Fatalf("Failed to create admin client: %s", err)
	}
	k.producer, err = kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
}

func (k *KafkaManager) InitializeSchemaRegistry(configFile string) {
	var err error
	conf := readConfig(configFile)

	k.schema, err = schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		conf["schemaRegistryURL"],
		conf["schemaRegistryUser"],
		conf["schemaRegistryPassword"]))

	if err != nil {
		log.Fatalf("Failed to create schema registry client: %s", err)
	}

	k.schemaRegistry = true
}

func (k *KafkaManager) Close() {
	k.admin.Close()
	k.producer.Flush(15 * 1000)
	k.producer.Close()
}

func (k *KafkaManager) Produce(key []byte, data []byte, o interface{}) {

	go listenToEventsFrom(k.producer, k.Topic)

	var ser serde.Serializer

	if k.schemaRegistry {
		var err error
		if k.Serializer == "avro" {
			ser, err = avro.NewSpecificSerializer(k.schema, serde.ValueSerde, avro.NewSerializerConfig())
		} else if k.Serializer == "avro-generic" {
			ser, err = avro.NewGenericSerializer(k.schema, serde.ValueSerde, avro.NewSerializerConfig())
		} else if k.Serializer == "protobuf" {
			//ser, err = protobuf.NewSerializer(k.schema, serde.ValueSerde, protobuf.NewSerializerConfig())
			log.Fatal("Protobuf not yet implemented")
		} else if k.Serializer == "json-schema" {
			ser, err = jsonschema.NewSerializer(k.schema, serde.ValueSerde, jsonschema.NewSerializerConfig())
		} else {
			log.Fatalf("Serializer '%v' not supported", k.Serializer)
		}
		if err != nil {
			log.Fatalf("Error creating serializer: %s\n", err)
		} else {

			t := types.GetType(k.TemplateType)
			err := json.Unmarshal(data, &t)

			if err != nil {
				log.Fatalf("Failed to unmarshal data: %s\n", err)
			}

			payload, err := ser.Serialize(k.Topic, t)
			if err != nil {
				log.Fatalf("Failed to serialize payload: %s\n", err)
			} else {
				data = payload
			}
		}
	}

	if strings.ToLower(string(key)) == "null" {
		key = nil
	}

	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.Topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          data,
		//Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)

	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			//time.Sleep(time.Second)
			//continue
		}
		log.Printf("Failed to produce message: %v\n", err)
	}

}

func (k *KafkaManager) CreateTopic(topic string) {
	k.CreateTopicFull(topic, 6, 3)
}

func (k *KafkaManager) CreateTopicFull(topic string, partitions int, rf int) {

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDuration, _ := time.ParseDuration("60s")

	results, err := k.admin.CreateTopics(ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: rf}},
		kafka.SetAdminOperationTimeout(maxDuration))

	if err != nil {
		log.Printf("Problem during the topic creation: %v\n", err)

	}

	// Check for specific topic errors
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Fatalf("Topic creation failed for %s: %v", result.Topic, result.Error.String())
		}
	}

	k.admin.Close()

}

func listenToEventsFrom(k *kafka.Producer, topic string) {

	for e := range k.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				//fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		case kafka.Error:
			log.Printf("Error: %v\n", ev)
		case *kafka.Stats:
			// https://github.com/confluentinc/librdkafka/blob/master/STATISTICS.md
			var stats map[string]interface{}
			err := json.Unmarshal([]byte(e.String()), &stats)
			if err != nil {
				return
			}
			log.Printf("%9.f bytes produced to topic %s \n", stats["txmsg_bytes"], topic)
		default:
			log.Printf("Ignored event: %s\n", ev)
		}
	}

}

func readConfig(configFile string) map[string]string {

	m := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error in closing file: %s", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	return m
}

func convertInKafkaConfig(m map[string]string) kafka.ConfigMap {
	var conf kafka.ConfigMap
	conf = make(map[string]kafka.ConfigValue)
	for k, v := range m {
		conf[k] = v
	}
	return conf
}
