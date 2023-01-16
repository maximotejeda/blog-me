package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestGetPost(t *testing.T) {
	p := &Post{Title: "hello"}
	res, err := Get(p, 1)
	assert.Empty(t, err)
	assert.Equal(t, "hello", res[0].Title)
	assert.Equal(t, 2, res[0].ID)

}
func TestGetTag(t *testing.T) {
	p := &Tag{Name: "golang"}
	res, err := Get(p, 1)
	assert.Empty(t, err)
	assert.GreaterOrEqual(t, 2, len(res))
}
func TestGetPosts(t *testing.T) {
	res, err := Get(nil, 1)
	assert.Empty(t, err)
	assert.GreaterOrEqual(t, 2, len(res))
}

func TestFailForeignKey(t *testing.T) {
	//Trying to insert unexistent foreign keys
	_, err := DB.Exec("insert into tag_rel(post_id, tag_id) values(200, 1000)")
	assert.Error(t, err)
}
