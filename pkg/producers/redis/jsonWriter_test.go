//go:build exclude

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

package redis

import "testing"

func TestConnectAndClose(t *testing.T) {
	configFile := "config.json.example"

	config, err := InitializeJsonWriter(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	client := Client{config: config}

	err = client.Connect()
	if err != nil {
		t.Fatalf("Error connecting to Redis: %v", err)
	}
	defer client.Disconnect()
}

func TestSetAndGet(t *testing.T) {
	configFile := "config.json.example"

	options, err := InitializeJsonWriter(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	client := Client{config: options}

	err = client.Connect()
	if err != nil {
		t.Fatalf("Error connecting to Redis: %v", err)
	}
	defer client.Disconnect()

	key := "loo:foo"
	value := "{\"firstname\":\"Luigi\",\"lastname\":\"Fugaro\"}"

	err = client.Set(key, value, "")
	if err != nil {
		t.Fatalf("Error setting key-value: %v", err)
	}

	retrievedValue, err := client.Get(key, "")
	if err != nil {
		t.Fatalf("Error getting key-value: %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected value '%s', got '%s'", value, retrievedValue)
	}
}
