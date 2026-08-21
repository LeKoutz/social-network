package state

import (
	"forum/src/models"
	"sort"
)

func (r *State) SetPosts(posts models.PostsType) *State {
	r.Posts = posts
	sort.Slice(r.Posts, func(i, j int) bool {
		return r.Posts[i].Timestamp > r.Posts[j].Timestamp
	})
	return r
}

func (r *State) EditPosts() *models.PostsType {
	return &r.Posts
}
