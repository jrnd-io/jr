//go:build exclude

package redis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const CRLF = "\r\n"
const JSON_GET = "JSON.GET"
const JSON_SET = "JSON.SET"
const DOLLAR = "$"
const STAR = "*"
const NEW_LINE = '\n'
const PLUS = "+"

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisClient struct {
	config RedisConfig
	conn   net.Conn
}

func Initialize(filename string) (RedisConfig, error) {
	var config RedisConfig
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (r *RedisClient) Connect() error {
	addr := fmt.Sprintf("%s:%s", r.config.Host, r.config.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	if r.config.Password != "" {
		authCmd := fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(r.config.Password), r.config.Password)
		if r.config.Username != "" {
			authCmd = fmt.Sprintf("*3\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(r.config.Username), r.config.Username, len(r.config.Password), r.config.Password)
		}
		_, err = conn.Write([]byte(authCmd))
		if err != nil {
			fmt.Println("Error sending AUTH command:", err)
			return err
		}
		reader := bufio.NewReader(conn)
		authResponse, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading AUTH response:", err)
			return err
		}

		if strings.HasPrefix(authResponse, "-") {
			msg := fmt.Sprint("Authentication failed:", strings.TrimSpace(authResponse[1:]))
			fmt.Println(msg)
			return fmt.Errorf(msg)
		}
	}

	r.conn = conn
	return nil
}

func (r *RedisClient) Disconnect() {
	r.conn.Close()
}

func (r *RedisClient) Set(key, value, path string) error {
	args := make([]string, 4)
	args[0] = JSON_SET
	args[1] = key
	args[2] = string(DOLLAR)
	args[3] = value

	if path != "" {
		args[2] = path
	}

	setCmd := fmt.Sprint(STAR, len(args), CRLF)
	for _, value := range args {
		setCmd = fmt.Sprint(setCmd, DOLLAR, len(value), CRLF, value, CRLF)
	}

	_, err := r.conn.Write([]byte(setCmd))
	if err != nil {
		return err
	}

	reader := bufio.NewReader(r.conn)
	response, err := reader.ReadString(NEW_LINE)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, PLUS) {
		msg := fmt.Sprintf("Unexpected response from Redis: %s", response)
		return fmt.Errorf(msg)
	}

	return nil
}

func (r *RedisClient) Get(key, path string) (string, error) {
	args := make([]string, 2, 3)
	args[0] = JSON_GET
	args[1] = key

	if path != "" {
		args[2] = path
	}

	setCmd := fmt.Sprint(STAR, len(args), CRLF)
	for _, value := range args {
		setCmd = fmt.Sprint(setCmd, DOLLAR, len(value), CRLF, value, CRLF)
	}

	_, err := r.conn.Write([]byte(setCmd))
	if err != nil {
		return "Error while invoking command", err
	}

	reader := bufio.NewReader(r.conn)
	response, err := reader.ReadString(NEW_LINE)
	if err != nil {
		return "Error while reading Redis response", err
	}

	if !strings.HasPrefix(response, DOLLAR) {
		msg := fmt.Sprintf("Unexpected response from Redis: %s", response)
		return msg, fmt.Errorf(msg)
	}

	valueLen, err := strconv.Atoi(response[1 : len(response)-2])
	if err != nil {
		return "", err
	}

	valueBytes := make([]byte, valueLen+2)
	_, err = reader.Read(valueBytes)
	if err != nil {
		return "", err
	}
	return string(valueBytes[:valueLen]), nil
}
