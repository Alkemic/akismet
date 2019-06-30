package main

import (
	"context"
	"flag"
	"log"

	"akismet"
)

var (
	key     = flag.String("key", "", "Akismet API key")
	blogUrl = flag.String("blog-url", "http://some-blog.com", "Blog URI (inc. scheme)")
)

func main() {
	flag.Parse()

	client, err := akismet.NewAkismet(*key, *blogUrl)
	if err != nil {
		log.Fatalf("error creating client instance: %v", err)
	}

	isSpam, err := client.Check(context.Background(), &akismet.Comment{
		CommentType:   "comment",
		CommentAuthor: "viagra-test-123",
		UserIP:        "8.8.8.8",
		UserAgent:     "Mozilla/6.1.6",
	})
	if err != nil {
		log.Fatalf("got error: %v", err)
	}

	log.Printf("is spam: %t", isSpam)
}
