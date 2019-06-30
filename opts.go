package akismet

import "net/http"

const akismetUrl = "https://%s.rest.akismet.com/1.1/%s"

type OptFn func(c *akismetClient)

func defaultOpts(c *akismetClient) {
	c.httpClient = &http.Client{}
	c.akismetUrl = akismetUrl
}

// WithHttpClient is client functional option to set custom httpClient.
func WithHttpClient(httpClient *http.Client) OptFn {
	return func(c *akismetClient) {
		c.httpClient = httpClient
	}
}
