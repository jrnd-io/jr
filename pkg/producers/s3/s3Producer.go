package s3

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AWSRegion string `json:"aws_region"`
	Bucket    string `json:"bucket"`
}

type S3Producer struct {
	client s3.S3
	bucket string
}

func (p *S3Producer) Initialize(configBytes []byte) {
	var config Config
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}

	sess, err := session.NewSession(&aws.Config{Region: &config.AWSRegion})

	if err != nil {
		log.Fatal().Err(err).Msg("Can't establish a session to S3")
		return
	}

	s3Client := s3.New(sess)

	p.client = *s3Client
	p.bucket = config.Bucket
}

func (p *S3Producer) Produce(k []byte, v []byte, o any) {

	bucket := p.bucket
	var key string

	if k == nil || len(k) == 0 {
		// generate a UUID as index
		id := uuid.New()
		key = id.String() + "/.json"
	} else {
		key = string(k) + "/.json"
	}

	buffer := bytes.NewReader(v)

	_, err := p.client.PutObject(&s3.PutObjectInput{
		Body:   buffer,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write data in s3")
	}
}

func (p *S3Producer) Close() error {
	log.Warn().Msg("S3 Client doesn't provide a close method!")
	return nil
}
