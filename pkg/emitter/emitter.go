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

package emitter

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/ctx"
	"github.com/ugol/jr/pkg/functions"
	"github.com/ugol/jr/pkg/producers/console"
	"github.com/ugol/jr/pkg/producers/elastic"
	"github.com/ugol/jr/pkg/producers/gcs"
	"github.com/ugol/jr/pkg/producers/http"
	"github.com/ugol/jr/pkg/producers/kafka"
	"github.com/ugol/jr/pkg/producers/mongoDB"
	"github.com/ugol/jr/pkg/producers/redis"
	"github.com/ugol/jr/pkg/producers/s3"
	"github.com/ugol/jr/pkg/producers/server"
	"github.com/ugol/jr/pkg/tpl"
)

type Emitter struct {
	Name             string        `mapstructure:"name"`
	Locale           string        `mapstructure:"locale"`
	Num              int           `mapstructure:"num"`
	Frequency        time.Duration `mapstructure:"frequency"`
	Duration         time.Duration `mapstructure:"duration"`
	Preload          int           `mapstructure:"preload"`
	ValueTemplate    string        `mapstructure:"valueTemplate"`
	EmbeddedTemplate string        `mapstructure:"embeddedTemplate"`
	KeyTemplate      string        `mapstructure:"keyTemplate"`
	OutputTemplate   string        `mapstructure:"outputTemplate"`
	Output           string        `mapstructure:"output"`
	Topic            string        `mapstructure:"topic"`
	Kcat             bool          `mapstructure:"kcat"`
	Oneline          bool          `mapstructure:"oneline"`
	Csv              string        `mapstructure:"csv"`
	Producer         Producer
	KTpl             tpl.Tpl
	VTpl             tpl.Tpl
}

func (e *Emitter) Initialize(conf configuration.GlobalConfiguration) {

	functions.InitCSV(e.Csv)

	templateName := e.ValueTemplate
	if e.EmbeddedTemplate == "" {
		path := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_SYSTEM_DIR, "templates"))
		templateFullPath := fmt.Sprintf("%s/%s.tpl", path, templateName)
		vt, err := os.ReadFile(templateFullPath)
		e.EmbeddedTemplate = string(vt)
		if err != nil {
			log.Printf("Template '%s' not found in %s\n", templateName, path)
		}
	}

	keyTpl, err := tpl.NewTpl("key", e.KeyTemplate, functions.FunctionsMap(), &ctx.JrContext)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create key template")
	}
	valueTpl, err := tpl.NewTpl("value", e.EmbeddedTemplate, functions.FunctionsMap(), &ctx.JrContext)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create value template")
	}

	e.KTpl = keyTpl
	e.VTpl = valueTpl

	o, _ := tpl.NewTpl("out", e.OutputTemplate, functions.FunctionsMap(), nil)
	if e.Output == "stdout" {
		e.Producer = &console.ConsoleProducer{OutputTpl: &o}
		return
	}

	if e.Output == "kafka" {
		e.Producer = createKafkaProducer(conf, e.Topic, templateName)
		return
	} else {
		if conf.SchemaRegistry {
			log.Warn().Msg("Ignoring schemaRegistry and/or serializer when output not set to kafka")
		}
	}

	if e.Output == "redis" {
		e.Producer = createRedisProducer(conf.RedisTtl, conf.RedisConfig)
		return
	}

	if e.Output == "mongo" || e.Output == "mongodb" {
		e.Producer = createMongoProducer(conf.MongoConfig)
		return
	}

	if e.Output == "elastic" {
		e.Producer = createElasticProducer(conf.ElasticConfig)
		return
	}

	if e.Output == "s3" {
		e.Producer = createS3Producer(conf.S3Config)
		return
	}

	if e.Output == "gcs" {
		e.Producer = createGCSProducer(conf.GCSConfig)
		return
	}

	if e.Output == "json" {
		e.Producer = &server.JsonProducer{OutputTpl: &o}
		return
	}

	if e.Output == "http" {
		e.Producer = createHTTPProducer(conf.HTTPConfig)
		return
	}

}

func (e *Emitter) Run(num int, o any) {

	for i := 0; i < num; i++ {

		k := e.KTpl.Execute()
		v := e.VTpl.Execute()
		kInValue := functions.GetV("KEY")

		if kInValue != "" {
			e.Producer.Produce([]byte(kInValue), []byte(v), o)
		} else {
			e.Producer.Produce([]byte(k), []byte(v), o)
		}
		ctx.JrContext.GeneratedObjects++
		ctx.JrContext.GeneratedBytes += int64(len(v))

	}

}

func createRedisProducer(ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.RedisProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createMongoProducer(mongoConfig string) Producer {
	mProducer := &mongoDB.MongoProducer{}
	mProducer.Initialize(mongoConfig)

	return mProducer
}

func createElasticProducer(elasticConfig string) Producer {
	eProducer := &elastic.ElasticProducer{}
	eProducer.Initialize(elasticConfig)

	return eProducer
}

func createS3Producer(s3Config string) Producer {
	sProducer := &s3.S3Producer{}
	sProducer.Initialize(s3Config)

	return sProducer
}

func createGCSProducer(gcsConfig string) Producer {
	gProducer := &gcs.GCSProducer{}
	gProducer.Initialize(gcsConfig)

	return gProducer
}

func createHTTPProducer(httpConfig string) Producer {
	httpProducer := &http.Producer{}
	httpProducer.Initialize(httpConfig)

	return httpProducer
}

func createKafkaProducer(conf configuration.GlobalConfiguration, topic string, templateType string) *kafka.KafkaManager {

	kManager := &kafka.KafkaManager{
		Serializer:   conf.Serializer,
		Topic:        topic,
		TemplateType: templateType,
	}

	kManager.Initialize(conf.KafkaConfig)

	if conf.SchemaRegistry {
		kManager.InitializeSchemaRegistry(conf.RegistryConfig)
	}
	if conf.AutoCreate {
		kManager.CreateTopic(topic)
	}
	return kManager
}
