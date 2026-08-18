package models

import (
	"forum/src/db"
	"forum/src/utils"
	"testing"
)

func setupTestReactionDB(t *testing.T) (UserType, PostType, CommentType) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := UserType{}
	user.Username = "reactuser"
	user.Email = "react@test.com"
	user.Hash = hash
	user.Gender = "male"
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	cat := CategoryType{}
	cat.Name = "reactcat"
	cat.Description = "Reaction category"
	cat.Add()

	post := PostType{}
	post.Title = "Reaction Post"
	post.Body = "Reaction Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	comment := CommentType{}
	comment.UserId = user.Id
	comment.PostId = post.Id
	comment.Body = "Reaction comment"
	comment.Add()

	return user, post, comment
}

func TestGetPostLikesByUserId(t *testing.T) {
	user, post, _ := setupTestReactionDB(t)

	db.InsertLikeToPost(user.Id, post.Id)

	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetPostLikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetPostLikesByUserId() returned %d reactions, want 1", len(reactions))
	}
	if len(reactions) > 0 && reactions[0].PostId != post.Id {
		t.Errorf("GetPostLikesByUserId() PostId = %d, want %d", reactions[0].PostId, post.Id)
	}
}

func TestGetPostLikesByUserIdEmpty(t *testing.T) {
	user, _, _ := setupTestReactionDB(t)

	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetPostLikesByUserId() error: %v", err)
	}
	if len(reactions) != 0 {
		t.Errorf("GetPostLikesByUserId() returned %d reactions, want 0", len(reactions))
	}
}

func TestGetPostDislikesByUserId(t *testing.T) {
	user, post, _ := setupTestReactionDB(t)

	db.InsertDislikeToPost(user.Id, post.Id)

	reactions, err := GetPostDislikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetPostDislikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetPostDislikesByUserId() returned %d reactions, want 1", len(reactions))
	}
	if len(reactions) > 0 && reactions[0].PostId != post.Id {
		t.Errorf("GetPostDislikesByUserId() PostId = %d, want %d", reactions[0].PostId, post.Id)
	}
}

func TestGetPostDislikesByUserIdEmpty(t *testing.T) {
	user, _, _ := setupTestReactionDB(t)

	reactions, err := GetPostDislikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetPostDislikesByUserId() error: %v", err)
	}
	if len(reactions) != 0 {
		t.Errorf("GetPostDislikesByUserId() returned %d reactions, want 0", len(reactions))
	}
}

func TestGetCommentLikesByUserId(t *testing.T) {
	user, _, comment := setupTestReactionDB(t)

	db.InsertLikeToComment(user.Id, comment.Id)

	reactions, err := GetCommentLikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetCommentLikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetCommentLikesByUserId() returned %d reactions, want 1", len(reactions))
	}
	if len(reactions) > 0 && reactions[0].CommentId != comment.Id {
		t.Errorf("GetCommentLikesByUserId() CommentId = %d, want %d", reactions[0].CommentId, comment.Id)
	}
}

func TestGetCommentLikesByUserIdEmpty(t *testing.T) {
	user, _, _ := setupTestReactionDB(t)

	reactions, err := GetCommentLikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetCommentLikesByUserId() error: %v", err)
	}
	if len(reactions) != 0 {
		t.Errorf("GetCommentLikesByUserId() returned %d reactions, want 0", len(reactions))
	}
}

func TestGetCommentDislikesByUserId(t *testing.T) {
	user, _, comment := setupTestReactionDB(t)

	db.InsertDislikeToComment(user.Id, comment.Id)

	reactions, err := GetCommentDisikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetCommentDisikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetCommentDisikesByUserId() returned %d reactions, want 1", len(reactions))
	}
	if len(reactions) > 0 && reactions[0].CommentId != comment.Id {
		t.Errorf("GetCommentDisikesByUserId() CommentId = %d, want %d", reactions[0].CommentId, comment.Id)
	}
}

func TestGetCommentDislikesByUserIdEmpty(t *testing.T) {
	user, _, _ := setupTestReactionDB(t)

	reactions, err := GetCommentDisikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetCommentDisikesByUserId() error: %v", err)
	}
	if len(reactions) != 0 {
		t.Errorf("GetCommentDisikesByUserId() returned %d reactions, want 0", len(reactions))
	}
}

func TestGetPostLikesByUserIdMultiple(t *testing.T) {
	user, post, _ := setupTestReactionDB(t)

	user2 := UserType{}
	user2.Username = "reactuser2"
	user2.Email = "react2@test.com"
	user2.Hash, _ = utils.HashPassword("password123")
	user2.Gender = "male"
	user2.Add()

	db.InsertLikeToPost(user.Id, post.Id)
	db.InsertLikeToPost(user2.Id, post.Id)

	var reactions ReactionsType
	err := reactions.GetPostLikesByUserId(user.Id)
	if err != nil {
		t.Fatalf("GetPostLikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetPostLikesByUserId() returned %d reactions for user1, want 1", len(reactions))
	}
	err = reactions.GetPostLikesByUserId(user2.Id)
	if err != nil {
		t.Fatalf("GetPostLikesByUserId() error: %v", err)
	}
	if len(reactions) != 1 {
		t.Errorf("GetPostLikesByUserId() returned %d reactions for user2, want 1", len(reactions))
	}
	count, err := db.SelectLikesCountByPostId(post.Id)
	if count != 2 {
		t.Errorf("GetPostLikesByUserId() returned %d reactions after second call, want 2 (cumulative append)", count)
	}
}

func TestDeleteReactionById(t *testing.T) {
	user, post, _ := setupTestReactionDB(t)

	db.InsertLikeToPost(user.Id, post.Id)

	var reactions ReactionsType
	reactions.GetPostLikesByUserId(user.Id)
	if len(reactions) != 1 {
		t.Fatalf("Setup failed: expected 1 like, got %d", len(reactions))
	}

	err := db.DeleteReactionById(reactions[0].Id)
	if err != nil {
		t.Fatalf("DeleteReactionById() error: %v", err)
	}

	reactions = nil
	reactions.GetPostLikesByUserId(user.Id)
	if len(reactions) != 0 {
		t.Errorf("RemoveReaction() left %d reactions, want 0", len(reactions))
	}
}
