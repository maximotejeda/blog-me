/*Mantains the folder of files updated*/
package files

import (
	db "blog-me/internal/database"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	currFiles map[string]time.Time
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Check file structure without blocking
	FilesEquality()
	go concurrentCheck(wg)

}

func FilesEquality() {
	actualFilesList, err := os.ReadDir(filesPath)
	actualFiles := make(map[string]time.Time)
	wg := &sync.WaitGroup{}
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range actualFilesList {
		if strings.HasSuffix(item.Name(), ".html") { // only proccess html

			info, _ := item.Info()
			actualFiles[item.Name()] = info.ModTime()
		}
	}

	if currFiles == nil {
		currFiles = actualFiles
		//ADD files to db
		for file := range currFiles {
			wg.Add(1)
			go fileOpAdd(file, db.Add, wg)
		}
		wg.Wait()
		return
	} else {
		prevLen, actualLen := len(currFiles), len(actualFiles)
		if prevLen < actualLen {
			for item := range actualFiles {
				if date, ok := currFiles[item]; !ok {
					wg.Add(1)
					currFiles[item] = date
					//TODO add file to DB
					go fileOpAdd(item, db.Add, wg)
				}
			}
		} else if actualLen < prevLen {
			for item := range currFiles {
				if _, ok := actualFiles[item]; !ok {
					delete(currFiles, item)
					//TODO Delete file from DB
					wg.Add(1)
					go fileOpDel(item, db.Delete, wg)
				}
			}

		}
		for idx, item := range currFiles {
			prevInfo := item
			actualInfo := actualFiles[idx]
			if prevInfo != actualInfo {
				currFiles[idx] = actualFiles[idx]
				//TODO update file in database
				wg.Add(1)
				go fileOpEdit(idx, db.Update, wg)
			}
		}

	}
	wg.Wait()
}

func concurrentCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-time.After(time.Hour * 12):
			fmt.Println("\rchecking update")
			FilesEquality()

		}
	}
}
