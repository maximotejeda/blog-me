package main

import (
	"blog-me/internal/assets"
	"blog-me/internal/blog"
	"log"
	"net/http"
	"os"
)

var (
	br    = blog.BlogRoute
	asets = assets.AssetsRoute
	port  = os.Getenv("SERVERPORT")
	addr  = os.Getenv("SERVERADDR")
)

func main() {
	http.Handle("/assets/", asets)
	http.Handle("/", br)
	log.Fatal(http.ListenAndServe(addr+":"+port, nil))
}
