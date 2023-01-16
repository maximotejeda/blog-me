package files

import (
	db "blog-me/internal/database"
	"fmt"
	"os"
	"sync"
)

// File reader
func readFile(name string) []byte {
	filePath := fmt.Sprintf("%s/%s", filesPath, name)
	// Leemos el archivo
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// check file content more than 0
	if len(file) <= 1000 {
		fmt.Println("empty file", filePath)
		return nil
	}
	return file
}

// As add returns the id this is the signature of the add
func fileOpAdd(name string, op func(interface{}) int, wg *sync.WaitGroup) {
	defer wg.Done()
	filePath := fmt.Sprintf("%s/%s", filesPath, name)
	fmt.Println("Adding File: ", filePath)
	file := readFile(name)
	if file == nil {
		return
	}
	// File read and with content now lets add to db

	post, tags := PostReader(file, filePath)
	if post != nil {
		id := op(post)
		if len(tags) > 0 {
			for _, item := range tags {
				op(&db.Tag{Name: item, PostID: id})
			}

		} else {
			op(&db.Tag{Name: "untaged", PostID: id})
		}
	}
}

// for update and delete
func fileOpEdit(name string, op func(interface{}), wg *sync.WaitGroup) {
	defer wg.Done()
	filePath := fmt.Sprintf("%s/%s", filesPath, name)
	fmt.Println("Updating File: ", filePath)
	file := readFile(name)
	if file == nil {
		return
	}
	post, _ := PostReader(file, filePath)
	if post != nil {
		op(post)
	}
}

func fileOpDel(name string, op func(interface{}), wg *sync.WaitGroup) {
	defer wg.Done()
	filePath := fmt.Sprintf("%s/%s", filesPath, name)
	fmt.Println("Deleting File: ", filePath)
	op(&db.Post{File: filePath})
}
