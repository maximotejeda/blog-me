package database

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAddPost(t *testing.T) {
	p := &Post{Title: "test post", Content: "Test content", Created: time.Now(), Updated: time.Now()}
	Add(p)
	pl, err := DB.Exec(fmt.Sprintf("SELECT * FROM post WHERE title = '%s'", p.Title))
	assert.NoError(t, err)
	num, _ := pl.RowsAffected()
	assert.Equal(t, 1, int(num))
}

func TestAddTag(t *testing.T) {
	p, err := Get(&Post{Title: "hello"}, 1)
	assert.NoError(t, err)
	fmt.Println(p)
	assert.Equal(t, "world", p[0].Content)
}

func TestAddTagRel(t *testing.T) {}
