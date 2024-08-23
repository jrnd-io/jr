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

package kafka

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/rules/encryption"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/rules/encryption/awskms"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/rules/encryption/azurekms"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/rules/encryption/gcpkms"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/rules/encryption/hcvault"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/jrnd-io/jr/pkg/types"

	"github.com/rs/zerolog/log"
)

type KafkaManager struct {
	producer       *kafka.Producer
	admin          *kafka.AdminClient
	schema         schemaregistry.Client
	schemaRegistry bool
	Topic          string
	Serializer     string
	TemplateType   string
	fleEnabled     bool
}

func (k *KafkaManager) Initialize(configFile string) {

	var err error
	conf := convertInKafkaConfig(readConfig(configFile))
	k.admin, err = kafka.NewAdminClient(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create admin client")
	}
	k.producer, err = kafka.NewProducer(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create producer")
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
		log.Fatal().Err(err).Msg("Failed to create schema registry client")
	}

	if k.Serializer == "avro" || k.Serializer == "avro-generic" {
		verifyCSFLE(conf, k)
	}

	k.schemaRegistry = true
}

func verifyCSFLE(conf map[string]string, k *KafkaManager) {
	if conf["kekName"] != "" && conf["kmsType"] != "" && conf["kmsKeyID"] != "" {
		registerProviders()

		// load avro schema file: CSFLE requires schema registration
		_, currentFilePath, _, _ := runtime.Caller(0)
		currentDir := filepath.Dir(currentFilePath)
		filePath := filepath.Join(currentDir, "../../types/"+k.TemplateType+".avsc")
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open file")
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read file")
		}
		contentString := string(content)

		// check presence of PII in schema
		substring := `"confluent:tags": [ "PII" ]`
		normalizedJSON := normalizeWhitespace(contentString)
		normalizedSubstring := normalizeWhitespace(substring)

		if strings.Contains(normalizedJSON, normalizedSubstring) {
			// upper-casing the first letter of the fields --> name - required by https://pkg.go.dev/github.com/actgardner/gogen-avro#readme-naming
			re := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)
			result := re.ReplaceAllStringFunc(contentString, func(match string) string {
				// extract the name part after "name: "
				parts := re.FindStringSubmatch(match)
				fmt.Print(len(parts))
				if len(parts) > 1 {
					name := parts[1]
					// capitalize the first letter of the name
					capitalized := capitalizeFirstLetter(name)
					// replace the original match with the new capitalized version
					return "\"name\":" + "\"" + capitalized + "\""
				}
				return match
			})

			// register the avro schema adding rule set: PII
			schema := schemaregistry.SchemaInfo{
				Schema:     result,
				SchemaType: "AVRO",
				RuleSet: &schemaregistry.RuleSet{
					DomainRules: []schemaregistry.Rule{
						{
							Name: "encryptPII",
							Kind: "TRANSFORM",
							Mode: "WRITEREAD",
							Type: "ENCRYPT",
							Tags: []string{"PII"},
							Params: map[string]string{
								"encrypt.kek.name":   conf["kekName"],
								"encrypt.kms.type":   conf["kmsType"],
								"encrypt.kms.key.id": conf["kmsKeyID"],
							},
							OnFailure: "ERROR,NONE",
						},
					},
				},
			}
			_, err = k.schema.Register(k.Topic+"-value", schema, true)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to register schema")
			}

			k.fleEnabled = true

		}

	}
}

func registerProviders() {
	awskms.Register()
	azurekms.Register()
	gcpkms.Register()
	hcvault.Register()
	encryption.Register()
}

func (k *KafkaManager) Close(_ context.Context) error {
	k.admin.Close()
	k.producer.Flush(15 * 1000)
	k.producer.Close()
	return nil
}

func (k *KafkaManager) Produce(_ context.Context, key []byte, data []byte, _ any) {

	go listenToEventsFrom(k.producer, k.Topic)

	var ser serde.Serializer

	if k.schemaRegistry {
		var err error

		if k.Serializer == "avro" || k.Serializer == "avro-generic" {
			serConfig := avrov2.NewSerializerConfig()
			// CSFLE requires auto register to false
			if k.fleEnabled {
				serConfig.AutoRegisterSchemas = false
				serConfig.UseLatestVersion = true
			}
			ser, err = avrov2.NewSerializer(k.schema, serde.ValueSerde, serConfig)
		} else if k.Serializer == "protobuf" {
			// ser, err = protobuf.NewSerializer(k.schema, serde.ValueSerde, protobuf.NewSerializerConfig())
			log.Fatal().Msg("Protobuf not yet implemented")
		} else if k.Serializer == "json-schema" {
			ser, err = jsonschema.NewSerializer(k.schema, serde.ValueSerde, jsonschema.NewSerializerConfig())
		} else {
			log.Fatal().Str("serializer", k.Serializer).Msg("Serializer not supported")
		}
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating serializer")
		} else {

			t := types.GetType(k.TemplateType)
			err := json.Unmarshal(data, &t)

			if err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal data")
			}

			payload, err := ser.Serialize(k.Topic, t)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to serialize payload")
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
		// Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)

	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			// time.Sleep(time.Second)
			// continue
		}
		log.Error().Err(err).Msg("Failed to produce message")
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
		log.Error().Err(err).Msg("Problem during the topic creation")

	}

	// Check for specific topic errors
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Fatal().
				Str("topic", result.Topic).
				Str("code", result.Error.Code().String()).
				Err(result.Error).
				Msg("Topic creation failed")
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
				log.Error().Err(m.TopicPartition.Error).Msg("Delivery failed")
			} else {
				// fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		case kafka.Error:
			log.Error().Err(ev).Msg("Kafka error")
		case *kafka.Stats:
			// https://github.com/confluentinc/librdkafka/blob/master/STATISTICS.md
			var stats map[string]interface{}
			err := json.Unmarshal([]byte(e.String()), &stats)
			if err != nil {
				return
			}
			txbytes := fmt.Sprintf("%9.f", stats["txmsg_bytes"])
			b, _ := strconv.Atoi(strings.TrimSpace(txbytes))

			if b > 0 {
				log.Info().
					Str("bytes", txbytes).
					Str("topic", topic).
					Msg("Bytes produced to topic")
			}
		default:
			log.Warn().Interface("ev", ev).Msg("Ignored event")
		}
	}

}

func readConfig(configFile string) map[string]string {

	m := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open configuration file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Error in closing file")
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
		log.Fatal().Err(err).Msg("Failed to read file")
	}
	return m
}

func convertInKafkaConfig(m map[string]string) kafka.ConfigMap {
	conf := make(map[string]kafka.ConfigValue)
	for k, v := range m {
		conf[k] = v
	}
	return conf
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func normalizeWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}
