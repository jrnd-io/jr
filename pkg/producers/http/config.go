package http

import "time"

type AuthType string
type Method string

const (
	BasicAuth  AuthType = "basic"
	BearerAuth AuthType = "bearer"
	APIKeyAuth AuthType = "api_key"
	DigestAuth AuthType = "digest"

	POST Method = "POST"
	PUT  Method = "PUT"
)

type Endpoint struct {
	URL     string `json:"url"`
	Method  Method `json:"method"`
	Timeout string `json:"timeout"`
	timeout time.Duration
}

type Session struct {
	UseCookieJar bool `json:"use_cookie_jar"`
}

type ErrorHandling struct {
	ExpectStatusCode int  `json:"expect_status_code"`
	IgnoreStatusCode bool `json:"ignore_status_code"`
}

type Headers map[string]string
type TLS struct {
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	CertFile           string `json:"cert_file"`
	KeyFile            string `json:"key_file"`
	RootCAFile         string `json:"root_ca_file"`
}

type APIKey struct {
	Header string `json:"header"`
	Value  string `json:"Value"`
}

type Bearer struct {
	Token string `json:"token"`
}

type Basic struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Authentication struct {
	Type   AuthType `json:"type"`
	Basic  Basic    `json:"basic"`
	Digest Basic    `json:"digest"`
	Bearer Bearer   `json:"bearer"`
	APIKey APIKey   `json:"api_key"`
}

type Config struct {
	Endpoint       Endpoint       `json:"endpoint"`
	Session        Session        `json:"session"`
	ErrorHandling  ErrorHandling  `json:"error_handling"`
	Headers        Headers        `json:"headers"`
	TLS            TLS            `json:"tls"`
	Authentication Authentication `json:"authentication"`
}
