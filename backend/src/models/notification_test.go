package models

import (
	"forum/src/db"
	"forum/src/utils"
	"testing"
)

func setupTestNotificationDB(t *testing.T) (UserType, UserType, PostType) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	hash, _ := utils.HashPassword("password123")

	owner := UserType{}
	owner.Username = "notifowner"
	owner.Email = "owner@test.com"
	owner.Hash = hash
	owner.Gender = "male"
	if err := owner.Add(); err != nil {
		t.Fatalf("Failed to create owner user: %v", err)
	}

	actor := UserType{}
	actor.Username = "notifactor"
	actor.Email = "actor@test.com"
	actor.Hash = hash
	actor.Gender = "male"
	if err := actor.Add(); err != nil {
		t.Fatalf("Failed to create actor user: %v", err)
	}

	cat := CategoryType{}
	cat.Name = "notifcat"
	cat.Description = "Notification category"
	cat.Add()

	post := PostType{}
	post.Title = "Notif Post"
	post.Body = "Notif Body"
	post.UserId = owner.Id
	post.Categories = CategoriesType{cat}
	post.Add()

	return owner, actor, post
}

func TestCreateNotification(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	n := NotificationType{}
	n.UserId = owner.Id
	n.ActorId = actor.Id
	n.Type = "comment"
	n.PostId = post.Id
	err := CreateNotification(n)
	if err != nil {
		t.Fatalf("CreateNotification() error: %v", err)
	}
}

func TestCreateNotificationSelfAction(t *testing.T) {
	owner, _, post := setupTestNotificationDB(t)

	n := NotificationType{}
	n.UserId = owner.Id
	n.ActorId = owner.Id
	n.Type = "like"
	n.PostId = post.Id
	err := CreateNotification(n)
	if err != nil {
		t.Errorf("CreateNotification() should not error for self-action, got: %v", err)
	}

	notifications, err := db.SelectNotificationsByUserId(owner.Id)
	if err != nil {
		t.Fatalf("GetNotificationsByUserId() error: %v", err)
	}
	for _, notif := range notifications {
		if notif.UserId == notif.ActorId {
			t.Error("Self-notification should not be created")
		}
	}
}

func TestCreateCommentNotification(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)
	post.User = owner

	comment := CommentType{}
	comment.UserId = actor.Id
	comment.PostId = post.Id
	comment.Body = "test comment"
	comment.Add()

	err := comment.CreateCommentNotification(post)
	if err != nil {
		t.Fatalf("CreateCommentNotification() error: %v", err)
	}

	notifications, err := db.SelectNotificationsByUserId(owner.Id)
	if err != nil {
		t.Fatalf("GetNotificationsByUserId() error: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].Type != "comment" {
		t.Errorf("Notification Type = %q, want %q", notifications[0].Type, "comment")
	}
}

func TestCreateReactionNotificationPost(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	post.User = owner
	err := post.CreateReactionNotification(actor.Id, "like")
	if err != nil {
		t.Fatalf("CreateReactionNotification() error: %v", err)
	}

	notifications, err := db.SelectNotificationsByUserId(owner.Id)
	if err != nil {
		t.Fatalf("GetNotificationsByUserId() error: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].Type != "like" {
		t.Errorf("Notification Type = %q, want %q", notifications[0].Type, "like")
	}
}

func TestCreateReactionNotificationComment(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)
	post.User = owner

	comment := CommentType{}
	comment.UserId = owner.Id
	comment.PostId = post.Id
	comment.Body = "react notif comment"
	comment.Add()

	err := comment.CreateReactionNotification(actor.Id, "commentLike")
	if err != nil {
		t.Fatalf("CreateReactionNotification() error: %v", err)
	}

	notifications, err := db.SelectNotificationsByUserId(owner.Id)
	if err != nil {
		t.Fatalf("GetNotificationsByUserId() error: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].Type != "commentLike" {
		t.Errorf("Notification Type = %q, want %q", notifications[0].Type, "commentLike")
	}
}

func TestUserGetNotifications(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	n := NotificationType{}
	n.UserId = owner.Id
	n.ActorId = actor.Id
	n.Type = "like"
	n.PostId = post.Id
	CreateNotification(n)

	err := owner.GetNotifications()
	if err != nil {
		t.Fatalf("GetNotifications() error: %v", err)
	}
	if len(owner.Notifications) != 1 {
		t.Errorf("GetNotifications() returned %d notifications, want 1", len(owner.Notifications))
	}
	if owner.UnreadNotificationsCount != 1 {
		t.Errorf("UnreadNotificationsCount = %d, want 1", owner.UnreadNotificationsCount)
	}
}

func TestUserGetNotificationsEmpty(t *testing.T) {
	owner, _, _ := setupTestNotificationDB(t)

	err := owner.GetNotifications()
	if err != nil {
		t.Fatalf("GetNotifications() error: %v", err)
	}
	if len(owner.Notifications) != 0 {
		t.Errorf("GetNotifications() returned %d notifications, want 0", len(owner.Notifications))
	}
}

func TestUserMarkNotificationAsRead(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	n := NotificationType{}
	n.UserId = owner.Id
	n.ActorId = actor.Id
	n.Type = "like"
	n.PostId = post.Id
	CreateNotification(n)
	owner.GetNotifications()

	if len(owner.Notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(owner.Notifications))
	}

	err := owner.MarkNotificationAsRead(owner.Notifications[0].Id)
	if err != nil {
		t.Fatalf("MarkNotificationAsRead() error: %v", err)
	}

	owner.Notifications = nil
	owner.UnreadNotificationsCount = 0
	owner.GetNotifications()
	if owner.UnreadNotificationsCount != 0 {
		t.Errorf("UnreadNotificationsCount after mark read = %d, want 0", owner.UnreadNotificationsCount)
	}
}

func TestUserMarkAllNotificationsAsRead(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	for i := 0; i < 3; i++ {
		n := NotificationType{}
		n.UserId = owner.Id
		n.ActorId = actor.Id
		n.Type = "like"
		n.PostId = post.Id
		CreateNotification(n)
	}

	owner.GetNotifications()
	if owner.UnreadNotificationsCount != 3 {
		t.Fatalf("UnreadNotificationsCount = %d, want 3", owner.UnreadNotificationsCount)
	}

	err := owner.MarkAllNotificationsAsRead()
	if err != nil {
		t.Fatalf("MarkAllNotificationsAsRead() error: %v", err)
	}

	owner.Notifications = nil
	owner.UnreadNotificationsCount = 0
	owner.GetNotifications()
	if owner.UnreadNotificationsCount != 0 {
		t.Errorf("UnreadNotificationsCount after mark all read = %d, want 0", owner.UnreadNotificationsCount)
	}
}

func TestMarkAsReadPost(t *testing.T) {
	owner, actor, post := setupTestNotificationDB(t)

	n := NotificationType{}
	n.UserId = owner.Id
	n.ActorId = actor.Id
	n.Type = "like"
	n.PostId = post.Id
	CreateNotification(n)
	owner.GetNotifications()

	err := owner.MarkAsReadPost(post)
	if err != nil {
		t.Fatalf("MarkAsReadPost() error: %v", err)
	}

	if owner.Notifications[0].Read != true {
		t.Error("MarkAsReadPost() should mark notification as read")
	}
	if owner.UnreadNotificationsCount != 0 {
		t.Errorf("UnreadNotificationsCount = %d, want 0", owner.UnreadNotificationsCount)
	}
}
