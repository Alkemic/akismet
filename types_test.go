package akismet

import (
	"testing"
	"time"
)

func TestCommentValidate(t *testing.T) {
	type check func(err error, t *testing.T)
	checks := func(cs ...check) []check { return cs }

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

	validComment := &Comment{
		UserIP:                 "8.8.8.8",
		UserAgent:              "Mozilla/6.1.6",
		CommentDateGMT:         time.Now().Format(time.RFC3339),
		CommentPostModifiedGMT: time.Now().Format(time.RFC3339),
	}

	tests := []struct {
		name    string
		comment *Comment
		checks  []check
	}{{
		name:    "error on empty comment",
		comment: &Comment{},
		checks: checks(
			hasErrorMsg("field user ip is required"),
		),
	}, {
		name: "error on missing user agent",
		comment: &Comment{
			UserIP: "8.8.8.8",
		},
		checks: checks(
			hasErrorMsg("field user agent is required"),
		),
	}, {
		name: "error on invalid create date",
		comment: &Comment{
			UserIP:         "8.8.8.8",
			UserAgent:      "Mozilla/6.1.6",
			CommentDateGMT: "asdad",
		},
		checks: checks(
			hasErrorMsg(`cannot parse created date: parsing time "asdad" as "2006-01-02T15:04:05Z07:00": cannot parse "asdad" as "2006"`),
		),
	}, {
		name: "error on invalid modified date",
		comment: &Comment{
			UserIP:                 "8.8.8.8",
			UserAgent:              "Mozilla/6.1.6",
			CommentPostModifiedGMT: "asdad",
		},
		checks: checks(
			hasErrorMsg(`cannot parse modified date: parsing time "asdad" as "2006-01-02T15:04:05Z07:00": cannot parse "asdad" as "2006"`),
		),
	}, {
		name:    "success on valid comment",
		comment: validComment,
		checks: checks(
			hasNoError,
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			for _, ch := range tt.checks {
				ch(err, t)
			}
		})
	}
}
