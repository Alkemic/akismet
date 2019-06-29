# akismet
A GO Akismet client, made for easy use and testing

# Installation

```$ go get github.com/Alkemic/akismet```

## Usage

First create an instance using ctor with :

```go
akismetClient := akismet.NewClient("akismet-key")
ctx := context.Background()
akismetClient.Check(ctx, &akismet.Comment{
    CommentType:   "comment",
    Blog:          "http://some-blog.com",
    CommentAuthor: "author",
    UserIP:        "8.8.8.8",
})
```
