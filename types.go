package akismet

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Comment struct represents all information that will be send to endpoint.
type Comment struct {
	UserIP        string
	UserAgent     string
	Referrer      string
	Permalink     string
	Type          string
	Author        string
	AuthorEmail   string
	AuthorURL     string
	Content       string
	Language      string
	Charset       string
	UserRole      string
	Created       string
	Modified      string
	IsTest        string
	RecheckReason string
}

func (c *Comment) toValues() *url.Values {
	p := &url.Values{}
	p.Add("user_ip", c.UserIP)
	p.Add("user_agent", c.UserAgent)
	if c.Referrer != "" {
		p.Add("referrer", c.Referrer)
	}
	if c.Permalink != "" {
		p.Add("permalink", c.Permalink)
	}
	if c.Type != "" {
		p.Add("comment_type", c.Type)
	}
	if c.Author != "" {
		p.Add("comment_author", c.Author)
	}
	if c.AuthorEmail != "" {
		p.Add("comment_author_email", c.AuthorEmail)
	}
	if c.AuthorURL != "" {
		p.Add("comment_author_url", c.AuthorURL)
	}
	if c.Content != "" {
		p.Add("comment_content", c.Content)
	}
	if c.Language != "" {
		p.Add("blog_lang", c.Language)
	}
	if c.Charset != "" {
		p.Add("blog_charset", c.Charset)
	}
	if c.UserRole != "" {
		p.Add("user_role", c.UserRole)
	}
	if c.Created != "" {
		p.Add("comment_date_gmt", c.Created)
	}
	if c.Modified != "" {
		p.Add("comment_post_modified_gmt", c.Modified)
	}
	if c.IsTest != "" {
		p.Add("is_test", c.IsTest)
	}
	if c.RecheckReason != "" {
		p.Add("recheck_reason", c.RecheckReason)
	}
	return p
}

// Validate checks if user ip and user agent are present, and if present validates create/update dates.
func (c *Comment) Validate() error {
	if c.UserIP == "" {
		return errors.New("field user ip is required")
	}
	if c.UserAgent == "" {
		return errors.New("field user agent is required")
	}
	if c.Created != "" {
		_, err := time.Parse(time.RFC3339, c.Created)
		if err != nil {
			return errors.Wrap(err, "cannot parse created date")
		}
	}
	if c.Modified != "" {
		_, err := time.Parse(time.RFC3339, c.Modified)
		if err != nil {
			return errors.Wrap(err, "cannot parse modified date")
		}
	}
	return nil
}
