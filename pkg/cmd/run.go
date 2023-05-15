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

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	functions2 "github.com/ugol/jr/pkg/functions"
	"github.com/ugol/jr/pkg/producers/console"
	"github.com/ugol/jr/pkg/producers/kafka"
	"github.com/ugol/jr/pkg/producers/redis"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Producer interface {
	Close()
	Produce(k []byte, v []byte)
}

type Configuration struct {
	templateNames    []string
	keyTemplate      string
	outputTemplate   string
	embeddedTemplate bool
	templateFileName bool
	kcat             bool
	output           string
	oneline          bool
	locale           string
	num              int
	frequency        time.Duration
	duration         time.Duration
	seed             int64
	kafkaConfig      string
	registryConfig   string
	topic            []string
	preload          []string
	preloadSize      []int
	templateDir      string
	autocreate       bool
	schemaRegistry   bool
	serializer       string
	redisTtl         time.Duration
	redisConfig      string
}

var runCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template. 
  Without any other flag, [template] is just the name of a template in the templates directory, which by default is in '$HOME/.jr/templates' Example: 
jr run net-device
  With the --template flag, [template] is a string containing a full template. Example:
jr run --template "{{name}}"
 With the -templateFileName flag [template] is a file name with a template. Example:
jr run --templateFileName ~/.jr/templates/net_device.tpl
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		keyTemplate, _ := cmd.Flags().GetString("key")
		outputTemplate, _ := cmd.Flags().GetString("outputTemplate")
		embeddedTemplate, _ := cmd.Flags().GetBool("template")
		templateFileName, _ := cmd.Flags().GetBool("templateFileName")
		kcat, _ := cmd.Flags().GetBool("kcat")
		output, _ := cmd.Flags().GetString("output")
		oneline, _ := cmd.Flags().GetBool("oneline")
		locale, _ := cmd.Flags().GetString("locale")

		num, _ := cmd.Flags().GetInt("num")
		frequency, _ := cmd.Flags().GetDuration("frequency")
		duration, _ := cmd.Flags().GetDuration("duration")
		seed, _ := cmd.Flags().GetInt64("seed")
		kafkaConfig, _ := cmd.Flags().GetString("kafkaConfig")
		registryConfig, _ := cmd.Flags().GetString("registryConfig")
		topic, _ := cmd.Flags().GetStringSlice("topic")
		preload, _ := cmd.Flags().GetStringSlice("preload")
		preloadSize, _ := cmd.Flags().GetIntSlice("preloadSize")

		templateDir, _ := cmd.Flags().GetString("templateDir")
		templateDir = os.ExpandEnv(templateDir)

		autocreate, _ := cmd.Flags().GetBool("autocreate")
		schemaRegistry, _ := cmd.Flags().GetBool("schemaRegistry")
		serializer, _ := cmd.Flags().GetString("serializer")

		redisTtl, _ := cmd.Flags().GetDuration("redis.ttl")
		redisConfig, _ := cmd.Flags().GetString("redisConfig")

		if kcat {
			oneline = true
			output = "stdout"
			outputTemplate = "{{.K}},{{.V}}\n"
		}

		conf := Configuration{
			templateNames:    args,
			keyTemplate:      keyTemplate,
			outputTemplate:   outputTemplate,
			embeddedTemplate: embeddedTemplate,
			templateFileName: templateFileName,
			kcat:             kcat,
			output:           output,
			topic:            topic,
			oneline:          oneline,
			locale:           locale,
			num:              num,
			frequency:        frequency,
			duration:         duration,
			seed:             seed,
			kafkaConfig:      kafkaConfig,
			registryConfig:   registryConfig,
			templateDir:      templateDir,
			autocreate:       autocreate,
			schemaRegistry:   schemaRegistry,
			serializer:       serializer,
			redisTtl:         redisTtl,
			redisConfig:      redisConfig,
			preload:          preload,
			preloadSize:      preloadSize,
		}
		doTemplates(conf)
	},
}

