package s3

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type Config struct {
	AWSRegion string `json:"aws_region"`
	Bucket    string `json:"bucket"`
}

type S3Producer struct {
	client s3.S3
	bucket string
}

func (p *S3Producer) Initialize(configFile string) {
	var config Config
	file, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to parse configuration parameters: %s", err)
	}

	sess, err := session.NewSession(&aws.Config{Region: &config.AWSRegion})

	if err != nil {
		log.Fatalf("Can't establish a session to S3: %s", err)
		return
	}

	s3Client := s3.New(sess)

	p.client = *s3Client
	p.bucket = config.Bucket
}

func (p *S3Producer) Produce(k []byte, v []byte, o interface{}) {

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
		log.Fatalf("Failed to write data in s3:\n%s", err)
	}
}

func (p *S3Producer) Close() error {
	log.Println("S3 Client doesn't provide a close method!")
	return nil
}
