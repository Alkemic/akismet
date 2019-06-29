package akismet

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

func TestNewAkismetCheck(t *testing.T) {
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
		comment            *Comment
		responseBody       string
		responseStatusCode int
		checks             []check
	}{{
		name: "success with 'true' as a response",
		comment: &Comment{
			CommentType: "comment",
		},
		responseBody:       "true",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(true),
			hasPayload("blog=&comment_type=comment"),
		),
	}, {
		name: "success with 'false' as a response",
		comment: &Comment{
			CommentType: "comment",
		},
		responseBody:       "false",
		responseStatusCode: 200,
		checks: checks(
			hasNoError,
			hasResult(false),
			hasPayload("blog=&comment_type=comment"),
		),
	}, {
		name: "error when status code is not OK",
		comment: &Comment{
			CommentType: "comment",
		},
		responseStatusCode: 418,
		checks: checks(
			hasCauseError(ErrNonOKStatusCode),
			hasErrorMsg("error during comment check request: got status code 418: akismet API returned non 200 status code"),
			hasPayload("blog=&comment_type=comment"),
		),
	}, {
		name: "error when got unusual response with 'false' as a response",
		comment: &Comment{
			CommentType: "comment",
		},
		responseBody:       "Missing required field: user_ip.",
		responseStatusCode: 200,
		checks: checks(
			hasCauseError(ErrUnusualResponse),
			hasErrorMsg("got response: 'Missing required field: user_ip.': got unusual response"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/asd/comment-check" {
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
				httpClient: &http.Client{},
				akismetUrl: ts.URL + "/%s/%s",
				key:        "asd",
			}
			result, err := cli.Check(context.Background(), tt.comment)
			for _, ch := range tt.checks {
				ch(result, err, buffer, t)
			}
		})
	}
}
