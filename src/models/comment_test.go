package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"strings"
	"testing"
)

func setupTestCommentDB(t *testing.T) (UserType, PostType) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := UserType{}
	user.Username = "commentauthor"
	user.Email = "comment@test.com"
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	cat := CategoryType{}
	cat.Name = "commentcat"
	cat.Description = "Comment test category"
	cat.Add()

	post := PostType{}
	post.Title = "Comment Test Post"
	post.Body = "Comment Test Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()
	return user, post
}

func TestCommentValidateComment(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr error
	}{
		{"valid comment", "Hello world", nil},
		{"empty comment", "", ferror.ErrorCommentEmpty},
		{"too long comment", strings.Repeat("a", 1001), ferror.ErrorCommentTooLong},
		{"max length", strings.Repeat("a", 1000), nil},
		{"single char", "a", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommentType{}
			c.Body = tt.body
			err := c.ValidateComment()
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateComment() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateComment() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCommentAdd(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Test comment"
	err := c.Add()
	if err != nil {
		t.Fatalf("Add() error: %v", err)
	}
	if c.Id == 0 {
		t.Error("Add() did not set comment ID")
	}

	c2 := &CommentType{}
	c2.UserId = user.Id
	c2.PostId = post.Id
	c2.Body = ""
	err = c2.Add()
	if err == nil {
		t.Error("Add() should fail with empty body")
	}
}

func TestCommentGetReactions(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Reaction comment"
	c.Add()

	err := c.GetReactions()
	if err != nil {
		t.Fatalf("GetReactions() error: %v", err)
	}
	if c.Likes != 0 {
		t.Errorf("GetReactions() Likes = %d, want 0", c.Likes)
	}
	if c.Dislikes != 0 {
		t.Errorf("GetReactions() Dislikes = %d, want 0", c.Dislikes)
	}

	db.AddLikeToComment(user.Id, c.Id)
	db.AddLikeToComment(99, c.Id)
	db.AddDislikeToComment(user.Id, c.Id)

	c.Likes = 0
	c.Dislikes = 0
	err = c.GetReactions()
	if err != nil {
		t.Fatalf("GetReactions() error: %v", err)
	}
	if c.Likes != 2 {
		t.Errorf("GetReactions() Likes = %d, want 2", c.Likes)
	}
	if c.Dislikes != 1 {
		t.Errorf("GetReactions() Dislikes = %d, want 1", c.Dislikes)
	}
}

func TestCommentGetReactionsByUserId(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Reaction UID comment"
	c.Add()

	err := c.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if c.Liked || c.Disliked {
		t.Error("GetReactionsByUserId() should return false before any reactions")
	}

	db.AddLikeToComment(user.Id, c.Id)
	c.Liked = false
	c.Disliked = false
	err = c.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if !c.Liked {
		t.Error("GetReactionsByUserId() should return true for Liked")
	}

	db.RemoveReaction(c.Id)
	db.AddDislikeToComment(user.Id, c.Id)
	c.Liked = false
	c.Disliked = false
	err = c.GetReactionsByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetReactionsByUserId() error: %v", err)
	}
	if c.Liked {
		t.Error("GetReactionsByUserId() should return false for Liked after dislike")
	}
	if !c.Disliked {
		t.Error("GetReactionsByUserId() should return true for Disliked")
	}
}

func TestCommentGetById(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "GetById comment"
	c.Add()

	c2 := &CommentType{}
	c2.Id = c.Id
	err := c2.GetById()
	if err != nil {
		t.Fatalf("GetById() error: %v", err)
	}
	if c2.Body != "GetById comment" {
		t.Errorf("GetById() Body = %q, want %q", c2.Body, "GetById comment")
	}

	c3 := &CommentType{}
	c3.Id = 99999
	err = c3.GetById()
	if err == nil {
		t.Error("GetById() should fail for nonexistent comment")
	}
}

func TestCommentUpdate(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Original comment"
	c.Add()

	c.Body = "Updated comment"
	err := c.Update()
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	c2 := &CommentType{}
	c2.Id = c.Id
	err = c2.SelectCommentById()
	if err != nil {
		t.Fatalf("SelectCommentById() error: %v", err)
	}
	if c2.Body != "Updated comment" {
		t.Errorf("Update() Body = %q, want %q", c2.Body, "Updated comment")
	}

	c3 := &CommentType{}
	c3.Id = c.Id
	c3.Body = ""
	err = c3.Update()
	if err == nil {
		t.Error("Update() should fail with empty body")
	}
}

func TestCommentDeleteById(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Delete me"
	c.Add()

	err := c.DeleteCommentById()
	if err != nil {
		t.Fatalf("DeleteCommentById() error: %v", err)
	}

	c2 := &CommentType{}
	c2.Id = c.Id
	err = c2.SelectCommentById()
	if err == nil {
		t.Error("Comment should not exist after deletion")
	}
}

func TestCommentSelectCommentById(t *testing.T) {
	user, post := setupTestCommentDB(t)

	c := &CommentType{}
	c.UserId = user.Id
	c.PostId = post.Id
	c.Body = "Select test"
	c.Add()

	c2 := &CommentType{}
	c2.Id = c.Id
	err := c2.SelectCommentById()
	if err != nil {
		t.Fatalf("SelectCommentById() error: %v", err)
	}
	if c2.Body != "Select test" {
		t.Errorf("SelectCommentById() Body = %q, want %q", c2.Body, "Select test")
	}
	if c2.UserId != user.Id {
		t.Errorf("SelectCommentById() UserId = %d, want %d", c2.UserId, user.Id)
	}
	if c2.PostId != post.Id {
		t.Errorf("SelectCommentById() PostId = %d, want %d", c2.PostId, post.Id)
	}
}
