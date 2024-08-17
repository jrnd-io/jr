//Copyright © 2022 Vincenzo Marchese <vincenzo.marchese@gmail.com>
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

package emitter

import (
	"github.com/ugol/jr/pkg/producers/elastic"
	"github.com/ugol/jr/pkg/producers/gcs"
	"github.com/ugol/jr/pkg/producers/http"
	"github.com/ugol/jr/pkg/producers/mongoDB"
	"github.com/ugol/jr/pkg/producers/redis"
	"github.com/ugol/jr/pkg/producers/s3"
)

const (
	ConsoleProducer = "console"
	MongoProducer   = "mongo"
	MongoDBProducer = "mongodb"
	RedisProducer   = "redis"
	ElasticProducer = "elastic"
	S3Producer      = "s3"
	GCSProducer     = "gcs"
	HTTPProducer    = "http"
)

type ProducerFactoryFunc func([]byte) Producer

var ProducerFactories = map[string]ProducerFactoryFunc{
	MongoProducer:   createMongoProducer,
	MongoDBProducer: createMongoProducer,
	RedisProducer:   createRedisProducer,
	ElasticProducer: createElasticProducer,
	S3Producer:      createS3Producer,
	GCSProducer:     createGCSProducer,
	HTTPProducer:    createHTTPProducer,
}

func createRedisProducer(config []byte) Producer {
	producer := &redis.RedisProducer{}
	producer.Initialize(config)
	return producer
}

func createMongoProducer(config []byte) Producer {
	producer := &mongoDB.MongoProducer{}
	producer.Initialize(config)

	return producer
}

func createElasticProducer(config []byte) Producer {
	producer := &elastic.ElasticProducer{}
	producer.Initialize(config)

	return producer
}

func createS3Producer(config []byte) Producer {
	producer := &s3.S3Producer{}
	producer.Initialize(config)

	return producer
}

func createGCSProducer(config []byte) Producer {
	producer := &gcs.GCSProducer{}
	producer.Initialize(config)

	return producer
}

func createHTTPProducer(config []byte) Producer {
	producer := &http.Producer{}
	producer.Initialize(config)

	return producer
}
