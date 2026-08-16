package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"testing"
)

func setupTestPostDB(t *testing.T) UserType {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := UserType{}
	user.Username = "postauthor"
	user.Email = "author@test.com"
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	cat := CategoryType{}
	cat.Name = "general"
	cat.Description = "General discussion"
	cat.Add()
	return user
}

func createTestCategory(t *testing.T, name string) CategoryType {
	t.Helper()
	cat := CategoryType{}
	cat.Name = name
	cat.Description = name + " description"
	if err := cat.Add(); err != nil {
		t.Fatalf("Failed to create category %q: %v", name, err)
	}
	return cat
}

func TestPostValidatePost(t *testing.T) {
	cat := CategoriesType{CategoryType{CategoryRowType: db.CategoryRowType{Id: 1, Name: "test", Description: "test"}}}

	tests := []struct {
		name    string
		title   string
		body    string
		cats    CategoriesType
		wantErr error
	}{
		{"valid post", "Title", "Body", cat, nil},
		{"empty title", "", "Body", cat, ferror.ErrorPostTitleEmpty},
		{"empty body", "Title", "", cat, ferror.ErrorPostBodyEmpty},
		{"no categories", "Title", "Body", CategoriesType{}, ferror.ErrorPostHasNoCategory},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostType{}
			p.Title = tt.title
			p.Body = tt.body
			p.Categories = tt.cats
			err := p.ValidatePost()
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidatePost() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePost() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPostAdd(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "addcat")}

	p := &PostType{}
	p.Title = "Test Post"
	p.Body = "Test Body"
	p.UserId = user.Id
	p.Categories = cat
	err := p.Add()
	if err != nil {
		t.Fatalf("Add() error: %v", err)
	}
	if p.Id == 0 {
		t.Error("Add() did not set post ID")
	}

	p2 := &PostType{}
	p2.Title = ""
	p2.Body = "Body"
	p2.UserId = user.Id
	p2.Categories = cat
	err = p2.Add()
	if err == nil {
		t.Error("Add() should fail with empty title")
	}
}

func TestPostGetReactions(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "reactcat")}

	p := &PostType{}
	p.Title = "Reaction Post"
	p.Body = "Reaction Body"
	p.UserId = user.Id
	p.Categories = cat
	p.Add()

	err := p.GetReactions()
	if err != nil {
		t.Fatalf("GetReactions() error: %v", err)
	}
	if p.Likes != 0 {
		t.Errorf("GetReactions() Likes = %d, want 0", p.Likes)
	}
	if p.Dislikes != 0 {
		t.Errorf("GetReactions() Dislikes = %d, want 0", p.Dislikes)
	}

	db.AddLikeToPost(user.Id, p.Id)
	db.AddLikeToPost(user.Id+1, p.Id)
	db.AddDislikeToPost(user.Id, p.Id)

	p.Likes = 0
	p.Dislikes = 0
	err = p.GetReactions()
	if err != nil {
		t.Fatalf("GetReactions() error: %v", err)
	}
	if p.Likes != 2 {
		t.Errorf("GetReactions() Likes = %d, want 2", p.Likes)
	}
	if p.Dislikes != 1 {
		t.Errorf("GetReactions() Dislikes = %d, want 1", p.Dislikes)
	}
}

func TestPostGetReactionsByUserId(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "reactuidcat")}

	p := &PostType{}
	p.Title = "Reaction UID Post"
	p.Body = "Reaction UID Body"
	p.UserId = user.Id
	p.Categories = cat
	p.Add()

	err := p.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if p.Liked || p.Disliked {
		t.Error("GetReactionsByUserId() should return false before any reactions")
	}

	db.AddLikeToPost(user.Id, p.Id)
	p.Liked = false
	p.Disliked = false
	err = p.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if !p.Liked {
		t.Error("GetReactionsByUserId() should return true for Liked")
	}
	if p.Disliked {
		t.Error("GetReactionsByUserId() should return false for Disliked")
	}

	db.RemoveLikeFromPost(user.Id, p.Id)
	db.AddDislikeToPost(user.Id, p.Id)
	p.Liked = false
	p.Disliked = false
	err = p.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if p.Liked {
		t.Error("GetReactionsByUserId() should return false for Liked after dislike")
	}
	if !p.Disliked {
		t.Error("GetReactionsByUserId() should return true for Disliked")
	}
}

