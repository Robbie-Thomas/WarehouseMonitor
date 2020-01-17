package modeltests

import (
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("Error seeding user and post  table %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}