func doTemplates(conf Configuration) {

	valueTemplate := make([][]byte, len(conf.templateNames))

	var err error

	if conf.embeddedTemplate {
		valueTemplate[0] = []byte(conf.templateNames[0])
	} else if conf.templateFileName {
		for i := range conf.templateNames {
			valueTemplate[i], err = os.ReadFile(os.ExpandEnv(conf.templateNames[i]))
			functions2.JrContext.TemplateType[i] = conf.templateNames[i]
		}
		functions2.JrContext.NumTemplates = len(conf.templateNames)
	} else {
		for i := range conf.templateNames {
			templatePath := fmt.Sprintf("%s/%s.tpl", conf.templateDir, conf.templateNames[i])
			valueTemplate[i], err = os.ReadFile(templatePath)
			functions2.JrContext.TemplateType[i] = conf.templateNames[i]
		}
		functions2.JrContext.NumTemplates = len(conf.templateNames)
	}
	if err != nil {
		log.Fatal(err)
	}

	outTemplate, err := template.New("out").Parse(conf.outputTemplate)
	if err != nil {
		log.Fatal(err)
	}

	key, err := template.New("key").Funcs(functions2.FunctionsMap()).Parse(conf.keyTemplate)
	if err != nil {
		log.Fatal(err)
	}

	value := template.New("value").Funcs(functions2.FunctionsMap())
	for i := 0; i < len(conf.templateNames); i++ {
		_, err = value.New(strconv.Itoa(i)).Parse(string(valueTemplate[i]))
		if err != nil {
			log.Fatal(err)
		}
	}

	producer := make([]Producer, len(conf.templateNames))

	if conf.output == "stdout" {
		for i := range conf.templateNames {
			producer[i] = &console.ConsoleProducer{OutTemplate: outTemplate}
		}
	}

	if conf.output == "kafka" {
		if len(conf.templateNames) != len(conf.topic) {
			log.Fatalf("There are %d templates and %d topics, every templates must have its own topic. \nFor example: jr run user net-device -o kafka -t \"test\",\"test1\"", len(conf.templateNames), len(conf.topic))
		}
		for i := range conf.templateNames {
			producer[i] = createKafkaProducer(conf.serializer, conf.topic, i, conf.kafkaConfig, conf.schemaRegistry, conf.registryConfig, conf.kcat, conf.autocreate)
		}
	} else {
		if conf.schemaRegistry {
			log.Println("Ignoring schemaRegistry and/or serializer when output not set to kafka")
		}
	}

	if conf.output == "redis" {
		for i := range conf.templateNames {
			producer[i] = createRedisProducer(conf.redisTtl, conf.redisConfig)
		}
	}

	if conf.output == "mongo" {
		log.Fatal("Not yet implemented")
	}

	configureJrContext(conf.seed, conf.num, conf.frequency, conf.locale, conf.templateDir)
	orderedParsedTemplates := orderValueTemplates(value, conf.templateNames)

	infinite := true
	if conf.duration > 0 {
		timer := time.NewTimer(conf.duration)

		go func() {
			<-timer.C
			infinite = false
		}()
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if conf.frequency != -1 {
	Infinite:
		for ok := true; ok; ok = infinite {
			select {
			case <-time.After(conf.frequency):
				for i := range conf.templateNames {
					generatorLoop(key, orderedParsedTemplates[i], conf.oneline, producer[i])
				}
			case <-ctx.Done():
				stop()
				break Infinite
			}
		}
	} else {
		for i := range conf.templateNames {
			generatorLoop(key, orderedParsedTemplates[i], conf.oneline, producer[i])
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

func configureJrContext(seed int64, num int, frequency time.Duration, locale string, templateDir string) {
	functions2.Random.Seed(seed)
	functions2.JrContext.Num = num
	functions2.JrContext.Range = make([]int, num)
	functions2.JrContext.Frequency = frequency
	functions2.JrContext.Locale = strings.ToLower(locale)
	functions2.JrContext.Seed = seed
	functions2.JrContext.TemplateDir = templateDir
	functions2.JrContext.CountryIndex = functions2.IndexOf(strings.ToUpper(locale), "country")
}

func generatorLoop(key *template.Template, value *template.Template, oneline bool, producer Producer) {
	for range functions2.JrContext.Range {
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

func createKafkaProducer(serializer string, topic []string, index int, kafkaConfig string, schemaRegistry bool, registryConfig string, kcat bool, autocreate bool) *kafka.KafkaManager {

	kManager := &kafka.KafkaManager{
		Serializer:   serializer,
		Topic:        topic[index],
		TemplateType: functions2.JrContext.TemplateType[index],
	}

	kManager.Initialize(kafkaConfig)

	if schemaRegistry {
		kManager.InitializeSchemaRegistry(registryConfig)
		if kcat {
			log.Println("Ignoring kcat when schemaRegistry is enabled")
		}
	}
	if autocreate {
		for i := range topic {
			kManager.CreateTopic(topic[i])
		}
	}
	return kManager
}

func writeStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(functions2.JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", functions2.JrContext.GeneratedObjects)
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", functions2.JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Number of templates (Objects): %d\n", functions2.JrContext.NumTemplates)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(functions2.JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}

func executeTemplate(key *template.Template, value *template.Template, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, functions2.JrContext); err != nil {
		log.Println(err)
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, functions2.JrContext); err != nil {
		log.Println(err)
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	functions2.JrContext.GeneratedObjects++
	functions2.JrContext.GeneratedBytes += int64(len(v))

	return k, v, err
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("num", "n", functions2.JrContext.Num, "Number of elements to create for each pass")
	runCmd.Flags().DurationP("frequency", "f", -1, "how much time to wait for next generation pass")
	runCmd.Flags().DurationP("duration", "d", 0, "If frequency is enabled, with Duration you can set a finite amount of time")

	runCmd.Flags().Int64("seed", functions2.JrContext.Seed, "Seed to init pseudorandom generator")

	runCmd.Flags().String("templateDir", functions2.JrContext.TemplateDir, "directory containing templates")
	runCmd.Flags().StringP("kafkaConfig", "F", "./kafka/config.properties", "Kafka configuration")
	runCmd.Flags().String("registryConfig", "./kafka/registry.properties", "Kafka configuration")
	runCmd.Flags().Bool("templateFileName", false, "If enabled, [template] must be a template file")
	runCmd.Flags().Bool("template", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")

	runCmd.Flags().StringSliceP("preload", "p", []string{""}, "Array of templates to preload")
	runCmd.Flags().IntSlice("preloadSize", []int{}, "Array of template sizes to preload")

	runCmd.Flags().StringP("key", "k", "key", "A template to generate a key")
	runCmd.Flags().StringSliceP("topic", "t", []string{"test"}, "Array of Kafka topic names")

	runCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	runCmd.Flags().StringP("output", "o", "stdout", "can be one of stdout, kafka, redis, mongo")
	runCmd.Flags().String("outputTemplate", "{{.V}}\n", "Formatting of K,V on standard output")
	runCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
	runCmd.Flags().BoolP("autocreate", "a", false, "if enabled, autocreate topics")
	runCmd.Flags().String("locale", functions2.JrContext.Locale, "Locale")

	runCmd.Flags().BoolP("schemaRegistry", "s", false, "If you want to use Confluent Schema Registry")
	runCmd.Flags().String("serializer", "json-schema", "Type of serializer: json-schema, avro-generic, avro, protobuf")
	runCmd.Flags().Duration("redis.ttl", 1*time.Minute, "If output is redis, ttl of the object")
	runCmd.Flags().String("redisConfig", "./redis/config.json", "Redis configuration")

}
