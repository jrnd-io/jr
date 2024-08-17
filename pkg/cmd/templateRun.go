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
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"github.com/ugol/jr/pkg/functions"
)

var templateRunCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template. 
  Without any other flag, [template] is just the name of a template in the templates directory, which is '$JR_SYSTEM_DIR/templates'. Example: 
jr template run net_device
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
		throughputString, _ := cmd.Flags().GetString("throughput")
		seed, _ := cmd.Flags().GetInt64("seed")
		topic, _ := cmd.Flags().GetString("topic")
		preload, _ := cmd.Flags().GetInt("preload")

		csv, _ := cmd.Flags().GetString("csv")

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

		throughput, err := emitter.ParseThroughput(throughputString)
		if err != nil {
			log.Panic().Err(err).Msg("Throughput format error")
		}

		if throughput > 0 {
			// @TODO
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				switch f.Name {
				case "producerConfig":
					configuration.GlobalCfg.ProducerConfig, _ = cmd.Flags().GetString(f.Name)
				case "registryConfig":
					configuration.GlobalCfg.RegistryConfig, _ = cmd.Flags().GetString(f.Name)
				case "autocreate":
					configuration.GlobalCfg.AutoCreate, _ = cmd.Flags().GetBool(f.Name)
				case "schemaRegistry":
					configuration.GlobalCfg.SchemaRegistry, _ = cmd.Flags().GetBool(f.Name)
				case "serializer":
					configuration.GlobalCfg.Serializer, _ = cmd.Flags().GetString(f.Name)
					/*
						case "redisTtl":
							configuration.GlobalCfg.RedisTtl, _ = cmd.Flags().GetDuration(f.Name)
						case "redisConfig":
							configuration.GlobalCfg.RedisConfig, _ = cmd.Flags().GetString(f.Name)
						case "mongoConfig":
							configuration.GlobalCfg.MongoConfig, _ = cmd.Flags().GetString(f.Name)
						case "elasticConfig":
							configuration.GlobalCfg.ElasticConfig, _ = cmd.Flags().GetString(f.Name)
						case "s3Config":
							configuration.GlobalCfg.S3Config, _ = cmd.Flags().GetString(f.Name)
						case "gcsConfig":
							configuration.GlobalCfg.GCSConfig, _ = cmd.Flags().GetString(f.Name)
						case "httpConfig":
							configuration.GlobalCfg.HTTPConfig, _ = cmd.Flags().GetString(f.Name)
					*/
				}
			}
		})

		e := emitter.Emitter{
			Name:             constants.DEFAULT_EMITTER_NAME,
			Locale:           locale,
			Num:              num,
			Frequency:        frequency,
			Duration:         duration,
			Preload:          preload,
			ValueTemplate:    vTemplate,
			EmbeddedTemplate: eTemplate,
			KeyTemplate:      keyTemplate,
			OutputTemplate:   outputTemplate,
			Output:           output,
			Topic:            topic,
			Kcat:             kcat,
			Oneline:          oneline,
			Csv:              csv,
		}

		functions.SetSeed(seed)
		es := map[string][]emitter.Emitter{constants.DEFAULT_EMITTER_NAME: {e}}
		RunEmitters([]string{e.Name}, es, false)
	},
}

func init() {
	templateCmd.AddCommand(templateRunCmd)
	templateRunCmd.Flags().IntP("num", "n", constants.NUM, "Number of elements to create for each pass")
	templateRunCmd.Flags().DurationP("frequency", "f", constants.FREQUENCY, "how much time to wait for next generation pass")
	templateRunCmd.Flags().DurationP("duration", "d", constants.INFINITE, "If frequency is enabled, with Duration you can set a finite amount of time")
	templateRunCmd.Flags().String("throughput", "", "You can set throughput, JR will calculate frequency automatically.")

	templateRunCmd.Flags().Int64("seed", time.Now().UTC().UnixNano(), "Seed to init pseudorandom generator")

	templateRunCmd.Flags().String("csv", "", "Path to csv file to use")

	templateRunCmd.Flags().String("registryConfig", "", "Kafka configuration")
	templateRunCmd.Flags().Bool("embedded", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")
	templateRunCmd.Flags().Int("preload", constants.DEFAULT_PRELOAD_SIZE, "Number of elements to create during the preload phase")

	templateRunCmd.Flags().StringP("key", "k", constants.DEFAULT_KEY, "A template to generate a key")
	templateRunCmd.Flags().StringP("topic", "t", constants.DEFAULT_TOPIC, "Kafka topic")

	templateRunCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	templateRunCmd.Flags().StringP("output", "o", constants.DEFAULT_OUTPUT, "can be one of stdout, kafka, redis, mongo, s3, gcs")
	templateRunCmd.Flags().String("outputTemplate", constants.DEFAULT_OUTPUT_TEMPLATE, "Formatting of K,V on standard output")
	templateRunCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
	templateRunCmd.Flags().BoolP("autocreate", "a", false, "if enabled, autocreate topics")
	templateRunCmd.Flags().String("locale", constants.LOCALE, "Locale")

	templateRunCmd.Flags().BoolP("schemaRegistry", "s", false, "If you want to use Confluent Schema Registry")
	templateRunCmd.Flags().String("serializer", "", "Type of serializer: json-schema, avro-generic, avro, protobuf")
	templateRunCmd.Flags().Duration("redis.ttl", -1, "If output is redis, ttl of the object")

	/*
		templateRunCmd.Flags().String("httpConfig", "", "HTTP configuration")
		templateRunCmd.Flags().String("redisConfig", "", "Redis configuration")
		templateRunCmd.Flags().String("mongoConfig", "", "MongoDB configuration")
		templateRunCmd.Flags().String("elasticConfig", "", "Elastic Search configuration")
		templateRunCmd.Flags().String("s3Config", "", "Amazon S3 configuration")
		templateRunCmd.Flags().String("gcsConfig", "", "Google GCS configuration")
	*/
	//	templateRunCmd.Flags().StringP("kafkaConfig", "F", "", "Kafka configuration")
	templateRunCmd.Flags().String("producerConfig", "", "Producer configuration")

}
