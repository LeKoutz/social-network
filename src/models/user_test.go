package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"strings"
	"testing"
)

func setupTestUserDB(t *testing.T) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
}

func createTestUser(t *testing.T, username, email, password string) UserType {
	t.Helper()
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	user := UserType{}
	user.Username = username
	user.Email = email
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestGetGuestUser(t *testing.T) {
	u := GetGuestUser()
	if u.Username != "guest" {
		t.Errorf("GetGuestUser().Username = %q, want %q", u.Username, "guest")
	}
	if u.LoggedIn != false {
		t.Errorf("GetGuestUser().LoggedIn = %v, want false", u.LoggedIn)
	}
}

func TestUserValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid alphanumeric", "testuser", false},
		{"valid with underscore", "test_user", false},
		{"valid with numbers", "user123", false},
		{"too short", "usr", true},
		{"too long", strings.Repeat("a", 51), true},
		{"with spaces", "test user", true},
		{"with special chars", "test@user", true},
		{"empty", "", true},
		{"exactly 4 chars", "test", false},
		{"exactly 50 chars", strings.Repeat("a", 50), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserType{}
			u.Username = tt.username
			err := u.ValidateUsername()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@test.com", false},
		{"valid with subdomain", "user@mail.test.com", false},
		{"invalid no at", "usertest.com", true},
		{"invalid no domain", "user@", true},
		{"empty", "", true},
		{"invalid chars", "user @test.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserType{}
			u.Email = tt.email
			err := u.ValidateEmail()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserAdd(t *testing.T) {
	setupTestUserDB(t)

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "valid user",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "duplicate email",
			username: "anotheruser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  ferror.ErrorEmailIsRegistered,
		},
		{
			name:     "invalid username",
			username: "ab",
			email:    "short@example.com",
			password: "password123",
			wantErr:  ferror.ErrorInvalidUsername,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, _ := utils.HashPassword(tt.password)
			u := &UserType{}
			u.Username = tt.username
			u.Email = tt.email
			u.Hash = hash
			err := u.Add()
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Add() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Add() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserAddOAuth(t *testing.T) {
	setupTestUserDB(t)

	user := UserType{}
	user.Username = "oauthuser"
	user.Email = "oauth@test.com"
	user.OAuthProvider = "google"
	err := user.AddOAuth()
	if err != nil {
		t.Fatalf("AddOAuth() unexpected error: %v", err)
	}

	err = user.AddOAuth()
	if err == nil {
		t.Error("AddOAuth() should fail for duplicate email")
	}

	user2 := UserType{}
	user2.Username = "oauthuser"
	user2.Email = "oauth2@test.com"
	user2.OAuthProvider = "google"
	err = user2.AddOAuth()
	if err == nil {
		t.Error("AddOAuth() should fail for duplicate username")
	}
}

func TestUserGetByEmail(t *testing.T) {
	setupTestUserDB(t)
	createTestUser(t, "emailuser", "email@test.com", "pass123456!")

	u := &UserType{}
	u.Identifier = "email@test.com"
	err := u.GetUserByIdentifier()
	if err != nil {
		t.Fatalf("GetUserByEmail() error: %v", err)
	}
	if u.Username != "emailuser" {
		t.Errorf("GetUserByEmail() Username = %q, want %q", u.Username, "emailuser")
	}

	u2 := &UserType{}
	u2.Identifier = "nonexistent@test.com"
	err = u2.GetUserByIdentifier()
	if err == nil {
		t.Error("GetUserByEmail() should fail for nonexistent email")
	}
}

func TestUserGetPasswordByEmail(t *testing.T) {
	setupTestUserDB(t)
	createTestUser(t, "passuser", "pass@test.com", "mypassword123!")

	u := &UserType{}
	u.Identifier = "pass@test.com"
	err := u.GetUserPasswordByIdentifier()
	if err != nil {
		t.Fatalf("GetUserPasswordByEmail() error: %v", err)
	}
	if len(u.Hash) == 0 {
		t.Error("GetUserPasswordByEmail() returned empty hash")
	}

	u2 := &UserType{}
	u2.Identifier = "nouser@test.com"
	err = u2.GetUserPasswordByIdentifier()
	if err == nil {
		t.Error("GetUserPasswordByEmail() should fail for nonexistent email")
	}
}

func TestUserSetAndGetSession(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "sessuser", "sess@test.com", "pass123456!")

	err := user.SetUserSession("test-session-key")
	if err != nil {
		t.Fatalf("SetUserSession() error: %v", err)
	}

	u := &UserType{}
	u.SessionId = "test-session-key"
	err = u.GetUserBySession()
	if err != nil {
		t.Fatalf("GetUserBySession() error: %v", err)
	}
	if u.Id != user.Id {
		t.Errorf("GetUserBySession() Id = %d, want %d", u.Id, user.Id)
	}

	u2 := &UserType{}
	u2.SessionId = "nonexistent-session"
	err = u2.GetUserBySession()
	if err == nil {
		t.Error("GetUserBySession() should fail for nonexistent session")
	}
}

func TestIsUniqueUsername(t *testing.T) {
	setupTestUserDB(t)
	createTestUser(t, "uniqueuser", "unique@test.com", "pass123456!")

	if IsUniqueUsername("uniqueuser") {
		t.Error("IsUniqueUsername() should return false for taken username")
	}
	if !IsUniqueUsername("otheruser") {
		t.Error("IsUniqueUsername() should return true for available username")
	}
}

func TestIsUniqueEmail(t *testing.T) {
	setupTestUserDB(t)
	createTestUser(t, "emailcheck", "check@test.com", "pass123456!")

	if IsUniqueEmail("check@test.com") {
		t.Error("IsUniqueEmail() should return false for registered email")
	}
	if !IsUniqueEmail("other@test.com") {
		t.Error("IsUniqueEmail() should return true for available email")
	}
}

func TestIsEmailRegistered(t *testing.T) {
	setupTestUserDB(t)
	createTestUser(t, "regcheck", "reg@test.com", "pass123456!")

	if !IsEmailRegistered("reg@test.com") {
		t.Error("IsEmailRegistered() should return true for registered email")
	}
	if IsEmailRegistered("unreg@test.com") {
		t.Error("IsEmailRegistered() should return false for unregistered email")
	}
}

func TestHasUserLikedPost(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "likeuser", "like@test.com", "pass123456!")

	cat := CategoryType{}
	cat.Name = "testcat"
	cat.Description = "test category"
	cat.Add()

	post := PostType{}
	post.Title = "Test Post"
	post.Body = "Test Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	liked, err := HasUserLikedPost(user.Id, post.Id)
	if err != nil {
		t.Fatalf("HasUserLikedPost() error: %v", err)
	}
	if liked {
		t.Error("HasUserLikedPost() should return false initially")
	}

	db.AddLikeToPost(user.Id, post.Id)

	liked, err = HasUserLikedPost(user.Id, post.Id)
	if err != nil {
		t.Fatalf("HasUserLikedPost() error: %v", err)
	}
	if !liked {
		t.Error("HasUserLikedPost() should return true after liking")
	}
}

func TestHasUserDislikedPost(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "dislikeuser", "dislike@test.com", "pass123456!")

	cat := CategoryType{}
	cat.Name = "testcat2"
	cat.Description = "test category 2"
	cat.Add()

	post := PostType{}
	post.Title = "Test Post"
	post.Body = "Test Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	disliked, err := HasUserDislikedPost(user.Id, post.Id)
	if err != nil {
		t.Fatalf("HasUserDislikedPost() error: %v", err)
	}
	if disliked {
		t.Error("HasUserDislikedPost() should return false initially")
	}

	db.AddDislikeToPost(user.Id, post.Id)

	disliked, err = HasUserDislikedPost(user.Id, post.Id)
	if err != nil {
		t.Fatalf("HasUserDislikedPost() error: %v", err)
	}
	if !disliked {
		t.Error("HasUserDislikedPost() should return true after disliking")
	}
}

func TestHasUserLikedComment(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "clikeuser", "clike@test.com", "pass123456!")

	cat := CategoryType{}
	cat.Name = "testcat3"
	cat.Description = "test category 3"
	cat.Add()

	post := PostType{}
	post.Title = "Test Post"
	post.Body = "Test Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	comment := CommentType{}
	comment.UserId = user.Id
	comment.PostId = post.Id
	comment.Body = "test comment"
	comment.Add()

	liked, err := HasUserLikedComment(user.Id, comment.Id)
	if err != nil {
		t.Fatalf("HasUserLikedComment() error: %v", err)
	}
	if liked {
		t.Error("HasUserLikedComment() should return false initially")
	}

	db.AddLikeToComment(user.Id, comment.Id)

	liked, err = HasUserLikedComment(user.Id, comment.Id)
	if err != nil {
		t.Fatalf("HasUserLikedComment() error: %v", err)
	}
	if !liked {
		t.Error("HasUserLikedComment() should return true after liking")
	}
}

