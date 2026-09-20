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
