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

package cassandra

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

type Producer struct {
	configuration Config

	session          *gocql.Session
	consistencyLevel gocql.Consistency
	timeout          time.Duration
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

	if config.Keyspace == "" {
		log.Fatal().Msg("Keyspace is required")
	}

	if len(config.Hosts) == 0 {
		log.Fatal().Msg("Hosts are required")
	}

	if config.Username == "" || config.Password == "" {
		log.Fatal().Msg("Username and password are both required")
	}

	if config.ConsistencyLevel == "" {
		config.ConsistencyLevel = "QUORUM"
	}

	consistencyLevel, err := gocql.MustParseConsistency(config.ConsistencyLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse consistency level")
	}

	if config.Timeout == "" {
		config.Timeout = "10s"
	}

	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse timeout, setting default to 10s")
		timeout = time.Second * 10
	}

	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create session")
	}

	p.configuration = config
	p.timeout = timeout
	p.session = session
	p.consistencyLevel = consistencyLevel

}

func (p *Producer) Produce(_ context.Context, _ []byte, v []byte, _ any) {

	stmt := fmt.Sprintf("INSERT INTO %s.%s JSON ?",
		p.configuration.Keyspace,
		p.configuration.Table)
	if err := p.session.Query(stmt, string(v)).
		Consistency(p.consistencyLevel).Exec(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute query")
	}
}

func (p *Producer) Close(_ context.Context) error {
	p.session.Close()
	return nil
}
