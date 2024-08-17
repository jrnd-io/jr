//Copyright © 2022 Vincenzo Marchese <vincenzo.marchese@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

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
