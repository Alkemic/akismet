package akismet

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// ErrNonOKStatusCode returned when Akismet API returns non OK status, which shouldn't happen on normal usage.
var ErrNonOKStatusCode = errors.New("akismet API returned non 200 status code")

func (a *akismetClient) post(ctx context.Context, url string, payload *url.Values) (string, error) {
	payload.Add("blog", a.blogUrl)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "error creating HTTP request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := a.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "cannot do HTTP request")
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrapf(ErrNonOKStatusCode, "got status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}

	return string(body), nil
}
