package router

import (
	"fmt"
	"net/http"
	"regexp"
)

type RegexpRouter struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

func NewRegexpRouter() *RegexpRouter {
	return &RegexpRouter{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

func (r *RegexpRouter) Add(regex string, callback http.HandlerFunc) {
	r.handlers[regex] = callback
	r.cache[regex] = regexp.MustCompile(regex)

}

func (r *RegexpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	check := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	for pattern, handler := range r.handlers {
		if r.cache[pattern].Match([]byte(check)) {
			handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}
