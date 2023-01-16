package assets

import (
	"blog-me/internal/router"
	"os"

	"net/http"
)

var (
	AssetsRoute *router.RegexpRouter = router.NewRegexpRouter()
	assets      string               = os.Getenv("ASSETSDIR")
)

func init() {
	AssetsRoute.Add("(GET|HEAD) /assets/(css|images|fonts|js)/[-_a-zA-Z0-9]{1,40}\\.(css|img|png|jpg|jpeg|js|svg)$", assetsServer)
}

// Route for serving images
func assetsServer(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(assets))
	http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
}