func TestHasUserDislikedComment(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "cdislikeuser", "cdislike@test.com", "pass123456!")

	cat := CategoryType{}
	cat.Name = "testcat4"
	cat.Description = "test category 4"
	cat.Add()

	post := PostType{}
	post.Title = "Test Post"
	post.Body = "Test Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	comment := CommentType{}
	comment.UserId = user.Id
	comment.PostId = post.Id
	comment.Body = "test comment"
	comment.Add()

	disliked, err := HasUserDislikedComment(user.Id, comment.Id)
	if err != nil {
		t.Fatalf("HasUserDislikedComment() error: %v", err)
	}
	if disliked {
		t.Error("HasUserDislikedComment() should return false initially")
	}

	db.AddDislikeToComment(user.Id, comment.Id)

	disliked, err = HasUserDislikedComment(user.Id, comment.Id)
	if err != nil {
		t.Fatalf("HasUserDislikedComment() error: %v", err)
	}
	if !disliked {
		t.Error("HasUserDislikedComment() should return true after disliking")
	}
}

func TestUserGetPosts(t *testing.T) {
	setupTestUserDB(t)
	user := createTestUser(t, "postuser", "post@test.com", "pass123456!")

	cat := CategoryType{}
	cat.Name = "testcat5"
	cat.Description = "test category 5"
	cat.Add()

	post := PostType{}
	post.Title = "User Post"
	post.Body = "User Post Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	posts, err := user.GetPosts()
	if err != nil {
		t.Fatalf("GetPosts() error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("GetPosts() returned %d posts, want 1", len(posts))
	}

	user2 := createTestUser(t, "postuser2", "post2@test.com", "pass123456!")
	posts2, err := user2.GetPosts()
	if err != nil {
		t.Fatalf("GetPosts() error: %v", err)
	}
	if len(posts2) != 0 {
		t.Errorf("GetPosts() returned %d posts for user with no posts, want 0", len(posts2))
	}
}

func TestUserCountUnreadNotifications(t *testing.T) {
	u := &UserType{}
	u.Notifications = NotificationsType{
		{NotificationRowType: db.NotificationRowType{Read: false}},
		{NotificationRowType: db.NotificationRowType{Read: true}},
		{NotificationRowType: db.NotificationRowType{Read: false}},
	}
	u.CountUnreadNotifications()
	if u.UnreadNotificationsCount != 2 {
		t.Errorf("CountUnreadNotifications() = %d, want 2", u.UnreadNotificationsCount)
	}
}

func TestUserGetUserByOAuthProviderAndEmail(t *testing.T) {
	setupTestUserDB(t)
	user := UserType{}
	user.Username = "oauthlookup"
	user.Email = "oauthlookup@test.com"
	user.OAuthProvider = "google"
	user.AddOAuth()

	u := &UserType{}
	u.OAuthProvider = "google"
	u.Email = "oauthlookup@test.com"
	err := u.GetUserByOAuthProviderAndEmail()
	if err != nil {
		t.Fatalf("GetUserByOAuthProviderAndEmail() error: %v", err)
	}
	if u.Username != "oauthlookup" {
		t.Errorf("GetUserByOAuthProviderAndEmail() Username = %q, want %q", u.Username, "oauthlookup")
	}

	u2 := &UserType{}
	u2.OAuthProvider = "github"
	u2.Email = "oauthlookup@test.com"
	err = u2.GetUserByOAuthProviderAndEmail()
	if err == nil {
		t.Error("GetUserByOAuthProviderAndEmail() should fail for wrong provider")
	}
}
