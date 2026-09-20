package controllers

import (
	"testing"

	"forum/src/db"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func TestCreateGroup(t *testing.T) {
	if err := db.InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := models.UserType{}
	user.Username = "groupcreator"
	user.Email = "groupcreator@test.com"
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("valid group", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(user)
		g := models.GroupType{}
		g.Title = "Test Group"
		g.Description = "Test description"
		g.OwnerUserId = user.Id
		s.SetGroup(g)

		err := CreateGroup(s)
		if err != nil {
			t.Fatalf("CreateGroup() error: %v", err)
		}
		if s.GetGroup().Id == 0 {
			t.Error("CreateGroup() did not set group ID")
		}
	})

	t.Run("empty title", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(user)
		g := models.GroupType{}
		g.OwnerUserId = user.Id
		s.SetGroup(g)

		err := CreateGroup(s)
		if err != ferror.ErrorGroupTitleEmpty {
			t.Errorf("CreateGroup() error = %v, want %v", err, ferror.ErrorGroupTitleEmpty)
		}
	})
}

func TestShowGroups(t *testing.T) {
	if err := db.InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}
	hash, _ := utils.HashPassword("password123")
	user := models.UserType{}
	user.Username = "grouplister"
	user.Email = "grouplister@test.com"
	user.Hash = hash
	if err := user.Add(); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	for _, title := range []string{"Alpha Group", "Beta Group"} {
		g := models.GroupType{}
		g.Title = title
		g.OwnerUserId = user.Id
		if err := g.Add(); err != nil {
			t.Fatalf("Failed to create group %q: %v", title, err)
		}
	}

	s := &state.State{}
	s.Init()
	s.SetUser(user)
	if err := ShowGroups(s); err != nil {
		t.Fatalf("ShowGroups() error: %v", err)
	}
	if len(s.GetGroups()) != 2 {
		t.Errorf("ShowGroups() len = %d, want 2", len(s.GetGroups()))
	}
}

func TestShowGroup(t *testing.T) {
	if err := db.InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}
	hash, _ := utils.HashPassword("password123")
	owner := models.UserType{}
	owner.Username = "groupowner"
	owner.Email = "groupowner@test.com"
	owner.Hash = hash
	if err := owner.Add(); err != nil {
		t.Fatalf("Failed to create owner: %v", err)
	}
	stranger := models.UserType{}
	stranger.Username = "stranger"
	stranger.Email = "stranger@test.com"
	stranger.Hash = hash
	if err := stranger.Add(); err != nil {
		t.Fatalf("Failed to create stranger: %v", err)
	}

	group := models.GroupType{}
	group.Title = "Members Area"
	group.Description = "Only members"
	group.OwnerUserId = owner.Id
	if err := group.Add(); err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	post := models.PostType{}
	post.Title = "Group update"
	post.Body = "Members only content"
	post.UserId = owner.Id
	post.GroupId = group.Id
	if err := post.Add(); err != nil {
		t.Fatalf("Failed to create group post: %v", err)
	}

	t.Run("member sees posts", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(owner)
		g := models.GroupType{}
		g.Id = group.Id
		s.SetGroup(g)

		err := ShowGroup(s)
		if err != nil {
			t.Fatalf("ShowGroup() error: %v", err)
		}
		if !s.GetGroup().Member {
			t.Error("ShowGroup() member = false for owner, want true")
		}
		if len(s.GetGroup().Posts) != 1 {
			t.Errorf("ShowGroup() group posts len = %d, want 1", len(s.GetGroup().Posts))
		}
	})

	t.Run("non-member sees no posts", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(stranger)
		g := models.GroupType{}
		g.Id = group.Id
		s.SetGroup(g)

		err := ShowGroup(s)
		if err != nil {
			t.Fatalf("ShowGroup() error: %v", err)
		}
		if s.GetGroup().Member {
			t.Error("ShowGroup() member = true for stranger, want false")
		}
		if len(s.GetGroup().Posts) != 0 {
			t.Errorf("ShowGroup() leaked %d posts to non-member", len(s.GetGroup().Posts))
		}
	})
}
