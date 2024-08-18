// Copyright © 2022 Vincenzo Marchese <vincenzo.marchese@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package azblobstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	configuration Config
	client        *azblob.Client
}

func (p *Producer) Initialize(configFile string) {
	cfgBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	config := Config{}
	if err := json.Unmarshal(cfgBytes, &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	if config.AccountName == "" {
		log.Fatal().Msg("AccountName is mandatory")
	}

	if config.PrimaryAccountKey == "" {
		log.Fatal().Msg("PrimaryAccountKey is mandatory")
	}

	p.configuration = config
	cred, err := azblob.NewSharedKeyCredential(config.AccountName, config.PrimaryAccountKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create shared key credential")
	}

	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net", config.AccountName), cred, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create  azure blob storage client")
	}

	if config.Container.Name == "" {
		log.Fatal().Msg("Container name is mandatory")

	}
	if config.Container.Create {
		_, err := client.CreateContainer(context.TODO(), config.Container.Name, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create container")
		}
	}

	p.client = client

}

func (p *Producer) Produce(k []byte, v []byte, _ any) {

	var key string
	if len(k) == 0 || strings.ToLower(string(k)) == "null" {
		// generate a UUID as index
		key = uuid.New().String()
	} else {
		key = string(k)
	}

	resp, err := p.client.UploadBuffer(
		context.TODO(),
		p.configuration.Container.Name,
		key,
		v,
		&azblob.UploadBufferOptions{
			Metadata: map[string]*string{
				"key": &key,
			},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to upload blob")
	}

	log.Trace().Str("key", key).Interface("upload_resp", resp).Msg("Uploaded blob")

}

func (p *Producer) Close() error {
	return nil
}
