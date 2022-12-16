package main

import (
	"blog-me/internal/blog"
	"net/http"
)

var br = blog.BlogRoute

func main() {
	http.Handle("/blog/", br)
	http.ListenAndServe(":8080", nil)
}
