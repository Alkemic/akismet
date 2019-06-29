package akismet

import (
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
)

func TestWithHttpClient(t *testing.T) {
	httpClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "foo.bar",
		},
	}}

	tests := []struct {
		name             string
		withHttpClientFn []OptFn
		expected         *http.Client
	}{{
		name:             "nil http client",
		withHttpClientFn: []OptFn{WithHttpClient(nil)},
		expected:         nil,
	}, {
		name:     "use default http client",
		expected: &http.Client{},
	}, {
		name:             "basic working example",
		withHttpClientFn: []OptFn{WithHttpClient(httpClient)},
		expected:         httpClient,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewAkismet("asd", "http://some-blog.com", tt.withHttpClientFn...)
			if !reflect.DeepEqual(client.httpClient, tt.expected) {
				t.Errorf("expected http client to be '%v', but got '%v'", tt.expected, client.httpClient)
			}
		})
	}
}
