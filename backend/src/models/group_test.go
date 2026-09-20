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