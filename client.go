package akismet

import (
	"context"
	stderr "errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	commentCheckEndpoint    = "comment-check"
	keyVerificationEndpoint = "verify-key"
	submitSpamEndpoint      = "submit-spam"
	submitHamEndpoint       = "submit-ham"

	spamHamResponse = "Thanks for making the web a better place."
)

var (
	// ErrUnusualResponse indicates that we got response that we were not expecting, i.e.: comment check can return
	// true or false, but also some error information.
	ErrUnusualResponse = stderr.New("got unusual response")
	// ErrAPIKeyRequired indicates that API key was not provided.
	ErrAPIKeyRequired = stderr.New("API key is required")
	// ErrAPIKeyRequired indicates that blog url was not provided.
	ErrBlogURLRequired = stderr.New("blog url is required")
	// ErrBlogURLIncorrect indicates that provided blog url is not valid, i.e. missing scheme.
	ErrBlogURLIncorrect = stderr.New("incorrect blog url")
)

type akismetClient struct {
	key        string
	blogUrl    string
	akismetUrl string
	httpClient *http.Client
}

// NewAkismet returns new instance of Akismet client with optional error.
func NewAkismet(key, blogUrl string, optFns ...OptFn) (*akismetClient, error) {
	if key == "" {
		return nil, ErrAPIKeyRequired
	}
	if blogUrl == "" {
		return nil, ErrBlogURLRequired
	}
	if _, err := url.ParseRequestURI(blogUrl); err != nil {
		return nil, errors.Wrap(ErrBlogURLIncorrect, err.Error())
	}
	client := &akismetClient{
		blogUrl: blogUrl,
		key:     key,
	}

	defaultOpts(client)

	for _, fn := range optFns {
		fn(client)
	}

	return client, nil
}

// Check calls Akismet's check comment endpoint and return true or false along with error that indicates error during process.
func (a *akismetClient) Check(ctx context.Context, c *Comment) (bool, error) {
	if err := c.Validate(); err != nil {
		return false, errors.Wrap(err, "error validating comment struct")
	}
	commentCheckUrl := fmt.Sprintf(a.akismetUrl, a.key, commentCheckEndpoint)
	payload := c.toValues()
	respBody, err := a.post(ctx, commentCheckUrl, payload)
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

// Verify call Akismet's key verification endpoint and return true or false along with error that indicates error during process.
func (a *akismetClient) Verify(ctx context.Context) (bool, error) {
	payload := &url.Values{}
	payload.Add("key", a.key)
	verifyUrl := fmt.Sprintf(a.akismetUrl, a.key, keyVerificationEndpoint)
	respBody, err := a.post(ctx, verifyUrl, payload)
	if err != nil {
		return false, errors.Wrap(err, "error during comment check request")
	}

	if respBody == "valid" {
		return true, nil
	}
	if respBody == "invalid" {
		return false, nil
	}

	return false, errors.Wrapf(ErrUnusualResponse, "got response: '%s'", respBody)
}

// SubmitSpam calls Akismet's submit spam endpoint and error that indicates error during process.
func (a *akismetClient) SubmitSpam(ctx context.Context, c *Comment) error {
	if err := c.Validate(); err != nil {
		return errors.Wrap(err, "error validating comment struct")
	}
	payload := c.toValues()
	verifyUrl := fmt.Sprintf(a.akismetUrl, a.key, submitSpamEndpoint)
	respBody, err := a.post(ctx, verifyUrl, payload)
	if err != nil {
		return errors.Wrap(err, "error during comment check request")
	}

	if respBody == spamHamResponse {
		return nil
	}

	return errors.Wrapf(ErrUnusualResponse, "got response: '%s'", respBody)
}

// SubmitSpam calls Akismet's submit spam endpoint and error that indicates error during process.
func (a *akismetClient) SubmitHam(ctx context.Context, c *Comment) error {
	if err := c.Validate(); err != nil {
		return errors.Wrap(err, "error validating comment struct")
	}
	payload := c.toValues()
	verifyUrl := fmt.Sprintf(a.akismetUrl, a.key, submitHamEndpoint)
	respBody, err := a.post(ctx, verifyUrl, payload)
	if err != nil {
		return errors.Wrap(err, "error during comment check request")
	}

	if respBody == spamHamResponse {
		return nil
	}

	return errors.Wrapf(ErrUnusualResponse, "got response: '%s'", respBody)
}
