// Copyright © 2024 JR team
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

package wamp

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/rs/zerolog/log"
)

type Config struct {
	WampURI  string `json:"wamp_uri"`
	Username string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
	Topic    string `json:"topic"`
	SerType  string `json:"serType"`
	Compress bool   `json:"compress"`
	Authid   string `json:"authid"`
}

type Producer struct {
	client client.Client
	realm  string
	topic  string
	authid string
}

func (p *Producer) Initialize(ctx context.Context, configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}
	var wampclient *client.Client

	// Get requested serialization.
	serialization := client.JSON
	switch config.SerType {
	case "json":
	case "msgpack":
		serialization = client.MSGPACK
	case "cbor":
		serialization = client.CBOR
	default:
		log.Fatal().Err(err).Msg("Invalid serialization, muse be one of: json, msgpack, cbor")
	}

	cfg := client.Config{
		Realm:         config.Realm,
		Serialization: serialization,
		HelloDetails: wamp.Dict{
			"authid": config.Authid,
		},
	}

	if config.Compress {
		cfg.WsCfg.EnableCompression = true
	}

	addr := config.WampURI

	wampclient, err = client.ConnectNet(context.Background(), addr, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to WAMP Router")
	}
	// defer wampclient.Close()

	p.realm = config.Realm
	p.topic = config.Topic
	p.authid = config.Authid

	p.client = *wampclient
}

func (p *Producer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	data := string(v)
	args := wamp.List{data}
	opts := wamp.Dict{
		"authid": p.authid,
	}
	err := p.client.Publish(p.topic, opts, args, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("publish error: %s", err)
	}
}

func (p *Producer) Close(ctx context.Context) error {
	err := p.client.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close WAMP connection")
	}
	return err
}
