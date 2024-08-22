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

package http_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	phttp "github.com/jrnd-io/jr/pkg/producers/http"
)

var defaultBody = []byte("{\"property\": \"value\"}")

type mockResponder struct {
	name          string
	t             *testing.T
	expectHeaders map[string]string
	status        int
	basic         string
	bearer        string
	apikey        string
}

func (m *mockResponder) serveHTTP(req *http.Request) (*http.Response, error) {

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		m.t.Errorf("%s: cannot read request body", m.name)
	}
	if diff := cmp.Diff(defaultBody, body); diff != "" {
		m.t.Errorf("%s: mismatch challenge (-want +got):\n%s", m.name, diff)
	}

	if m.expectHeaders != nil {
		// flatten request headers
		reqHeaders := make(map[string]string)
		for k, v := range req.Header {
			// testing only  headers starting with "Test-jr"
			if strings.HasPrefix(k, "Test-Jr") {
				reqHeaders[k] = v[0]
			}
		}

		// canonical names for headers
		wantHeaders := make(map[string]string)
		for k, v := range m.expectHeaders {
			wantHeaders[http.CanonicalHeaderKey(k)] = v
		}
		if diff := cmp.Diff(wantHeaders, reqHeaders); diff != "" {
			m.t.Errorf("%s: mismatch challenge (-want +got):\n%s", m.name, diff)
		}
	}

	// check auth
	if m.basic != "" {
		header := req.Header.Get("Authorization")
		if diff := cmp.Diff(header, fmt.Sprintf("Basic %s", m.basic)); diff != "" {
			m.t.Errorf("%s: mismatch challenge (-want +got):\n%s", m.name, diff)
		}
	}
	if m.bearer != "" {
		header := req.Header.Get("Authorization")
		if diff := cmp.Diff(header, fmt.Sprintf("Bearer %s", m.bearer)); diff != "" {
			m.t.Errorf("%s: mismatch challenge (-want +got):\n%s", m.name, diff)
		}
	}
	if m.apikey != "" {
		header := req.Header.Get("api-key")
		if diff := cmp.Diff(header, m.apikey); diff != "" {
			m.t.Errorf("%s: mismatch challenge (-want +got):\n%s", m.name, diff)
		}
	}
	return httpmock.NewStringResponse(m.status, ""), nil

}

func TestProducer(t *testing.T) {
	fakeUrl := "https://jr.io"

	testCases := []struct {
		name    string
		config  phttp.Config
		headers map[string]string
		apiKey  string
		bearer  string
		basic   string
		status  int
	}{
		{
			name: "test_simple_PUT",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.PUT,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "test_simple_POST",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
			},
			status: http.StatusOK,
		},
		{
			name: "test_ignore_status_code",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				ErrorHandling: phttp.ErrorHandling{
					IgnoreStatusCode: true,
				},
			},
			status: http.StatusTeapot,
		},
		{
			name: "test_expect_status_code",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				ErrorHandling: phttp.ErrorHandling{
					ExpectStatusCode: http.StatusTeapot,
				},
			},
			status: http.StatusTeapot,
		},
		{
			name: "test_with_headers",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				Headers: map[string]string{
					"Test-Jrheader01": "value01",
					"Test-Jrheader02": "value02",
				},
			},
			status: http.StatusOK,
			headers: map[string]string{
				"test-jrheader01": "value01",
				"test-jrheader02": "value02",
			},
		},
		{
			name: "test_with_basic",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				Authentication: phttp.Authentication{
					Type: phttp.BasicAuth,
					Basic: phttp.Basic{
						Username: "user",
						Password: "password",
					},
				},
			},
			status: http.StatusOK,
			basic:  base64.StdEncoding.EncodeToString([]byte("user:password")),
		},
		{
			name: "test_with_bearer",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				Authentication: phttp.Authentication{
					Type: phttp.BearerAuth,
					Bearer: phttp.Bearer{
						Token: "sometoken",
					},
				},
			},
			status: http.StatusOK,
			bearer: "sometoken",
		},
		{
			name: "test_with_api_key",
			config: phttp.Config{
				Endpoint: phttp.Endpoint{
					URL:    fakeUrl,
					Method: phttp.POST,
				},
				Authentication: phttp.Authentication{
					Type: phttp.APIKeyAuth,
					APIKey: phttp.APIKey{
						Header: "api-key",
						Value:  "jrapikey",
					},
				},
			},
			status: http.StatusOK,
			apiKey: "jrapikey",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			producer := phttp.Producer{}
			producer.InitializeFromConfig(tc.config)
			httpmock.ActivateNonDefault(producer.GetClient().GetClient())
			httpmock.Reset()

			mr := &mockResponder{
				name:          tc.name,
				t:             t,
				expectHeaders: tc.headers,
				status:        tc.status,
				basic:         tc.basic,
				bearer:        tc.bearer,
				apikey:        tc.apiKey,
			}
			httpmock.RegisterResponder(string(tc.config.Endpoint.Method),
				fakeUrl,
				mr.serveHTTP)

			producer.Produce(context.TODO(), []byte("key"), defaultBody, nil)
			httpmock.DeactivateAndReset()
		})
	}

}
