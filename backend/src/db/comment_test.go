package db

import (
	"forum/src/utils"
	"testing"
)

func setupCommentRowTestDB(t *testing.T) (int64, int64) {
	t.Helper()
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("InitDB(%q) failed: %v", ":memory:", err)
	}
	user := UserRowType{}
	user.Username = "cimguser"
	user.Email = "cimguser@test.com"
	user.Hash = "x"
	if err := user.InsertUserWithHash(); err != nil {
		t.Fatalf("create user: %v", err)
	}
	post := PostRowType{}
	post.Title = "Comment image post"
	post.Body = "Body"
	post.UserId = user.Id
	if err := post.InsertPost(); err != nil {
		t.Fatalf("create post: %v", err)
	}
	return user.Id, post.Id
}

func TestCommentRowImagePath(t *testing.T) {
	userId, postId := setupCommentRowTestDB(t)

	c := CommentRowType{}
	c.PostId = postId
	c.UserId = userId
	c.Body = "With image"
	c.ImagePath = "uploads/images/cimg.png"
	if err := c.InsertComment(); err != nil {
		t.Fatalf("InsertComment() error: %v", err)
	}
	if c.Id == 0 {
		t.Fatal("InsertComment() did not set ID")
	}

	var got CommentRowType
	got.Id = c.Id
	if err := got.SelectCommentById(); err != nil {
		t.Fatalf("SelectCommentById() error: %v", err)
	}
	if got.ImagePath != "uploads/images/cimg.png" {
		t.Errorf("SelectCommentById() ImagePath = %q, want %q", got.ImagePath, "uploads/images/cimg.png")
	}
	if got.PostId != postId || got.UserId != userId || got.Body != "With image" {
		t.Errorf("SelectCommentById() = %+v, want post_id=%d user_id=%d body=%q", got, postId, userId, "With image")
	}

	rows, err := SelectCommentsByUserId(userId)
	if err != nil {
		t.Fatalf("SelectCommentsByUserId() error: %v", err)
	}
	if len(rows) != 1 || rows[0].ImagePath != "uploads/images/cimg.png" {
		t.Errorf("SelectCommentsByUserId() emails = %+v, want 1 row with ImagePath %q", rows, "uploads/images/cimg.png")
	}

	got.ImagePath = "uploads/images/cimg-edited.png"
	got.Body = "Edited"
	if err := got.UpdateCommentById(); err != nil {
		t.Fatalf("UpdateCommentById() error: %v", err)
	}
	var got2 CommentRowType
	got2.Id = c.Id
	if err := got2.SelectCommentById(); err != nil {
		t.Fatalf("SelectCommentById() error: %v", err)
	}
	if got2.ImagePath != "uploads/images/cimg-edited.png" || got2.Body != "Edited" {
		t.Errorf("UpdateCommentById() -> %+v, want ImagePath %q Body %q", got2, "uploads/images/cimg-edited.png", "Edited")
	}
}

func TestCommentRowImagePathLegacyNull(t *testing.T) {
	userId, postId := setupCommentRowTestDB(t)

	r, err := db.Exec(
		"INSERT INTO comments (post_id, user_id, body, timestamp) VALUES (?, ?, ?, ?)",
		postId,
		userId,
		"Legacy comment",
		utils.GetCurrentTimestamp(),
	)
	if err != nil {
		t.Fatalf("legacy insert: %v", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		t.Fatalf("legacy LastInsertId: %v", err)
	}

	c := CommentRowType{}
	c.Id = id
	if err := c.SelectCommentById(); err != nil {
		t.Fatalf("SelectCommentById() error: %v", err)
	}
	if c.ImagePath != "" {
		t.Errorf("SelectCommentById() ImagePath = %q, want empty for legacy NULL row", c.ImagePath)
	}
	p := PostRowType{}
	p.Id = postId
	rows, err := p.SelectCommentsAndUsernameByPostId()
	if err != nil {
		t.Fatalf("SelectCommentsAndUsernameByPostId() error: %v", err)
	}
	if len(rows) != 1 || rows[0].ImagePath != "" {
		t.Errorf("SelectCommentsAndUsernameByPostId() = %+v, want 1 row with empty ImagePath", rows)
	}
}
