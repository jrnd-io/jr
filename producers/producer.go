package producers

import (
	"fmt"
	"github.com/ugol/jr/producers/redis"
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
	//		return kafka.Producer{}, nil
	case "redis":
		return redis.Producer{}, nil
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
