# akismet
A GO Akismet client, made for easy use and testing

# Installation

```$ go get github.com/Alkemic/akismet```

## Usage

Validate you key:

```go
akismetClient, _ := akismet.NewClient("akismet-key", "http://some-blog.com")
ctx := context.Background()
validated, err := akismetClient.Valid(ctx)

if err != nil {
	// handle error
}
```

Check if comment is a SPAM:

```go
akismetClient, _ := akismet.NewClient("akismet-key", "http://some-blog.com")
ctx := context.Background()
isSpam, err := akismetClient.Check(ctx, &akismet.Comment{
    CommentType:   "comment",
    Blog:          "http://some-blog.com",
    CommentAuthor: "author",
    UserIP:        "8.8.8.8",
})

if err != nil {
	// handle error
}
```

Submit SPAM:

```go
akismetClient, _ := akismet.NewClient("akismet-key", "http://some-blog.com")
ctx := context.Background()
err := akismetClient.SubmitSpam(ctx, &akismet.Comment{
    CommentType:   "comment",
    Blog:          "http://some-blog.com",
    CommentAuthor: "author",
    UserIP:        "8.8.8.8",
})

if err != nil {
	// handle error
}
```

Submit HAM (aka false positive):

```go
akismetClient, _ := akismet.NewClient("akismet-key", "http://some-blog.com")
ctx := context.Background()
err := akismetClient.SubmitHam(ctx, &akismet.Comment{
    CommentType:   "comment",
    Blog:          "http://some-blog.com",
    CommentAuthor: "author",
    UserIP:        "8.8.8.8",
})

if err != nil {
	// handle error
}
```

## Advanced usage

You can use your own `http.Client` instance with calls to API, just use `WithHttpClient` functional option:
```go
customHttpClient := &http.Client{
	// your options here
}
akismet.NewClient("akismet-key", "http://some-blog.com", WithHttpClient(customHttpClient))
```
