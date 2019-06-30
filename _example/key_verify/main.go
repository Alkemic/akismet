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

	verified, err := client.Verify(context.Background())
	if err != nil {
		log.Fatalf("got error: %v", err)
	}

	log.Printf("key verification: %t", verified)
}
