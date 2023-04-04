package jr

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"jr/jr/types"
	"log"
	"os"
	"strings"
	"time"
)

var adminClient *kafka.AdminClient
var schemaClient schemaregistry.Client
var schemaRegistry bool

func Initialize(configFile string) (*kafka.Producer, error) {

	var err error
	conf := convertInKafkaConfig(ReadConfig(configFile))

	adminClient, err = kafka.NewAdminClient(&conf)
	if err != nil {
		log.Fatalf("Failed to create admin client: %s", err)
	}
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	return p, err
}

func convertInKafkaConfig(m map[string]string) kafka.ConfigMap {
	var conf kafka.ConfigMap
	conf = make(map[string]kafka.ConfigValue)
	for k, v := range m {
		conf[k] = v
	}
	return conf
}

func InitializeSchemaRegistry(configFile string) error {
	var err error
	conf := ReadConfig(configFile)

	schemaClient, err = schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		conf["schemaRegistryURL"],
		conf["schemaRegistryUser"],
		conf["schemaRegistryPassword"]))

	if err != nil {
		log.Fatalf("Failed to create schema registry client: %s", err)
	}

	schemaRegistry = true
	return err
}

func Close(p *kafka.Producer) {
	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}

func ReadConfig(configFile string) map[string]string {

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

func Produce(p *kafka.Producer, key []byte, data []byte, topic string, serializer string, templateType string) {

	go func() {
		for e := range p.Events() {
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
				log.Printf("%9.f bytes produced to Kafka\n", stats["txmsg_bytes"])
			default:
				log.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	var ser serde.Serializer

	if schemaRegistry {
		var err error
		if serializer == "avro" {
			//ser, err = avro.NewSpecificSerializer(schemaClient, serde.ValueSerde, avro.NewSerializerConfig())
			log.Fatal("Avro not yet implemented")
		} else if serializer == "avro-generic" {
			ser, err = avro.NewGenericSerializer(schemaClient, serde.ValueSerde, avro.NewSerializerConfig())
		} else if serializer == "protobuf" {
			//ser, err = protobuf.NewSerializer(schemaClient, serde.ValueSerde, protobuf.NewSerializerConfig())
			log.Fatal("Protobuf not yet implemented")
		} else if serializer == "json-schema" {
			ser, err = jsonschema.NewSerializer(schemaClient, serde.ValueSerde, jsonschema.NewSerializerConfig())
		}
		if err != nil {
			log.Fatalf("Error creating serializer: %s\n", err)
		} else {

			t := getType(templateType)
			err := json.Unmarshal(data, t)

			if err != nil {
				log.Fatalf("Failed to unmarshal data: %s\n", err)
			}

			payload, err := ser.Serialize(topic, t)
			if err != nil {
				log.Fatalf("Failed to serialize payload: %s\n", err)
			} else {
				data = payload
			}
		}
	}

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
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

func getType(templateType string) interface{} {

	var netDevice types.NetDevice
	var user types.User

	switch templateType {
	case "net-device":
		return &netDevice
	case "user":
		return &user
	}
	return nil
}

func CreateTopic(topic string) {
	CreateTopicFull(topic, 6, 3)
}

func CreateTopicFull(topic string, partitions int, rf int) {

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDuration, _ := time.ParseDuration("60s")

	results, err := adminClient.CreateTopics(ctx,
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

	adminClient.Close()

}
