//go:build exclude

package redis

import "testing"

func TestConnectAndClose(t *testing.T) {
	configFile := "config.json.example"

	config, err := InitializeJsonWriter(configFile)
	if err != nil {
		t.Fatalf("Error reading configuration file: %v", err)
	}

	client := RedisClient{config: config}

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

	client := RedisClient{config: options}

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
