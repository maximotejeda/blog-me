package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbDIR  = os.Getenv("DBDIR")
	dbFile = os.Getenv("DBFILE")
	DB     *sql.DB
)

// Post Model to interact with db
// main source of get a post is by title
type Post struct {
	ID      int        `json:"id"`
	Title   string     `json:"title"`
	Content string     `json:"content"`
	File    string     `json:"file"`
	Tags    string     `json:"tags"`
	Created time.Time  `json:"created"`
	Updated time.Time  `json:"updated"`
	Deleted *time.Time `json:"deleted"`
}

type Tag struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	PostID int    `json:"post_id"`
	TagID  int    `json:"tag_id"`
}

func init() {
	createDB()
	createTables()
}

func createTables() {
	// create blog table
	_, err := DB.Exec("create table if not exists post(id integer primary key autoincrement,title varchar(250) unique not null, content varchar not null,file varchar unique not null, created datetime not null, updated datetime not null, deleted datetime)")
	if err != nil {
		log.Fatal("unable to create table in createTables blog: ", err)
	}
	// create tag table
	_, err = DB.Exec("create table if not exists tag(id integer primary key autoincrement, name varchar(100) unique not null)")
	if err != nil {
		log.Fatal("unable to create table in createTables tag")
	}
	//create relation between tables
	_, err = DB.Exec("create table if not exists tag_rel(post_id integer not null, tag_id integer not null, unique(post_id, tag_id), foreign key(post_id) references post(id) on delete cascade, foreign key(tag_id) references tag(id) on delete cascade)")
	if err != nil {
		log.Fatal("unable to create table in createTables tag_rel")
	}

	// this must failt to ensure foreign keys
	_, err = DB.Exec("insert into tag_rel(post_id, tag_id) values(1,1)")
	if err != nil {
		fmt.Println("keys working", err)
	}

	_, err = DB.Exec(`create view if not exists post_tag
							as select post.id, post.title, post.content, post.created, post.updated, post.deleted,
								group_concat(tag.name) as tags 
								from post
								join tag
								join tag_rel
								on tag_rel.post_id = post.id and tag_rel.tag_id = tag.id 
								group by post.id
	`)
	if err != nil {
		fmt.Println("create view: ", err)
	}

}
func createDB() {
	var err error
	// remove each db in case of reuse dunno if should refresh db n each reset dont think is a good idea
	/*err = os.RemoveAll(dbDIR)
	if err != nil {
		panic(err)
	}*/
	if dbDIR == "" {
		panic("error db folder env variable not set")
	}
	if dbFile == "" {
		panic("error db file env variable not set")
	}
	err = os.Mkdir(dbDIR, 0777)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("./db folder alredy exist")
		} else {
			panic(err)
		}
	}

	connString := fmt.Sprintf("%s/%s?cache=%s&mode=%s?_foreign_keys=on", dbDIR, dbFile, "shared", "memory")
	DB, err = sql.Open("sqlite3", connString)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		fmt.Println("===>err")
	}
	fmt.Println("Opening ", connString)
}
