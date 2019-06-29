package akismet

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Comment struct represents all information that will be send to endpoint.
type Comment struct {
	UserIP                 string
	UserAgent              string
	Referrer               string
	Permalink              string
	CommentType            string
	CommentAuthor          string
	CommentAuthorEmail     string
	CommentAuthorURL       string
	CommentContent         string
	BlogLang               string
	BlogCharset            string
	UserRole               string
	CommentDateGMT         string
	CommentPostModifiedGMT string
	IsTest                 string
	RecheckReason          string
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
	if c.CommentType != "" {
		p.Add("comment_type", c.CommentType)
	}
	if c.CommentAuthor != "" {
		p.Add("comment_author", c.CommentAuthor)
	}
	if c.CommentAuthorEmail != "" {
		p.Add("comment_author_email", c.CommentAuthorEmail)
	}
	if c.CommentAuthorURL != "" {
		p.Add("comment_author_url", c.CommentAuthorURL)
	}
	if c.CommentContent != "" {
		p.Add("comment_content", c.CommentContent)
	}
	if c.BlogLang != "" {
		p.Add("blog_lang", c.BlogLang)
	}
	if c.BlogCharset != "" {
		p.Add("blog_charset", c.BlogCharset)
	}
	if c.UserRole != "" {
		p.Add("user_role", c.UserRole)
	}
	if c.CommentDateGMT != "" {
		p.Add("comment_date_gmt", c.CommentDateGMT)
	}
	if c.CommentPostModifiedGMT != "" {
		p.Add("comment_post_modified_gmt", c.CommentPostModifiedGMT)
	}
	if c.IsTest != "" {
		p.Add("is_test", c.IsTest)
	}
	if c.RecheckReason != "" {
		p.Add("recheck_reason", c.RecheckReason)
	}
	return p
}

func (c *Comment) Validate() error {
	if c.UserIP == "" {
		return errors.New("field user ip is required")
	}
	if c.UserAgent == "" {
		return errors.New("field user agent is required")
	}
	if c.CommentDateGMT != "" {
		_, err := time.Parse(time.RFC3339, c.CommentDateGMT)
		if err != nil {
			return errors.Wrap(err, "cannot parse created date")
		}
	}
	if c.CommentPostModifiedGMT != "" {
		_, err := time.Parse(time.RFC3339, c.CommentPostModifiedGMT)
		if err != nil {
			return errors.Wrap(err, "cannot parse modified date")
		}
	}
	return nil
}
