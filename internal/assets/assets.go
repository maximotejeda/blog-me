package assets

import (
	"blog-me/internal/router"
	"fmt"
	"net/http"
)

var (
	AssetsRoute = router.NewRegexpRouter()
	//assetsRoute = os.Getenv("ASSESTSDIR")
)

func init() {
	AssetsRoute.Add("(GET|HEAD) /assets/(css|images|fonts)/[-_a-zA-Z0-9]{1,40}\\.(css|img|png|jpg|jpeg)$", assetsServer)
}

// Route for serving images
func assetsServer(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./posts/assets/"))
	fmt.Println(r.URL.Path)
	http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
}
