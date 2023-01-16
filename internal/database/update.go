package database

import (
	"fmt"
)

func Update(val interface{}) {
	switch val.(type) {
	// in the case of post only content and tag can be modified
	// updated will be modified
	case *Post:
		post := val.(*Post)
		if post.Title == "" {
			fmt.Println("Title is empty nothing can be found")
			return
		}
		// first get the post
		postL, _ := Get(&Post{Title: post.Title}, 1)
		if len(postL) < 1 {
			fmt.Println("no post found")
			return
		}
		post1 := postL[0]
		if post.Content == "" || post.Content == post1.Content {
			fmt.Println("nothing to update")
			return
		}
		fmt.Println(post.Tags, post1.Tags)
		stmt := fmt.Sprintf(`UPDATE post
				SET content = '%s',
				updated = current_timestamp
				WHERE file = '%s' AND deleted IS NULL`, post.Content, post.File)
		_, err := DB.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}

	case *Tag:
	// TODO not implemented
	default:
	}
}