func TestPostGetComments(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "commentcat")}

	p := &PostType{}
	p.Title = "Comment Post"
	p.Body = "Comment Body"
	p.UserId = user.Id
	p.Categories = cat
	p.Add()

	err := p.GetComments()
	if err != nil {
		t.Fatalf("GetComments() error: %v", err)
	}
	if len(p.Comments) != 0 {
		t.Errorf("GetComments() returned %d comments, want 0", len(p.Comments))
	}

	c := CommentType{}
	c.UserId = user.Id
	c.PostId = p.Id
	c.Body = "test comment"
	c.Add()

	p.Comments = nil
	err = p.GetComments()
	if err != nil {
		t.Fatalf("GetComments() error: %v", err)
	}
	if len(p.Comments) != 1 {
		t.Errorf("GetComments() returned %d comments, want 1", len(p.Comments))
	}
	if p.Comments[0].Body != "test comment" {
		t.Errorf("GetComments() Body = %q, want %q", p.Comments[0].Body, "test comment")
	}
}

func TestPostDelete(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "delcat")}

	p := &PostType{}
	p.Title = "Delete Post"
	p.Body = "Delete Body"
	p.UserId = user.Id
	p.Categories = cat
	p.Add()

	err := p.Delete()
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	p2 := &PostType{}
	p2.Id = p.Id
	err = p2.SelectPostById()
	if err == nil {
		t.Error("Post should not exist after deletion")
	}
}

func TestPostUpdate(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "upcat")}

	p := &PostType{}
	p.Title = "Original Title"
	p.Body = "Original Body"
	p.UserId = user.Id
	p.Categories = cat
	p.Add()

	p.Title = "Updated Title"
	p.Body = "Updated Body"
	err := p.Update()
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	p2 := &PostType{}
	p2.Id = p.Id
	err = p2.SelectPostById()
	if err != nil {
		t.Fatalf("SelectPostById() error: %v", err)
	}
	if p2.Title != "Updated Title" {
		t.Errorf("Update() Title = %q, want %q", p2.Title, "Updated Title")
	}
	if p2.Body != "Updated Body" {
		t.Errorf("Update() Body = %q, want %q", p2.Body, "Updated Body")
	}

	p3 := &PostType{}
	p3.Title = ""
	p3.Body = "Body"
	p3.Id = p.Id
	p3.Categories = cat
	err = p3.Update()
	if err == nil {
		t.Error("Update() should fail with empty title")
	}
}

func TestPostsGetPosts(t *testing.T) {
	user := setupTestPostDB(t)
	cat := CategoriesType{createTestCategory(t, "allcat")}

	p1 := &PostType{}
	p1.Title = "Post 1"
	p1.Body = "Body 1"
	p1.UserId = user.Id
	p1.Categories = cat
	p1.Add()

	p2 := &PostType{}
	p2.Title = "Post 2"
	p2.Body = "Body 2"
	p2.UserId = user.Id
	p2.Categories = cat
	p2.Add()

	var posts PostsType
	err := posts.GetPosts()
	if err != nil {
		t.Fatalf("GetPosts() error: %v", err)
	}
	if len(posts) < 2 {
		t.Errorf("GetPosts() returned %d posts, want >= 2", len(posts))
	}
}

func TestPostsGetPostsByCategoryId(t *testing.T) {
	user := setupTestPostDB(t)
	cat1 := createTestCategory(t, "filtercat1")
	cat2 := createTestCategory(t, "filtercat2")

	p1 := &PostType{}
	p1.Title = "Cat1 Post"
	p1.Body = "Cat1 Body"
	p1.UserId = user.Id
	p1.Categories = CategoriesType{cat1}
	p1.Add()

	p2 := &PostType{}
	p2.Title = "Cat2 Post"
	p2.Body = "Cat2 Body"
	p2.UserId = user.Id
	p2.Categories = CategoriesType{cat2}
	p2.Add()

	var posts PostsType
	err := posts.GetPostsByCategoryId(cat1.Id)
	if err != nil {
		t.Fatalf("GetPostsByCategoryId() error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("GetPostsByCategoryId() returned %d posts, want 1", len(posts))
	}
	if len(posts) > 0 && posts[0].Title != "Cat1 Post" {
		t.Errorf("GetPostsByCategoryId() Title = %q, want %q", posts[0].Title, "Cat1 Post")
	}
}
