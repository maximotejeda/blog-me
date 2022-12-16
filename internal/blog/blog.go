package blog

import (
	"blog-me/internal/router"
	"net/http"
)

var BlogRoute = router.NewRegexpRouter()

func init() {
	BlogRoute.Add("(GET|HEAD) /blog/", baseRoute)

}

// base route wher the index of the blog will be placed ordered for more recent
func baseRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Blog principal page"))
}
