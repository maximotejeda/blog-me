package blog

import (
	"blog-me/internal/router"
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	files string = os.Getenv("FILESDIR")
)

// / Sorting by date structure
type byDate []fs.DirEntry

type page struct {
	Where string
	Post  string
	Posts []string
}

func (b byDate) Less(i, j int) bool {
	x, _ := b[i].Info()
	y, _ := b[j].Info()

	return x.ModTime().Unix() > y.ModTime().Unix()
}
func (b byDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byDate) Len() int      { return len(b) }

//End of sorting by date

// sort by name Structure
type byName []fs.DirEntry

func (n byName) Less(i, j int) bool {
	x := strings.ToLower(n[i].Name())
	y := strings.ToLower(n[j].Name())

	return x[:1] < y[:1]
}
func (n byName) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n byName) Len() int      { return len(n) }

// End of sort by name

var BlogRoute = router.NewRegexpRouter()
var t *template.Template

func init() {
	if files == "" {
		log.Fatal("FILESDIR env variable not set")
	}
	BlogRoute.Add("(GET|HEAD) /$", baseRoute)
	BlogRoute.Add("(GET|HEAD) /title/", blogPostTitle)
	t = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/content.html",
		"templates/head.html",
		"templates/header.html",
		"templates/footer.html"))
}

// base route wher the index of the blog will be placed ordered for more recent
// will dispatch the Template for the / {root} page of the blog
func baseRoute(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	dirEntry, err := os.ReadDir(files)
	params := r.URL.Query()
	query := params.Get("query")
	page := page{Where: "base", Post: "", Posts: []string{}}
	if query == "" || query == "date" {
		sort.Sort(byDate(dirEntry))
	} else if query == "name" {
		sort.Sort(byName(dirEntry))
	}

	posts := []string{}
	regexper := regexp.MustCompile(`.*\.org$`)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range dirEntry {
		if regexper.Match([]byte(item.Name())) {

			posts = append(posts, strings.Replace(item.Name(), ".org", "", 1))

		}
	}
	page.Posts = posts
	err = t.ExecuteTemplate(&b, "index.html", page)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	b.WriteTo(w)

}

// Find post by title
func blogPostTitle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	title := strings.Replace(path, "/title/", "", 1)
	if title == "" {
		http.NotFound(w, r)
	}
	w.Write([]byte(title))
}

// Find post by Label
func blogPostLabel(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	title := strings.Replace(path, "/label/", "", 1)
	if title == "" {
		http.NotFound(w, r)
	}
}
