package main

import (
	"blog-me/internal/assets"
	"blog-me/internal/blog"
	"net/http"
)

var br = blog.BlogRoute
var asets = assets.AssetsRoute

func main() {
	http.Handle("/blog/", br)
	http.Handle("/assets/", asets)
	http.ListenAndServe(":8080", nil)
}
