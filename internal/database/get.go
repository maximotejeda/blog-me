package database

import (
	"fmt"
	"log"
	"strings"
)

func Get(val interface{}, page int) (result []Post, err error) {
	table := "post_tag" // view created with the joins concatening the tags name
	switch val.(type) {
	case *Post: // if the interface is a POST
		post := val.(*Post)
		if post.Title == "" {
			fmt.Println("Title is empty nothing can be found")
			return nil, fmt.Errorf("nothing to query empty title")
		}
		stmtSTR := fmt.Sprintf("SELECT id, title, tags, content, created, updated FROM %s WHERE title = '%s' AND deleted IS NULL", table, strings.ToLower(post.Title))
		// execute statement
		stmt, err := DB.Prepare(stmtSTR)
		if err != nil {
			log.Println("Get ===>>>", err)
		}
		res := stmt.QueryRow()
		p := Post{}
		err = res.Scan(&p.ID, &p.Title, &p.Tags, &p.Content, &p.Created, &p.Updated)
		if err != nil {
			fmt.Println("scan: ", err)
		}
		result = append(result, p)
		return result, nil
	case *Tag: // if the interface is a Tag
		tag := val.(*Tag)
		if tag.Name == "" {
			fmt.Println("Tag name cant be empty when searching by tag")
			return nil, fmt.Errorf("nothing to query empty tag")
		}
		// here i return a list of posts
		stmt := fmt.Sprintf(`SELECT id, title, content, tags, created, updated FROM %[2]s WHERE tags like '%[1]s%[3]s%[1]s' ORDER BY updated DESC limit 10 OFFSET %[4]d`, "%", table, tag.Name, (page*10)-10)
		res, err := DB.Query(stmt)
		if err != nil {
			fmt.Println("Get tag: ", err)
			return nil, fmt.Errorf("error tag: %v", err)
		}
		defer res.Close()
		for res.Next() {
			p := Post{}
			res.Scan(&p.ID, &p.Title, &p.Content, &p.Tags, &p.Created, &p.Updated)
			result = append(result, p)
		}
		return result, nil
	default:
		stmt := fmt.Sprintf("SELECT id, title, content, tags, created, updated FROM %s ORDER BY updated DESC LIMIT 10 OFFSET %d", table, (page*10)-10)
		res, err := DB.Query(stmt)
		if err != nil {
			fmt.Println("Get tag: ", err)
			return nil, fmt.Errorf("error tag: %v", err)
		}
		defer res.Close()
		for res.Next() {
			p := Post{}
			err = res.Scan(&p.ID, &p.Title, &p.Content, &p.Tags, &p.Created, &p.Updated)
			if err != nil {
				fmt.Println("Get posts: ", err)
				return nil, fmt.Errorf("error posts: %v", err)
			}
			result = append(result, p)
		}
		return result, nil
	}

}
func GetPages() int {
	var result int
	stmtSTR := "SELECT COUNT(*) FROM post_tag"
	stmt, _ := DB.Prepare(stmtSTR)
	res := stmt.QueryRow()

	err := res.Scan(&result)
	if err != nil {
		fmt.Println("ERR quering scaning pages====> ", err)
	}
	return result
}
