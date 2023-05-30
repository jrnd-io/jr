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
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"github.com/ugol/jr/pkg/functions"
	"os"
	"time"
)

var templateRunCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template. 
  Without any other flag, [template] is just the name of a template in the templates directory, which by default is in '$HOME/.jr/templates' Example: 
jr template run net-device
  With the --embedded flag, [template] is a string containing a full template. Example:
jr template run --template "{{name}}"
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		keyTemplate, _ := cmd.Flags().GetString("key")
		outputTemplate, _ := cmd.Flags().GetString("outputTemplate")
		embeddedTemplate, _ := cmd.Flags().GetBool("embedded")
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
		topic, _ := cmd.Flags().GetString("topic")
		preload, _ := cmd.Flags().GetInt("preload")

		templateDir, _ := cmd.Flags().GetString("templateDir")
		templateDir = os.ExpandEnv(templateDir)

		autocreate, _ := cmd.Flags().GetBool("autocreate")
		schemaRegistry, _ := cmd.Flags().GetBool("schemaRegistry")
		serializer, _ := cmd.Flags().GetString("serializer")

		redisTtl, _ := cmd.Flags().GetDuration("redis.ttl")
		redisConfig, _ := cmd.Flags().GetString("redisConfig")
		mongoConfig, _ := cmd.Flags().GetString("mongoConfig")
		elasticConfig, _ := cmd.Flags().GetString("elasticConfig")
		s3Config, _ := cmd.Flags().GetString("s3Config")

		if kcat {
			oneline = true
			output = "stdout"
			outputTemplate = constants.DEFAULT_OUTPUT_KCAT_TEMPLATE
		}

		var vTemplate, eTemplate string
		if embeddedTemplate {
			vTemplate = ""
			eTemplate = args[0]
		} else {
			vTemplate = args[0]
			eTemplate = ""
		}

		configuration.GlobalCfg = configuration.GlobalConfiguration{
			Seed:           seed,
			AutoCreate:     autocreate,
			KafkaConfig:    kafkaConfig,
			RegistryConfig: registryConfig,
			TemplateDir:    templateDir,
			SchemaRegistry: schemaRegistry,
			Serializer:     serializer,
			RedisTtl:       redisTtl,
			RedisConfig:    redisConfig,
			MongoConfig:    mongoConfig,
			ElasticConfig:  elasticConfig,
			S3Config:       s3Config,
		}

		e := emitter.Emitter{
			Name:             name,
			Locale:           locale,
			Num:              num,
			Frequency:        frequency,
			Duration:         duration,
			Preload:          preload,
			ValueTemplate:    vTemplate,
			EmbeddedTemplate: eTemplate,
			KeyTemplate:      keyTemplate,
			OutputTemplate:   outputTemplate, //@TODO FIX
			Output:           output,
			Topic:            topic,
			Kcat:             kcat,
			Oneline:          oneline,
		}

		functions.Random.Seed(seed)
		uuid.SetRand(functions.Random)

		es := []emitter.Emitter{e}

		emitter.Initialize([]string{e.Name}, es)
		emitter.DoLoop(es)
		emitter.CloseProducers(es)
		time.Sleep(100 * time.Millisecond)
		emitter.WriteStats()

	},
}

func init() {
	templateCmd.AddCommand(templateRunCmd)
	templateRunCmd.Flags().IntP("num", "n", constants.NUM, "Number of elements to create for each pass")
	templateRunCmd.Flags().DurationP("frequency", "f", constants.FREQUENCY, "how much time to wait for next generation pass")
	templateRunCmd.Flags().DurationP("duration", "d", constants.DURATION, "If frequency is enabled, with Duration you can set a finite amount of time")

	templateRunCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")

	templateRunCmd.Flags().String("templateDir", os.ExpandEnv(constants.TEMPLATEDIR), "directory containing templates")
	templateRunCmd.Flags().StringP("kafkaConfig", "F", constants.KAFKA_CONFIG, "Kafka configuration")
	templateRunCmd.Flags().String("registryConfig", constants.REGISTRY_CONFIG, "Kafka configuration")
	templateRunCmd.Flags().Bool("embedded", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")
	templateRunCmd.Flags().Int("preload", constants.DEFAULT_PRELOAD_SIZE, "Number of elements to create during the preload phase")

	templateRunCmd.Flags().StringP("key", "k", constants.DEFAULT_KEY, "A template to generate a key")
	templateRunCmd.Flags().StringP("topic", "t", constants.DEFAULT_TOPIC, "Kafka topic")

	templateRunCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	templateRunCmd.Flags().StringP("output", "o", constants.DEFAULT_OUTPUT, "can be one of stdout, kafka, redis, mongo")
	templateRunCmd.Flags().String("outputTemplate", constants.DEFAULT_OUTPUT_TEMPLATE, "Formatting of K,V on standard output")
	templateRunCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
	templateRunCmd.Flags().BoolP("autocreate", "a", false, "if enabled, autocreate topics")
	templateRunCmd.Flags().String("locale", constants.LOCALE, "Locale")

	templateRunCmd.Flags().BoolP("schemaRegistry", "s", false, "If you want to use Confluent Schema Registry")
	templateRunCmd.Flags().String("serializer", constants.DEFAULT_SERIALIZER, "Type of serializer: json-schema, avro-generic, avro, protobuf")
	templateRunCmd.Flags().Duration("redis.ttl", constants.REDIS_TTL, "If output is redis, ttl of the object")
	templateRunCmd.Flags().String("redisConfig", constants.REDIS_CONFIG, "Redis configuration")
	templateRunCmd.Flags().String("mongoConfig", constants.MONGO_CONFIG, "MongoDB configuration")
	templateRunCmd.Flags().String("elasticConfig", constants.ELASTIC_CONFIG, "Elastic Search configuration")
	templateRunCmd.Flags().String("s3Config", constants.S3_CONFIG, "Amazon S3 configuration")

}
