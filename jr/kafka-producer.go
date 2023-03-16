package jr

import (
	"bufio"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"os"
	"strings"
)

func Initialize(configFile string) (*kafka.Producer, error) {
	conf := ReadConfig(configFile)
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	return p, err
}

func Close(p *kafka.Producer) {
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

func Produce(p *kafka.Producer, key []byte, data []byte, topic string) {

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or
				// permanent failure after retries have been exhausted.
				// Application level retries won't help since the client
				// is already configured to do that.
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					//fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

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
		fmt.Printf("Failed to produce message: %v\n", err)
	}

}
