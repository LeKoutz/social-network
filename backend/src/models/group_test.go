package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"strings"
	"testing"
)

func setupTestGroupDB(t *testing.T) UserType {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := UserType{}
	user.Username = "groupowner"
	user.Email = "group@test.com"
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestGroupValidateGroup(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"valid title", "Go Developers", nil},
		{"empty title", "", ferror.ErrorGroupTitleEmpty},
		{"max length", strings.Repeat("a", 127), nil},
		{"too long title", strings.Repeat("a", 128), ferror.ErrorGroupTitleTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GroupType{}
			g.Title = tt.title
			err := g.ValidateGroup()
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateGroup() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateGroup() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGroupAdd(t *testing.T) {
	user := setupTestGroupDB(t)

	g := &GroupType{}
	g.Title = "Go Developers"
	g.Description = "For Go enthusiasts"
	g.OwnerUserId = user.Id
	err := g.Add()
	if err != nil {
		t.Fatalf("Add() error: %v", err)
	}
	if g.Id == 0 {
		t.Error("Add() did not set group ID")
	}

	g2 := &GroupType{}
	g2.Title = ""
	g2.OwnerUserId = user.Id
	err = g2.Add()
	if err == nil {
		t.Error("Add() should fail with empty title")
	}

	g3 := &GroupType{}
	g3.Title = "Go Developers"
	g3.OwnerUserId = user.Id
	err = g3.Add()
	if err == nil {
		t.Error("Add() should fail with duplicate title")
	} else if err != ferror.ErrorGroupTitleAlreadyExists {
		t.Errorf("Add() duplicate title error = %v, want %v", err, ferror.ErrorGroupTitleAlreadyExists)
	}
}

func TestGroupGetGroups(t *testing.T) {
	user := setupTestGroupDB(t)

	for _, title := range []string{"Group One", "Group Two"} {
		g := &GroupType{}
		g.Title = title
		g.OwnerUserId = user.Id
		if err := g.Add(); err != nil {
			t.Fatalf("Add() error for %q: %v", title, err)
		}
	}

	var groups GroupsType
	if err := groups.GetGroups(); err != nil {
		t.Fatalf("GetGroups() error: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("GetGroups() len = %d, want 2", len(groups))
	}
	found := map[string]bool{}
	for _, group := range groups {
		found[group.Title] = true
		if group.OwnerUsername != user.Username {
			t.Errorf("GetGroups() OwnerUsername = %q, want %q", group.OwnerUsername, user.Username)
		}
	}
	for _, title := range []string{"Group One", "Group Two"} {
		if !found[title] {
			t.Errorf("GetGroups() missing group %q", title)
		}
	}
}

func TestGroupGetById(t *testing.T) {
	user := setupTestGroupDB(t)

	g := &GroupType{}
	g.Title = "Club Alpha"
	g.Description = "Members only"
	g.OwnerUserId = user.Id
	if err := g.Add(); err != nil {
		t.Fatalf("Add() error: %v", err)
	}

	got := &GroupType{}
	got.Id = g.Id
	if err := got.GetById(); err != nil {
		t.Fatalf("GetById() error: %v", err)
	}
	if got.Title != g.Title || got.Description != g.Description {
		t.Errorf("GetById() = %+v, want title %q description %q", got, g.Title, g.Description)
	}
	if got.OwnerUsername != user.Username {
		t.Errorf("GetById() OwnerUsername = %q, want %q", got.OwnerUsername, user.Username)
	}
}

func TestGroupIsMember(t *testing.T) {
	user := setupTestGroupDB(t)

	g := &GroupType{}
	g.Title = "Closed Club"
	g.OwnerUserId = user.Id
	if err := g.Add(); err != nil {
		t.Fatalf("Add() error: %v", err)
	}

	member, err := g.IsMember(user.Id)
	if err != nil {
		t.Fatalf("IsMember() error: %v", err)
	}
	if !member {
		t.Error("IsMember(owner) = false, want true")
	}

	stranger := &UserType{}
	stranger.Username = "stranger"
	stranger.Email = "stranger@test.com"
	stranger.Hash, _ = utils.HashPassword("password123")
	if err := stranger.Add(); err != nil {
		t.Fatalf("failed to create stranger: %v", err)
	}
	member, err = g.IsMember(stranger.Id)
	if err != nil {
		t.Fatalf("IsMember() error: %v", err)
	}
	if member {
		t.Error("IsMember(stranger) = true, want false")
	}
}

func TestGroupGetPostsByGroupId(t *testing.T) {
	user := setupTestGroupDB(t)

	g := &GroupType{}
	g.Title = "Group Hub"
	g.OwnerUserId = user.Id
	if err := g.Add(); err != nil {
		t.Fatalf("Add() group error: %v", err)
	}

	groupPost := &PostType{}
	groupPost.Title = "Group only post"
	groupPost.Body = "Secret content"
	groupPost.UserId = user.Id
	groupPost.GroupId = g.Id
	if err := groupPost.Add(); err != nil {
		t.Fatalf("Add() group post error: %v", err)
	}

	cat := CategoryType{}
	cat.Name = "general"
	cat.Description = "General discussion"
	if err := cat.Add(); err != nil {
		t.Fatalf("Add() category error: %v", err)
	}
	publicPost := &PostType{}
	publicPost.Title = "Public post"
	publicPost.Body = "Everyone sees this"
	publicPost.UserId = user.Id
	publicPost.Categories = CategoriesType{cat}
	if err := publicPost.Add(); err != nil {
		t.Fatalf("Add() public post error: %v", err)
	}

	var groupPosts PostsType
	if err := groupPosts.GetPostsByGroupId(g.Id); err != nil {
		t.Fatalf("GetPostsByGroupId() error: %v", err)
	}
	if len(groupPosts) != 1 {
		t.Errorf("GetPostsByGroupId() len = %d, want 1", len(groupPosts))
	}
	if len(groupPosts) > 0 {
		if groupPosts[0].GroupId != g.Id {
			t.Errorf("group post GroupId = %d, want %d", groupPosts[0].GroupId, g.Id)
		}
		if groupPosts[0].GroupTitle != g.Title {
			t.Errorf("group post GroupTitle = %q, want %q", groupPosts[0].GroupTitle, g.Title)
		}
	}

	var all PostsType
	if err := all.GetPosts(); err != nil {
		t.Fatalf("GetPosts() error: %v", err)
	}
	for _, post := range all {
		if post.GroupId != 0 {
			t.Errorf("public feed leaked group post id %d", post.Id)
		}
	}
}
