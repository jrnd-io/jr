package configuration

import "time"

var GlobalCfg GlobalConfiguration

type GlobalConfiguration struct {
	Seed             int64
	TemplateDir      string
	KafkaConfig      string
	SchemaRegistry   bool
	RegistryConfig   string
	Serializer       string
	AutoCreate       bool
	RedisTtl         time.Duration
	RedisConfig      string
	MongoConfig      string
	ElasticConfig    string
	S3Config         string
	Url              string
	EmbeddedTemplate bool
	FileNameTemplate bool
}
