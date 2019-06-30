package akismet

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestNewAkismet(t *testing.T) {
	type check func(client *akismetClient, err error, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(_ *akismetClient, err error, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(_ *akismetClient, err error, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(_ *akismetClient, err error, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}
	hasClient := func(expClient *akismetClient) check {
		return func(client *akismetClient, _ error, t *testing.T) {
			t.Helper()
			if !reflect.DeepEqual(expClient, client) {
				t.Errorf("Expected Akismet client to be '%v', but got '%v'", expClient, client)
			}
		}
	}

	tests := []struct {
		name    string
		key     string
		blogUrl string
		checks  []check
	}{{
		name: "error when no key provided",
		checks: checks(
			hasCauseError(ErrAPIKeyRequired),
			hasErrorMsg("API key is required"),
			hasClient(nil),
		),
	}, {
		name: "error when no blog url provided",
		key:  "deadbeef",
		checks: checks(
			hasCauseError(ErrBlogURLRequired),
			hasErrorMsg("blog url is required"),
			hasClient(nil),
		),
	}, {
		name:    "error when invalid (not parsable) blog url provided",
		key:     "deadbeef",
		blogUrl: "deadbeef",
		checks: checks(
			hasCauseError(ErrBlogURLIncorrect),
			hasErrorMsg("parse deadbeef: invalid URI for request: incorrect blog url"),
			hasClient(nil),
		),
	}, {
		name:    "success",
		key:     "deadbeef",
		blogUrl: "http://some-blog.com",
		checks: checks(
			hasNoError,
			hasClient(&akismetClient{
				key:        "deadbeef",
				blogUrl:    "http://some-blog.com",
				akismetUrl: "https://%s.rest.akismet.com/1.1/%s",
				httpClient: &http.Client{},
			}),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := NewAkismet(tt.key, tt.blogUrl)
			for _, ch := range tt.checks {
				ch(cli, err, t)
			}
		})
	}
}

func TestAkismetCheck(t *testing.T) {
	type check func(result bool, err error, payload []byte, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(_ bool, err error, payload []byte, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(_ bool, err error, payload []byte, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(_ bool, err error, payload []byte, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}
	hasResult := func(exp bool) check {
		return func(result bool, _ error, payload []byte, t *testing.T) {
			t.Helper()
			if result != exp {
				t.Errorf("Expected result to be '%t', but got '%t'", exp, result)
			}
		}
	}
	hasPayload := func(expPayload string) check {
		return func(_ bool, _ error, payload []byte, t *testing.T) {
			t.Helper()
			if string(payload) != expPayload {
				t.Errorf("Expected requst payload to be \n'%s', but got \n'%s'", expPayload, string(payload))
			}
		}
	}

	validComment := &Comment{
		UserIP:    "0.0.0.0",
		UserAgent: "Mozilla/6.16",
	}
	filledComment := &Comment{
		UserIP:                 "1",
		UserAgent:              "2",
		Referrer:               "3",
		Permalink:              "4",
		CommentType:            "5",
		CommentAuthor:          "6",
		CommentAuthorEmail:     "7",
		CommentAuthorURL:       "8",
		CommentContent:         "9",
		BlogLang:               "10",
		BlogCharset:            "11",
		UserRole:               "12",
		CommentDateGMT:         time.Date(2019, 6, 30, 13, 43, 12, 0, time.UTC).Format(time.RFC3339),
		CommentPostModifiedGMT: time.Date(2019, 6, 30, 14, 43, 12, 0, time.UTC).Format(time.RFC3339),
		IsTest:                 "13",
		RecheckReason:          "14",
	}

	tests := []struct {
		name               string
		comment            *Comment
		responseBody       string
		responseStatusCode int
		checks             []check
	}{{
		name:               "success with 'true' as a response",
		comment:            validComment,
		responseBody:       "true",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(true),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "success with 'false' as a response",
		comment:            validComment,
		responseBody:       "false",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(false),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "success with full serialization",
		comment:            filledComment,
		responseBody:       "true",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(true),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&blog_charset=11&blog_lang=10&comment_author=6&comment_author_email=7&comment_author_url=8&comment_content=9&comment_date_gmt=2019-06-30T13%3A43%3A12Z&comment_post_modified_gmt=2019-06-30T14%3A43%3A12Z&comment_type=5&is_test=13&permalink=4&recheck_reason=14&referrer=3&user_agent=2&user_ip=1&user_role=12"),
		),
	}, {
		name:               "error when status code is not OK",
		comment:            validComment,
		responseStatusCode: 418,
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("error during comment check request: got status code 418: akismet API returned non 200 status code"),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "error when got unusual response",
		comment:            validComment,
		responseBody:       "Missing required field: user_ip.",
		responseStatusCode: 200,
		checks: checks(
			hasCauseError(ErrUnusualResponse),
			hasErrorMsg("got response: 'Missing required field: user_ip.': got unusual response"),
		),
	}, {
		name:         "error when comment validation fails",
		comment:      &Comment{},
		responseBody: "Missing required field: user_ip.",
		checks: checks(
			hasErrorMsg("error validating comment struct: field user ip is required"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/deadbeef/comment-check" {
					t.Errorf("Unexpected endpoint used: %s", r.URL.Path)
				}
				defer r.Body.Close()
				var err error
				if buffer, err = ioutil.ReadAll(r.Body); err != nil {
					t.Fatalf("got error reading payload body from request: '%v'", err)
				}

				w.WriteHeader(tt.responseStatusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer ts.Close()

			cli := &akismetClient{
				key:        "deadbeef",
				blogUrl:    "http://some-blog.com",
				httpClient: &http.Client{},
				akismetUrl: ts.URL + "/%s/%s",
			}
			result, err := cli.Check(context.Background(), tt.comment)
			for _, ch := range tt.checks {
				ch(result, err, buffer, t)
			}
		})
	}
}

func TestAkismetVerify(t *testing.T) {
	type check func(result bool, err error, payload []byte, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(_ bool, err error, payload []byte, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(_ bool, err error, payload []byte, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(_ bool, err error, payload []byte, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}
	hasResult := func(exp bool) check {
		return func(result bool, _ error, payload []byte, t *testing.T) {
			t.Helper()
			if result != exp {
				t.Errorf("Expected result to be '%t', but got '%t'", exp, result)
			}
		}
	}
	hasPayload := func(expPayload string) check {
		return func(_ bool, _ error, payload []byte, t *testing.T) {
			t.Helper()
			if string(payload) != expPayload {
				t.Errorf("Expected requst payload to be '%s', but got '%s'", expPayload, string(payload))
			}
		}
	}

	tests := []struct {
		name               string
		responseBody       string
		responseStatusCode int
		checks             []check
	}{{
		name:               "success when 'valid' as a response",
		responseBody:       "valid",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(true),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&key=deadbeef"),
		),
	}, {
		name:               "error when 'invalid' as a response",
		responseBody:       "invalid",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(false),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&key=deadbeef"),
		),
	}, {
		name:               "error when status code is not OK",
		responseStatusCode: 418,
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("error during comment check request: got status code 418: akismet API returned non 200 status code"),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&key=deadbeef"),
		),
	}, {
		name:               "error when got unusual response",
		responseBody:       "some response",
		responseStatusCode: 200,
		checks: checks(
			hasCauseError(ErrUnusualResponse),
			hasErrorMsg("got response: 'some response': got unusual response"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/deadbeef/verify-key" {
					t.Errorf("Unexpected endpoint used: %s", r.URL.Path)
				}
				defer r.Body.Close()
				var err error
				if buffer, err = ioutil.ReadAll(r.Body); err != nil {
					t.Fatalf("got error reading payload body from request: '%v'", err)
				}

				w.WriteHeader(tt.responseStatusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer ts.Close()

			cli := &akismetClient{
				key:        "deadbeef",
				blogUrl:    "http://some-blog.com",
				httpClient: &http.Client{},
				akismetUrl: ts.URL + "/%s/%s",
			}
			result, err := cli.Verify(context.Background())
			for _, ch := range tt.checks {
				ch(result, err, buffer, t)
			}
		})
	}
}

func TestAkismetSubmitSpam(t *testing.T) {
	type check func(err error, payload []byte, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(err error, payload []byte, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(err error, payload []byte, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(err error, payload []byte, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}
	hasPayload := func(expPayload string) check {
		return func(_ error, payload []byte, t *testing.T) {
			t.Helper()
			if string(payload) != expPayload {
				t.Errorf("Expected requst payload to be '%s', but got '%s'", expPayload, string(payload))
			}
		}
	}

	validComment := &Comment{
		UserIP:    "0.0.0.0",
		UserAgent: "Mozilla/6.16",
	}

	validResponse := "Thanks for making the web a better place."

	tests := []struct {
		name               string
		comment            *Comment
		responseBody       string
		responseStatusCode int
		checks             []check
	}{{
		name:               "success when got valid response",
		comment:            validComment,
		responseBody:       validResponse,
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "error when status code is not OK",
		comment:            validComment,
		responseStatusCode: 418,
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("error during comment check request: got status code 418: akismet API returned non 200 status code"),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "error when got unusual response",
		comment:            validComment,
		responseBody:       "Missing required field: user_ip.",
		responseStatusCode: 200,
		checks: checks(
			hasCauseError(ErrUnusualResponse),
			hasErrorMsg("got response: 'Missing required field: user_ip.': got unusual response"),
		),
	}, {
		name:         "error when comment validation fails",
		comment:      &Comment{},
		responseBody: "Missing required field: user_ip.",
		checks: checks(
			hasErrorMsg("error validating comment struct: field user ip is required"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/deadbeef/submit-spam" {
					t.Errorf("Unexpected endpoint used: %s", r.URL.Path)
				}
				defer r.Body.Close()
				var err error
				if buffer, err = ioutil.ReadAll(r.Body); err != nil {
					t.Fatalf("got error reading payload body from request: '%v'", err)
				}

				w.WriteHeader(tt.responseStatusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer ts.Close()

			cli := &akismetClient{
				key:        "deadbeef",
				blogUrl:    "http://some-blog.com",
				httpClient: &http.Client{},
				akismetUrl: ts.URL + "/%s/%s",
			}
			err := cli.SubmitSpam(context.Background(), tt.comment)
			for _, ch := range tt.checks {
				ch(err, buffer, t)
			}
		})
	}
}

func TestAkismetSubmitHam(t *testing.T) {
	type check func(err error, payload []byte, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(err error, payload []byte, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(err error, payload []byte, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(err error, payload []byte, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}
	hasPayload := func(expPayload string) check {
		return func(_ error, payload []byte, t *testing.T) {
			t.Helper()
			if string(payload) != expPayload {
				t.Errorf("Expected requst payload to be '%s', but got '%s'", expPayload, string(payload))
			}
		}
	}

	validComment := &Comment{
		UserIP:    "0.0.0.0",
		UserAgent: "Mozilla/6.16",
	}

	validResponse := "Thanks for making the web a better place."

	tests := []struct {
		name               string
		comment            *Comment
		responseBody       string
		responseStatusCode int
		checks             []check
	}{{
		name:               "success when got valid response",
		comment:            validComment,
		responseBody:       validResponse,
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "error when status code is not OK",
		comment:            validComment,
		responseStatusCode: 418,
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("error during comment check request: got status code 418: akismet API returned non 200 status code"),
			hasPayload("blog=http%3A%2F%2Fsome-blog.com&user_agent=Mozilla%2F6.16&user_ip=0.0.0.0"),
		),
	}, {
		name:               "error when got unusual response",
		comment:            validComment,
		responseBody:       "Missing required field: user_ip.",
		responseStatusCode: 200,
		checks: checks(
			hasCauseError(ErrUnusualResponse),
			hasErrorMsg("got response: 'Missing required field: user_ip.': got unusual response"),
		),
	}, {
		name:         "error when comment validation fails",
		comment:      &Comment{},
		responseBody: "Missing required field: user_ip.",
		checks: checks(
			hasErrorMsg("error validating comment struct: field user ip is required"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/deadbeef/submit-ham" {
					t.Errorf("Unexpected endpoint used: %s", r.URL.Path)
				}
				defer r.Body.Close()
				var err error
				if buffer, err = ioutil.ReadAll(r.Body); err != nil {
					t.Fatalf("got error reading payload body from request: '%v'", err)
				}

				w.WriteHeader(tt.responseStatusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer ts.Close()

			cli := &akismetClient{
				key:        "deadbeef",
				blogUrl:    "http://some-blog.com",
				httpClient: &http.Client{},
				akismetUrl: ts.URL + "/%s/%s",
			}
			err := cli.SubmitHam(context.Background(), tt.comment)
			for _, ch := range tt.checks {
				ch(err, buffer, t)
			}
		})
	}
}

type transportMock struct {
	roundTripResp *http.Response
	roundTripErr  error
}

func (m *transportMock) RoundTrip(*http.Request) (*http.Response, error) {
	return m.roundTripResp, m.roundTripErr
}

type nopCloserError struct{}

func (nopCloserError) Read(p []byte) (n int, err error) {
	return 0, errors.New("mocked error from nop closer")
}

func (nopCloserError) Close() error { return nil }

func TestAkismetMaliciousResponses(t *testing.T) {
	type check func(err error, t *testing.T)
	checks := func(cs ...check) []check { return cs }

	hasCauseError := func(exp error) check {
		return func(err error, t *testing.T) {
			t.Helper()
			if errors.Cause(err) != exp {
				t.Errorf("Expected error cause to be '%v', but got '%v'", exp, err)
			}
		}
	}
	hasErrorMsg := func(expMsg string) check {
		return func(err error, t *testing.T) {
			t.Helper()
			if err == nil || err.Error() != expMsg {
				t.Errorf("Expected error cause to be '%v', but got '%v'", expMsg, err)
			}
		}
	}
	hasNoError := func(err error, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected error to be nil, but got '%v'", err)
		}
	}

	tests := []struct {
		name       string
		httpClient *http.Client
		url        string
		checks     []check
	}{{
		name: "success when 'valid' as a response",
		httpClient: &http.Client{
			Transport: &transportMock{
				roundTripResp: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("ok")),
				},
			},
		},
		checks: checks(
			hasNoError,
		),
	}, {
		name: "error when processing new request",
		url:  "^!@#$%^&*()",
		checks: checks(
			hasErrorMsg(`error creating HTTP request: parse ^!@#$%^&*(): invalid URL escape "%^&"`),
		),
	}, {
		name: "error when can't do request",
		httpClient: &http.Client{
			Transport: &transportMock{
				roundTripErr: errors.New("mocked error from transport"),
			},
		},
		checks: checks(
			hasErrorMsg("cannot do HTTP request: Post : mocked error from transport"),
		),
	}, {
		name: "error when can't read body",
		httpClient: &http.Client{
			Transport: &transportMock{
				roundTripResp: &http.Response{
					StatusCode: 200,
					Body:       nopCloserError{},
				},
			},
		},
		checks: checks(
			hasErrorMsg("can't read response body: mocked error from nop closer"),
		),
	}, {
		name: "error when got status code other than OK",
		httpClient: &http.Client{
			Transport: &transportMock{
				roundTripResp: &http.Response{
					StatusCode: 418,
					Body:       ioutil.NopCloser(bytes.NewBufferString("ok")),
				},
			},
		},
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("got status code 418: akismet API returned non 200 status code"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &akismetClient{
				blogUrl:    "http://some-blog.com",
				httpClient: tt.httpClient,
			}
			_, err := cli.post(context.Background(), tt.url, &url.Values{})
			for _, ch := range tt.checks {
				ch(err, t)
			}
		})
	}
}
