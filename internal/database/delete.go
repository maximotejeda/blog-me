package database

import (
	"fmt"
)

func Delete(val interface{}) {
	switch val.(type) {
	case *Post:
		post := val.(*Post)
		if post.File == "" {
			fmt.Println("File is empty nothing can be found")
			return
		}
		stmt := fmt.Sprintf(`DELETE FROM post 
					WHERE file = '%s' AND deleted IS NULL`, post.File)
		_, err := DB.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
	case *Tag:
		tag := val.(*Tag)
		if tag.Name == "" {
			fmt.Println("Tag name cant be empty when searching by tag")
			return
		}
		stmt := fmt.Sprintf(`DELETE FROM tag_rel
					WHERE tag_id = %d AND post_id = %d AND name = %s
		`, tag.TagID, tag.PostID, tag.Name)
		fmt.Println(stmt)
	default:
	}
}
