package state

import (
	"forum/src/models"
	"sort"
)

func (r *State) SetPosts(posts models.PostsType) *State {
	r.Posts = posts
	sort.Slice(r.Posts, func(i, j int) bool {
		return r.Posts[i].TimestampString > r.Posts[j].TimestampString
	})
	return r
}

func (r *State) EditPosts() *models.PostsType {
	return &r.Posts
}
