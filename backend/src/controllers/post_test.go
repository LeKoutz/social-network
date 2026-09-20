package controllers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"forum/src/db"
	"forum/src/ferror"
	"forum/src/models"
	"forum/src/state"
	"forum/src/utils"
)

func newActionRequest(t *testing.T, action string) *http.Request {
	t.Helper()
	form := url.Values{}
	form.Set("action", action)
	req, err := http.NewRequest("POST", "/api/post/react", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := req.ParseForm(); err != nil {
		t.Fatalf("Failed to parse form: %v", err)
	}
	return req
}

func TestPostReactionGroupAccess(t *testing.T) {
	if err := db.InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}
	hash, _ := utils.HashPassword("password123")
	owner := models.UserType{}
	owner.Username = "reactionowner"
	owner.Email = "reactionowner@test.com"
	owner.Hash = hash
	if err := owner.Add(); err != nil {
		t.Fatalf("Failed to create owner: %v", err)
	}
	stranger := models.UserType{}
	stranger.Username = "reactionstranger"
	stranger.Email = "reactionstranger@test.com"
	stranger.Hash = hash
	if err := stranger.Add(); err != nil {
		t.Fatalf("Failed to create stranger: %v", err)
	}

	group := models.GroupType{}
	group.Title = "Reaction Club"
	group.Description = "Members react here"
	group.OwnerUserId = owner.Id
	if err := group.Add(); err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	post := models.PostType{}
	post.Title = "Group post"
	post.Body = "Members only"
	post.UserId = owner.Id
	post.GroupId = group.Id
	if err := post.Add(); err != nil {
		t.Fatalf("Failed to create group post: %v", err)
	}

	t.Run("member can react", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(owner)
		p := models.PostType{}
		p.Id = post.Id
		s.SetPost(p)
		s.SetRequest(newActionRequest(t, "like"))

		err := PostReaction(s)
		if err != nil {
			t.Fatalf("PostReaction() member error: %v", err)
		}
	})

	t.Run("non-member cannot react", func(t *testing.T) {
		s := &state.State{}
		s.Init()
		s.SetUser(stranger)
		p := models.PostType{}
		p.Id = post.Id
		s.SetPost(p)
		s.SetRequest(newActionRequest(t, "like"))

		err := PostReaction(s)
		if err != ferror.ErrorGroupMembershipRequired {
			t.Errorf("PostReaction() stranger error = %v, want %v", err, ferror.ErrorGroupMembershipRequired)
		}
	})
}