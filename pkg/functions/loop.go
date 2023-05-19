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

package functions

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ugol/jr/pkg/producers/console"
	"github.com/ugol/jr/pkg/producers/kafka"
	"github.com/ugol/jr/pkg/producers/redis"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

type Producer interface {
	Close()
	Produce(k []byte, v []byte)
}

func DoTemplates(conf Configuration) {

	valueTemplate := make([][]byte, len(conf.TemplateNames))
	preloadTemplate := make([][]byte, len(conf.Preload))

	var err error

	if conf.EmbeddedTemplate {
		valueTemplate[0] = []byte(conf.TemplateNames[0])
	} else if conf.TemplateFileName {
		for i := range conf.TemplateNames {
			valueTemplate[i], err = os.ReadFile(os.ExpandEnv(conf.TemplateNames[i]))
			JrContext.TemplateType[i] = conf.TemplateNames[i]
		}
		JrContext.NumTemplates = len(conf.TemplateNames)
	} else {
		for i := range conf.TemplateNames {
			templatePath := fmt.Sprintf("%s/%s.tpl", conf.TemplateDir, conf.TemplateNames[i])
			valueTemplate[i], err = os.ReadFile(templatePath)
			JrContext.TemplateType[i] = conf.TemplateNames[i]
		}
		JrContext.NumTemplates = len(conf.TemplateNames)
	}

	for i := range conf.Preload {
		templatePath := fmt.Sprintf("%s/%s.tpl", conf.TemplateDir, conf.Preload[i])
		preloadTemplate[i], err = os.ReadFile(templatePath)
		JrContext.PreloadTemplateType[i] = conf.Preload[i]
	}

	if err != nil {
		log.Fatal(err)
	}

	outTemplate, err := template.New("out").Parse(conf.OutputTemplate)
	if err != nil {
		log.Fatal(err)
	}

	key, err := template.New("key").Funcs(FunctionsMap()).Parse(conf.KeyTemplate)
	if err != nil {
		log.Fatal(err)
	}

	value := template.New("value").Funcs(FunctionsMap())
	for i := 0; i < len(conf.TemplateNames); i++ {
		_, err = value.New(strconv.Itoa(i)).Parse(string(valueTemplate[i]))
		if err != nil {
			log.Fatal(err)
		}
	}

	producer := make([]Producer, len(conf.TemplateNames))

	if conf.Output == "stdout" {
		for i := range conf.TemplateNames {
			producer[i] = &console.ConsoleProducer{OutTemplate: outTemplate}
		}
	}

	if conf.Output == "kafka" {
		if len(conf.TemplateNames) != len(conf.Topic) {
			log.Fatalf("There are %d templates and %d topics, every templates must have its own topic. \nFor example: jr run user net-device -o kafka -t \"test\",\"test1\"", len(conf.TemplateNames), len(conf.Topic))
		}
		for i := range conf.TemplateNames {
			producer[i] = createKafkaProducer(conf, i)
		}
	} else {
		if conf.SchemaRegistry {
			log.Println("Ignoring schemaRegistry and/or serializer when output not set to kafka")
		}
	}

	if conf.Output == "redis" {
		for i := range conf.TemplateNames {
			producer[i] = createRedisProducer(conf.RedisTtl, conf.RedisConfig)
		}
	}

	if conf.Output == "mongo" {
		log.Fatal("Not yet implemented")
	}

	orderedParsedTemplates := orderValueTemplates(value, conf.TemplateNames)

	infinite := true
	if conf.Duration > 0 {
		timer := time.NewTimer(conf.Duration)

		go func() {
			<-timer.C
			infinite = false
		}()
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if conf.Frequency != -1 {
	Infinite:
		for ok := true; ok; ok = infinite {
			select {
			case <-time.After(conf.Frequency):
				for i := range conf.TemplateNames {
					generatorLoop(key, orderedParsedTemplates[i], conf.Oneline, producer[i])
				}
			case <-ctx.Done():
				stop()
				break Infinite
			}
		}
	} else {
		for i := range conf.TemplateNames {
			generatorLoop(key, orderedParsedTemplates[i], conf.Oneline, producer[i])
		}
	}

	closeProducers(producer)

	time.Sleep(100 * time.Millisecond)
	writeStats()
}

func closeProducers(producer []Producer) {
	for _, p := range producer {
		p.Close()
	}
}

func orderValueTemplates(valueTemplate *template.Template, templateNames []string) []*template.Template {
	parsedTemplates := make([]*template.Template, len(templateNames))
	orderedParsedTemplates := make([]*template.Template, len(templateNames))
	copy(parsedTemplates, valueTemplate.Templates())

	for i := range templateNames {
		index, _ := strconv.Atoi(parsedTemplates[i].Name())
		orderedParsedTemplates[index] = parsedTemplates[i]
	}
	return orderedParsedTemplates
}

func generatorLoop(key *template.Template, value *template.Template, oneline bool, producer Producer) {
	for range JrContext.Range {
		k, v, _ := executeTemplate(key, value, oneline)
		producer.Produce([]byte(k), []byte(v))
	}
}

func createRedisProducer(ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.RedisProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createKafkaProducer(conf Configuration, index int) *kafka.KafkaManager {

	kManager := &kafka.KafkaManager{
		Serializer:   conf.Serializer,
		Topic:        conf.Topic[index],
		TemplateType: JrContext.TemplateType[index],
	}

	kManager.Initialize(conf.KafkaConfig)

	if conf.SchemaRegistry {
		kManager.InitializeSchemaRegistry(conf.RegistryConfig)
		if conf.Kcat {
			log.Println("Ignoring kcat when schemaRegistry is enabled")
		}
	}
	if conf.Autocreate {
		for i := range conf.Topic {
			kManager.CreateTopic(conf.Topic[i])
		}
	}
	return kManager
}

func writeStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", JrContext.GeneratedObjects)
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Number of templates (Objects): %d\n", JrContext.NumTemplates)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}

func executeTemplate(key *template.Template, value *template.Template, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, JrContext); err != nil {
		log.Println(err)
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, JrContext); err != nil {
		log.Println(err)
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	JrContext.GeneratedObjects++
	JrContext.GeneratedBytes += int64(len(v))

	return k, v, err
}
