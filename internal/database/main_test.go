package database

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	/*DB, err = sql.Open("sqlite3", "./test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("could not open db: %v", err)
	}*/
	DB.Exec("PRAGMA foreign_keys = ON")
	createTables()
	_, err = DB.Exec("insert into post(title, content, file, created, updated, deleted) values('hello', 'world', 'test' ,current_timestamp, current_timestamp, null)")
	_, err = DB.Exec("insert into post(title, content, file, created, updated, deleted) values('second', 'post', 'tester' ,current_timestamp, current_timestamp, null)")
	if err != nil {
		fmt.Println("blog: ", err)
	}
	// insert a tag
	_, err = DB.Exec("insert into tag(name) values('golang'),('javascript'),('react'),('git')")
	if err != nil {

		fmt.Println("insert test tags: ", err)
	}
	// insert relation
	_, err = DB.Exec("insert into tag_rel(post_id, tag_id) values(2,1),(2,2),(3,1),(3,2)")
	if err != nil {
		fmt.Println("insert test tags_rel: ", err)
	}
	defer func() {
		os.RemoveAll(dbDIR)
		DB.Close()
	}()
	return m.Run(), nil

}
