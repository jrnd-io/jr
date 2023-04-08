package redis

import (
	"fmt"
	"log"
)

type Producer interface {
	Initialize(configFile string) error
	Close()
	Produce(params ...any)
}

func ProducerFactory(producerType string) (Producer, error) {
	switch producerType {
	//	case "kafka":
	//		return &kafka.Producer{}, nil
	case "redis":
		return &RedisProducer{}, nil
	default:
		return nil, fmt.Errorf("Invalid producer type: %s", producerType)
	}
}

func main() {
	redisProducer, err := ProducerFactory("Redis")
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	err = redisProducer.Initialize("config.json")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	redisProducer.Produce([]byte("luigi"), []byte("fugaro"), 0)
	redisProducer.Close()
}
