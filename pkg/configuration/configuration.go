package configuration

import "time"

var GlobalCfg GlobalConfiguration

type GlobalConfiguration struct {
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

type Configuration struct {
	TemplateNames    []string
	KeyTemplate      string
	OutputTemplate   string
	EmbeddedTemplate bool
	TemplateFileName bool
	Kcat             bool
	Output           string
	Oneline          bool
	Locale           string
	Num              int
	Frequency        time.Duration
	Duration         time.Duration
	Seed             int64
	KafkaConfig      string
	RegistryConfig   string
	Topic            []string
	Preload          []string
	PreloadSize      []int
	TemplateDir      string
	Autocreate       bool
	SchemaRegistry   bool
	Serializer       string
	RedisTtl         time.Duration
	RedisConfig      string
	MongoConfig      string
	ElasticConfig    string
	S3Config         string
	Url              string
}
