package files

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type processHtml func(n *html.Node, s string) bool

// return an attribute for a node
func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

// return the first node with class
func GetClass(n *html.Node, s string) bool {
	if n.Type == html.ElementNode {
		si, ok := getAttribute(n, "class")
		if ok && si == s {
			return true
		}
	}
	return false
}

// test element ID
func GetID(n *html.Node, s string) bool {
	if n.Type == html.ElementNode {
		si, ok := getAttribute(n, "id")
		if ok && si == s {
			return true
		}
	}
	return false
}

// test element tag
func GetTag(n *html.Node, s string) bool {
	return n.Data == s
}

// Traverse Data
func TraverseNode(n *html.Node, s string, p processHtml) *html.Node {
	if p(n, s) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := TraverseNode(c, s, p)
		if result != nil {
			return result
		}
	}
	return nil
}

// Traver by id the node bring the first item with id
func GetElementByID(n *html.Node, id string) *html.Node {
	return TraverseNode(n, id, GetID)
}

// traverse the node By class return the first item with that id
func GetElementByClass(n *html.Node, class string) *html.Node {
	return TraverseNode(n, class, GetClass)
}

// bring the first item with the html tag
func GetElementByTag(n *html.Node, tag string) *html.Node {
	return TraverseNode(n, tag, GetTag)
}

// convert node elements to string
func RenderElement(n *html.Node) (string, error) {
	s := &bytes.Buffer{}
	err := html.Render(s, n)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return s.String(), nil
}

func tagsExtractor(n *html.Node) []string {
	q := goquery.NewDocumentFromNode(n)
	tags := []string{}
	tager := make(map[string]struct{})
	q.Find(".tag").Each(func(i int, s *goquery.Selection) {
		for i := range s.Nodes {
			//r := strings.Split(s.Eq(i).Text(), " ")
			res := s.Eq(i).Text()
			nbs := '\u00a0'
			results := strings.Map(func(r rune) rune {
				if r == nbs {
					return '\u0020'
				}
				return r
			}, res)
			values := strings.Split(results, " ")
			for _, i := range values {
				tager[i] = struct{}{}
			}

		}

	})
	for i := range tager {
		tags = append(tags, i)
	}
	if len(tags) == 0 {
		tags = append(tags, "untaged")
	}
	return tags

}

// Get creation from file
func GetCreation(n *html.Node) (created, updated time.Time) {
	q := goquery.NewDocumentFromNode(n)
	dates := make(map[string]time.Time)
	q.Find(".date").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		strList := strings.Split(str, ":")
		if len(strList) > 1 {
			date, err := time.Parse("2006-01-02", strings.Replace(strList[1][0:11], " ", "", 2))
			if err != nil {
				fmt.Println(err)
				return
			}
			dates[strList[0]] = date

		} else {
			dates["Created"] = time.Now()
			dates["Date"] = time.Now()
		}

	})
	// date is the stablished date by me updated is the time the file was parsed
	created = dates["Date"]
	updated = dates["Created"]
	return
}
