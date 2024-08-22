// Copyright © 2024 JR team
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
package wasm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"

	wazapi "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Producer struct {
	lock sync.Mutex

	r wazero.Runtime
	m wazapi.Module
	f wazapi.Function

	stdin  *bytes.Buffer
	stderr *bytes.Buffer
}

func (p *Producer) Initialize(ctx context.Context, configFile string) {
	cfgBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	config := Config{}
	if err := json.Unmarshal(cfgBytes, &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	p.InitializeFromConfig(ctx, config)
}

func (p *Producer) InitializeFromConfig(ctx context.Context, config Config) {
	p.lock = sync.Mutex{}
	p.r = wazero.NewRuntime(ctx)

	// initialize WASI for stdin/out
	if _, err := wasi.NewBuilder(p.r).Instantiate(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to configure WASI")
	}

	moduleBytes, err := os.ReadFile(config.ModulePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read WASM module file")
	}

	p.stdin = bytes.NewBuffer(nil)
	p.stderr = bytes.NewBuffer(nil)

	mCfg := wazero.NewModuleConfig()
	mCfg = mCfg.WithStdin(p.stdin)
	mCfg = mCfg.WithStderr(p.stderr)

	if config.BindStdout {
		mCfg = mCfg.WithStdout(os.Stdout)
	}

	m, err := p.r.InstantiateWithConfig(ctx, moduleBytes, mCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create WASM module")
	}

	p.m = m
	p.f = p.m.ExportedFunction("produce")
}

func (p *Producer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.stdin.Reset()
	p.stderr.Reset()

	data, err := json.Marshal(map[string][]byte{
		"k": k,
		"v": v,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to serialize WASM request")
	}

	p.stdin.Write(data)
	ret, err := p.f.Call(ctx, uint64(len(data)))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to invoke WASM function")
	}

	if len(ret) == 1 && ret[0] > 0 {
		err = p.extractError(ret[0])
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to execute WASM function")
		}
	}
}

func (p *Producer) Close(ctx context.Context) error {
	if p.r == nil {
		if err := p.r.Close(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Producer) extractError(len uint64) error {
	if len == 0 {
		return nil
	}

	out := p.stderr.Bytes()
	if out == nil {
		return nil
	}

	return errors.New(string(out))
}
