package jr

import (
	"bufio"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"strings"
)

var conf kafka.ConfigMap
var p *kafka.Producer
var err error

func Initialize(configFile string) {
	conf = ReadConfig(configFile)
	p, err = kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
}

func Produce(key []byte, data []byte, topic string) {

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          data,
	}, nil)

}

func Close() {
	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}

func ReadConfig(configFile string) kafka.ConfigMap {

	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

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
