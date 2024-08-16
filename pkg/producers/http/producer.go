package http

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Producer struct {
	configuration Config

	certificate tls.Certificate
	client      *resty.Client
	cookiejar   http.CookieJar
}

func (p *Producer) Initialize(configFile string) {
	cfgBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	config := Config{}
	if err := json.Unmarshal(cfgBytes, &config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	p.InitializeFromConfig(config)
}

func (p *Producer) InitializeFromConfig(config Config) {

	var err error
	p.configuration = config
	if p.configuration.Endpoint.Timeout == "" {
		p.configuration.Endpoint.timeout = time.Second * 10
	} else {
		p.configuration.Endpoint.timeout, err = time.ParseDuration(p.configuration.Endpoint.Timeout)
		if err != nil {
			log.Fatalf("Failed to parse timeout: %v", err)
		}

	}

	if p.configuration.ErrorHandling.ExpectStatusCode == 0 {
		p.configuration.ErrorHandling.ExpectStatusCode = 200
	}

	if p.configuration.TLS.CertFile != "" && p.configuration.TLS.KeyFile == "" {
		log.Fatalf("CertFile is set but KeyFile is not")
	}
	if p.configuration.TLS.CertFile == "" && p.configuration.TLS.KeyFile != "" {
		log.Fatalf("KeyFile is set but CertFile is not")
	}

	if p.configuration.TLS.CertFile != "" {
		p.certificate, err = tls.LoadX509KeyPair(p.configuration.TLS.CertFile, p.configuration.TLS.KeyFile)
		if err != nil {
			log.Fatalf("Failed to load certificate: %v", err)
		}
	}

	if p.configuration.Session.UseCookieJar {
		p.cookiejar, err = cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Failed to create cookie jar: %v", err)
		}
	}

	p.client = resty.New().
		SetTimeout(p.configuration.Endpoint.timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: p.configuration.TLS.InsecureSkipVerify,
		}).
		SetHeaders(p.configuration.Headers)

	if p.configuration.Session.UseCookieJar {
		p.client.SetCookieJar(p.cookiejar)
	}

	if p.configuration.TLS.CertFile != "" {
		p.client.SetCertificates(p.certificate)
	}

	if p.configuration.TLS.RootCAFile != "" {
		p.client.SetRootCertificate(p.configuration.TLS.RootCAFile)
	}

	switch p.configuration.Authentication.Type {
	case BasicAuth:
		p.client.SetBasicAuth(p.configuration.Authentication.Basic.Username,
			p.configuration.Authentication.Basic.Password)
	case BearerAuth:
		p.client.SetAuthToken(p.configuration.Authentication.Bearer.Token)
	case APIKeyAuth:
		p.client.SetHeader(p.configuration.Authentication.APIKey.Header,
			p.configuration.Authentication.APIKey.Value)
	case DigestAuth:
		p.client.SetDigestAuth(p.configuration.Authentication.Digest.Username,
			p.configuration.Authentication.Digest.Password)
	default:

	}

	if p.configuration.Endpoint.Method == "" {
		p.configuration.Endpoint.Method = POST
	}

}

func (p *Producer) Produce(k []byte, v []byte, o any) {

	var err error

	// creating request
	req := p.client.R().
		SetBody(v)

	resp := &resty.Response{}
	switch p.configuration.Endpoint.Method {
	case POST:
		resp, err = req.Post(p.configuration.Endpoint.URL)
	case PUT:
		resp, err = req.Put(p.configuration.Endpoint.URL)
	default:
		resp, err = req.Post(p.configuration.Endpoint.URL)
	}

	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode() != p.configuration.ErrorHandling.ExpectStatusCode &&
		!p.configuration.ErrorHandling.IgnoreStatusCode {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode())
	}

}

func (p *Producer) Close() error {
	return nil
}

func (p *Producer) GetClient() *resty.Client {
	return p.client
}
