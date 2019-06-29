package akismet

import (
	"context"
	stderr "errors"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	commentCheckEndpoint = "comment-check"
)

var (
	// ErrUnusualResponse indicates that we got response that we were not expecting, i.e.: comment check can return
	// true or false, but also some error information.
	ErrUnusualResponse = stderr.New("got unusual response")
)

type akismetClient struct {
	key        string
	blogUrl    string
	akismetUrl string
	httpClient *http.Client
}

// NewAkismet returns new instance of Akismet client.
func NewAkismet(key, blogUrl string, optFns ...OptFn) *akismetClient {
	if key == "" {
		panic("API key is required")
	}
	client := &akismetClient{
		blogUrl: blogUrl,
		key:     key,
	}

	defaultOpts(client)

	for _, fn := range optFns {
		fn(client)
	}

	return client
}

// Check call Akismet's check comment endpoint and return true or false along with error that indicates error during process.
func (a *akismetClient) Check(ctx context.Context, c *Comment) (bool, error) {
	if err := c.Validate(); err != nil {
		return false, errors.Wrap(err, "error validating comment struct")
	}
	url := fmt.Sprintf(a.akismetUrl, a.key, commentCheckEndpoint)
	payload := c.toValues()
	respBody, err := a.post(ctx, url, payload)
	if err != nil {
		return true, errors.Wrap(err, "error during comment check request")
	}

	if respBody == "true" {
		return true, nil
	}
	if respBody == "false" {
		return false, nil
	}

	return true, errors.Wrapf(ErrUnusualResponse, "got response: '%s'", respBody)
}
