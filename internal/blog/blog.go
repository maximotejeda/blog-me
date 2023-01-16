package blog

import (
	"blog-me/internal/database"
	db "blog-me/internal/database"
	_ "blog-me/internal/files"
	"blog-me/internal/router"
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed templates/*
var templates embed.FS

type page struct {
	Where      string
	Post       db.Post
	Posts      []db.Post
	Pages      int
	ActualPage int
}

var BlogRoute = router.NewRegexpRouter()
var t *template.Template

func init() {
	BlogRoute.Add("(GET|HEAD) /", baseRoute)
	t, _ = template.New("").Funcs(template.FuncMap{
		"replace": strings.Replace,
		"split":   strings.Split,
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"unSafeHTML": html.UnescapeString,
		"formatDate": formatDate,
		"incPage":    incPage,
		"decPage":    decPage,
	}).ParseFS(templates,
		"templates/*.html")

}

// base route wher the index of the blog will be placed ordered for more recent
// will dispatch the Template for the / {root} page of the blog
func baseRoute(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	params := r.URL.Query()
	noSTR := params.Get("page")
	title := params.Get("title")
	tag := params.Get("label")
	var no int
	var err error
	if noSTR == "" {
		no = 1
	} else {
		no, err = strconv.Atoi(noSTR)
		if err != nil {
			no = 1
		}

	}
	if title != "" {
		blogPostTitle(w, r)
		return
	}
	if tag != "" {
		blogPostLabel(w, r)
		return
	}

	page := page{Where: "base", Post: database.Post{}, Posts: nil, Pages: 0.0, ActualPage: 0}

	page.Posts, _ = db.Get(nil, no)
	cantPost := db.GetPages()
	pages := 0
	if cantPost > 0 {
		pages = int(math.Ceil(float64(cantPost) / 10))
	}

	page.ActualPage = no
	page.Pages = pages
	err = t.ExecuteTemplate(&b, "index.html", page)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	b.WriteTo(w)

}

// Find post by title
func blogPostTitle(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	//path := r.URL.Path
	params := r.URL.Query()
	title := params.Get("title")
	//title := strings.Replace(path, "/title/", "", 1)
	page := page{Where: title, Post: db.Post{}, Posts: nil}
	if title == "" {
		http.NotFound(w, r)
		return
	}
	posts, _ := db.Get(&db.Post{Title: title}, 1)
	page.Post = posts[0]

	err := t.ExecuteTemplate(&b, "index.html", page)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	b.WriteTo(w)
}

// Find post by Label
func blogPostLabel(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	params := r.URL.Query()
	noSTR := params.Get("page")
	label := params.Get("label")
	var (
		no  int
		err error
	)

	if noSTR == "" {
		no = 1
	} else {
		no, err = strconv.Atoi(noSTR)
		if err != nil {
			no = 1
		}

	}
	if label == "" {
		http.NotFound(w, r)
	}
	posts, _ := db.Get(&db.Tag{Name: label}, no)
	page := page{Where: "base", Post: database.Post{}, Posts: nil}
	page.Posts = posts
	cantPost := db.GetPages()
	pages := 0
	if cantPost > 0 {
		pages = int(math.Ceil(float64(cantPost) / 10))

	}
	page.ActualPage = no
	page.Pages = pages

	err = t.ExecuteTemplate(&b, "index.html", page)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	b.WriteTo(w)
}

func formatDate(t time.Time) string {
	formTime := time.Now().Unix() - t.Unix()
	switch {
	case formTime > 0 && formTime < 60:
		return fmt.Sprintf("%d Seconds Ago.", formTime)
	case formTime > 60 && formTime < 3600:
		return fmt.Sprintf("%d Minutes Ago.", formTime/60)
	case formTime > 3600 && formTime < 86400:
		return fmt.Sprintf("%d Hours Ago.", formTime/3600)
	case formTime > 86400 && formTime < 2_542_000:
		return fmt.Sprintf("%d Days Ago.", formTime/86400)
	case formTime > 2_542_000:
		return fmt.Sprintf("%d Mounth Ago", formTime/2542000)
	default:
		return "Unknown"
	}
}

func incPage(n int) int {
	return n + 1
}
func decPage(n int) int {
	return n - 1
}
