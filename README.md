# akismet
A GO Akismet client, made for easy use and testing

# Installation

```$ go get github.com/Alkemic/akismet```

## Usage

Check if comment is a SPAM:

```go
akismetClient := akismet.NewClient("akismet-key", "http://some-blog.com")
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
