// Here reside the files behavior to abstract from the handler
package files

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	filesPath string = os.Getenv("FILESDIR")
)

// / Sorting by date structure
type byDate []fs.DirEntry

func (b byDate) Less(i, j int) bool {
	x, _ := b[i].Info()
	y, _ := b[j].Info()

	return x.ModTime().Unix() > y.ModTime().Unix()
}
func (b byDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byDate) Len() int      { return len(b) }

//End sort by date

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

type Post struct {
	Name         string
	Date         string
	Title        template.HTML
	ContentTable template.HTML
	Content      template.HTML
	Label        template.HTML
	Postamble    template.HTML
}

func init() {
	if filesPath == "" {
		log.Fatal("FILESDIR env variable not set")
	}
}

func ReadFile(title string) (p Post, err error) {
	fileLocation := fmt.Sprintf("%s/%s.html", filesPath, strings.ToLower(title))
	info, err := os.Stat(fileLocation)
	if err != nil {
		fmt.Println("err on info ", info)
		fmt.Printf("%s/%s.org", filesPath, title)
		return
	}
	file, err := os.ReadFile(filesPath + "/" + strings.ToLower(title) + ".html")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error aqui " + title)
		fmt.Printf("%s/%s.html", filesPath, title)

		return Post{}, err
	}
	p = ContentCreator(file)
	caser := cases.Title(language.Spanish)
	p.Name = caser.String(title)
	p.Date = formatDate(float64(time.Now().UnixMilli() - info.ModTime().UnixMilli()))
	//p.Content = template.HTML(ContentCreator(file))
	//p.Label = ""
	return p, nil
}

func ReadDir(query string) ([]fs.DirEntry, error) {
	dirEntry, err := os.ReadDir(filesPath)
	if err != nil {
		log.Fatal("Cant read directory")
	}
	if query == "" || query == "date" {
		sort.Sort(byDate(dirEntry))
	} else if query == "name" {
		sort.Sort(byName(dirEntry))
	}
	return dirEntry, nil

}

func OrgFitered(dirEntry []fs.DirEntry) ([]Post, error) {
	if len(dirEntry) == 0 {
		return nil, fmt.Errorf("empty list nothing to show")
	}
	posts := []Post{}
	regexper := regexp.MustCompile(`.*\.html$`)
	caser := cases.Title(language.Spanish)
	for _, item := range dirEntry {
		if regexper.Match([]byte(item.Name())) {
			date, _ := item.Info()
			daysSince := formatDate(float64(time.Now().UnixMilli() - date.ModTime().UnixMilli()))
			title := strings.Replace(item.Name(), ".html", "", 1)
			posts = append(posts, Post{Name: caser.String(title), Date: daysSince, Content: "", Label: ""})

		}
	}
	return posts, nil
}

func formatDate(since float64) string {
	sinceHours := since / 3_600_000
	value := ""
	switch {
	case sinceHours < 24 && since > 0:
		value = fmt.Sprintf("%0.1f Hours Ago", sinceHours)
	case sinceHours > 24:
		value = fmt.Sprintf("%0.1f Days Ago", sinceHours/24)
	default:
		value = "Unknown Date"
	}
	return value
}

// ContentCreator will parse html file extract the required content from the generated html
// Will get information as title date author
func ContentCreator(file []byte) (p Post) {
	doc, err := html.Parse(strings.NewReader(string(file)))
	if err != nil {
		fmt.Println("error ===> ", err)
	}

	cont := bytes.Buffer{}
	content := GetElementByID(doc, "content")
	err = html.Render(&cont, content)
	if err != nil {
		fmt.Println("Error =====> ", err)
	}
	p.Content = template.HTML(cont.String())

	label := &bytes.Buffer{}
	tags := GetElementByClass(content, "tag")
	err = html.Render(label, tags)
	if err != nil {
		fmt.Println("Error =====> ", err)
	}
	p.Label = template.HTML(label.String())
	/*

		tit := &bytes.Buffer{}
		title := GetElementByClass(doc, "title")
		err = html.Render(tit, title)
		if err != nil {
			fmt.Println("Error =====> ", err)
		}
		p.Title = template.HTML(tit.String())
	*/
	tableContents := GetElementByID(doc, "table-of-contents")
	tableCont := &bytes.Buffer{}
	err = html.Render(tableCont, tableContents)
	if err != nil {
		fmt.Println("Error =====> ", err)
	}
	p.ContentTable = template.HTML(tableCont.String())

	postamble := GetElementByID(doc, "postamble")
	foot := &bytes.Buffer{}
	err = html.Render(foot, postamble)
	if err != nil {
		fmt.Println("Error =====> ", err)
	}
	p.Postamble = template.HTML(foot.String())

	return p
}
