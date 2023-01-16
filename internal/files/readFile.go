package files

import (
	"fmt"
	"strings"

	db "blog-me/internal/database"

	"golang.org/x/net/html"
)

func PostReader(file []byte, path string) (p *db.Post, t []string) {
	doc, err := html.Parse(strings.NewReader(string(file)))
	if err != nil {
		fmt.Println("error ===> ", err)
	}

	po := db.Post{}
	po.File = path

	// Ill get all the body of the file and proccess possition
	content := GetElementByTag(doc, "body")
	cont, err := RenderElement(content)
	if err != nil {
		fmt.Println("error rendering content")
		return
	}

	// Get Title
	title := GetElementByClass(content, "title")
	//tags := GetElementByClass(content, "tag")
	tags := tagsExtractor(content)
	created, updated := GetCreation(content)
	po.Created = created
	po.Updated = updated

	// conver body in div  with a class
	contentStr := strings.Replace(cont, "body>", `div class="internal root-post">`, 1)
	po.Content = strings.Replace(contentStr, "</body>", `</div>`, 1)
	if title != nil {
		po.Title = title.FirstChild.Data
	} else {
		fmt.Println("Error title class not found")
		return
	}

	return &po, tags
}
