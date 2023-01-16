package database

import (
	"fmt"
	"strings"
)

// Remember use pointer to avoid pass value by copy
// return an integer that is the id of the actual row in the table
func Add(val interface{}) int {
	switch val.(type) {
	case *Post:
		table := "post"
		post := val.(*Post)
		if post.Title == "" {
			fmt.Println("Title is empty nothing can be found")
			return 0
		}
		stmt := fmt.Sprintf("INSERT INTO %s(title, content, file, created, updated, deleted) VALUES('%s','%s', '%s','%v', '%v',NULL)", table, strings.ToLower(post.Title), post.Content, post.File, post.Created.Format("2006-01-02T15:04:05.999Z"), post.Updated.Format("2006-01-02T15:04:05.999Z"))
		res, err := DB.Exec(stmt)
		if err != nil {
			fmt.Println("add =>>", err)
			return 0
		}
		id, _ := res.LastInsertId()
		return int(id)
	case *Tag:
		table := "tag"
		tag := val.(*Tag)
		if tag.Name == "" {
			fmt.Println("Tag name cant be empty when searching by tag")
			return 0
		}
		stmt := fmt.Sprintf("INSERT INTO %s(name) VALUES('%s')", table, tag.Name)
		table1 := "tag_rel"
		res, err := DB.Exec(stmt)
		// we try to add the tag to table
		if err != nil {
			//  check if is a constrint err then get the actual id of the tag
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				fmt.Println("unique contraint Tag alredy exists", tag.Name)
				stmt2 := fmt.Sprintf("select id from tag where name = '%s'", tag.Name)
				res := DB.QueryRow(stmt2)
				err = res.Scan(&tag.TagID)
				if err != nil {
					fmt.Println("add tag get ==>", err)
				}
			} else {
				fmt.Println("add tag name =>>", err.Error())
				return 0
			}
		} else {
			id, _ := res.LastInsertId()
			tag.TagID = int(id)
		}
		// ensure not exist the same post with the same tag

		stmt1 := fmt.Sprintf("INSERT OR IGNORE INTO %s(post_id, tag_id) VALUES(%d, %d)", table1, tag.PostID, tag.TagID)
		// brig the newly added tag
		_, err = DB.Exec(stmt1)
		if err != nil {
			fmt.Println(err)
		}
		return tag.TagID

	default:
		fmt.Println("Add type not recognized")
	}
	return 0
}
