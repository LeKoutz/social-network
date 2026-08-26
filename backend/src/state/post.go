package state

import "forum/src/models"

func (r *State) SetEditPost(b bool) {
	r.EditingPost = b
}

func (r *State) InitPost() {
	var post models.PostType
	r.Posts = models.PostsType{post}
}

func (r *State) EditPost() *models.PostType {
	if r.Posts == nil {
		r.InitPost()
	}
	return &r.Posts[0]
}

func (r *State) GetPost() models.PostType {
	return *r.EditPost()
}

func (r *State) SetPost(post models.PostType) *State {
	r.Posts = models.PostsType{post}
	return r
}
