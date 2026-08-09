package state

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStateInit(t *testing.T) {
	s := &State{}
	result := s.Init()
	if result != s {
		t.Error("Init() should return same State pointer")
	}
	if s.WebsiteName != "Forum" {
		t.Errorf("Init() WebsiteName = %q, want %q", s.WebsiteName, "Forum")
	}
	if s.User.Username != "guest" {
		t.Errorf("Init() User.Username = %q, want %q", s.User.Username, "guest")
	}
}

func TestStateSetView(t *testing.T) {
	s := &State{}
	result := s.SetView("test_view")
	if result != s {
		t.Error("SetView() should return same State pointer")
	}
	if s.View != "test_view" {
		t.Errorf("SetView() View = %q, want %q", s.View, "test_view")
	}
}

func TestStateSetError(t *testing.T) {
	s := &State{}
	result := s.SetError(ferror.Error{Message: "test error"})
	if result != s {
		t.Error("SetError() should return same State pointer")
	}
	if s.Error.Message != "test error" {
		t.Errorf("SetError() Message = %q, want %q", s.Error.Message, "test error")
	}
}

func TestStateSetErrorConsume(t *testing.T) {
	s := &State{}
	result := s.SetErrorConsume(ferror.ErrorNotFound)
	if result != s {
		t.Error("SetErrorConsume() should return same State pointer")
	}
	if !s.Error.Has {
		t.Error("SetErrorConsume() should set Error.Has")
	}
	if s.Error.StatusCode != http.StatusNotFound {
		t.Errorf("SetErrorConsume() StatusCode = %d, want %d", s.Error.StatusCode, http.StatusNotFound)
	}
}

func TestStateSetRequest(t *testing.T) {
	s := &State{}
	req := httptest.NewRequest("GET", "/test", nil)
	result := s.SetRequest(req)
	if result != s {
		t.Error("SetRequest() should return same State pointer")
	}
	if s.Request != req {
		t.Error("SetRequest() did not set Request correctly")
	}
}

func TestStateSetResponse(t *testing.T) {
	s := &State{}
	w := httptest.NewRecorder()
	result := s.SetResponse(w)
	if result != s {
		t.Error("SetResponse() should return same State pointer")
	}
}

func TestStateGetRequest(t *testing.T) {
	s := &State{}
	req := httptest.NewRequest("GET", "/test", nil)
	s.Request = req
	if s.GetRequest() != req {
		t.Error("GetRequest() did not return correct request")
	}
}

func TestStateGetError(t *testing.T) {
	s := &State{}
	errPtr := s.GetError()
	if errPtr == nil {
		t.Error("GetError() returned nil")
	}
}

func TestStateGetRedirect(t *testing.T) {
	s := &State{}
	s.Redirect = "/test/redirect"
	if s.GetRedirect() != "/test/redirect" {
		t.Errorf("GetRedirect() = %q, want %q", s.GetRedirect(), "/test/redirect")
	}
}

func TestStateSetRedirect(t *testing.T) {
	s := &State{}
	result := s.SetRedirect("/new/redirect")
	if result != s {
		t.Error("SetRedirect() should return same State pointer")
	}
	if s.Redirect != "/new/redirect" {
		t.Errorf("SetRedirect() Redirect = %q, want %q", s.Redirect, "/new/redirect")
	}
}

func TestStateSetMessage(t *testing.T) {
	s := &State{}
	msg := models.Message{Has: true, Type: "Success", Content: "test"}
	result := s.SetMessage(msg)
	if result != s {
		t.Error("SetMessage() should return same State pointer")
	}
	if s.Message.Content != "test" {
		t.Errorf("SetMessage() Content = %q, want %q", s.Message.Content, "test")
	}
}

func TestStateSetWebsiteName(t *testing.T) {
	s := &State{}
	result := s.SetWebsiteName("My Forum")
	if result != s {
		t.Error("SetWebsiteName() should return same State pointer")
	}
	if s.WebsiteName != "My Forum" {
		t.Errorf("SetWebsiteName() WebsiteName = %q, want %q", s.WebsiteName, "My Forum")
	}
}

func TestStateReturnController(t *testing.T) {
	s := &State{}
	sc := s.ReturnController()
	if sc == nil {
		t.Error("ReturnController() returned nil")
	}
}

func TestStateReturnHandler(t *testing.T) {
	s := &State{}
	sh := s.ReturnHandler()
	if sh == nil {
		t.Error("ReturnHandler() returned nil")
	}
}

func TestStateReturnView(t *testing.T) {
	s := &State{}
	sv := s.ReturnView()
	if sv == nil {
		t.Error("ReturnView() returned nil")
	}
}

func TestStateEditUser(t *testing.T) {
	s := &State{}
	s.User = models.GetGuestUser()
	u := s.EditUser()
	if u == nil {
		t.Error("EditUser() returned nil")
	}
	u.Username = "modified"
	if s.User.Username != "modified" {
		t.Error("EditUser() should modify original user")
	}
}

func TestStateGetUser(t *testing.T) {
	s := &State{}
	s.User = models.GetGuestUser()
	u := s.GetUser()
	if u.Username != "guest" {
		t.Errorf("GetUser() Username = %q, want %q", u.Username, "guest")
	}
}

func TestStateSetUser(t *testing.T) {
	s := &State{}
	user := models.UserType{}
	user.Username = "newuser"
	result := s.SetUser(user)
	if result != s {
		t.Error("SetUser() should return same State pointer")
	}
	if s.User.Username != "newuser" {
		t.Errorf("SetUser() Username = %q, want %q", s.User.Username, "newuser")
	}
}

func TestStateEditPost(t *testing.T) {
	s := &State{}
	p := s.EditPost()
	if p == nil {
		t.Error("EditPost() returned nil")
	}
}

func TestStateGetPost(t *testing.T) {
	s := &State{}
	p := s.GetPost()
	if p.Id != 0 {
		t.Errorf("GetPost() Id = %d, want 0", p.Id)
	}
}

func TestStateSetPost(t *testing.T) {
	s := &State{}
	post := models.PostType{}
	post.Title = "test"
	result := s.SetPost(post)
	if result != s {
		t.Error("SetPost() should return same State pointer")
	}
	if s.Posts[0].Title != "test" {
		t.Errorf("SetPost() Title = %q, want %q", s.Posts[0].Title, "test")
	}
}

func TestStateEditPosts(t *testing.T) {
	s := &State{}
	posts := s.EditPosts()
	if posts == nil {
		t.Error("EditPosts() returned nil")
	}
}

func TestStateSetPosts(t *testing.T) {
	s := &State{}
	posts := models.PostsType{
		{PostRowType: db.PostRowType{Title: "Post B", TimestampString: "1700000000"}},
		{PostRowType: db.PostRowType{Title: "Post A", TimestampString: "1700000010"}},
	}
	result := s.SetPosts(posts)
	if result != s {
		t.Error("SetPosts() should return same State pointer")
	}
	if len(s.Posts) != 2 {
		t.Errorf("SetPosts() len = %d, want 2", len(s.Posts))
	}
}

func TestStateEditComment(t *testing.T) {
	s := &State{}
	c := s.EditComment()
	if c == nil {
		t.Error("EditComment() returned nil")
	}
}

func TestStateGetComment(t *testing.T) {
	s := &State{}
	c := s.GetComment()
	if c.Id != 0 {
		t.Errorf("GetComment() Id = %d, want 0", c.Id)
	}
}

func TestStateSetEditCommentId(t *testing.T) {
	s := &State{}
	s.SetEditCommentId(42)
	if s.EditCommentId != 42 {
		t.Errorf("SetEditCommentId() EditCommentId = %d, want 42", s.EditCommentId)
	}
}

func TestStateEditCategory(t *testing.T) {
	s := &State{}
	c := s.EditCategory()
	if c == nil {
		t.Error("EditCategory() returned nil")
	}
}

func TestStateGetCategory(t *testing.T) {
	s := &State{}
	c := s.GetCategory()
	if c.Id != 0 {
		t.Errorf("GetCategory() Id = %d, want 0", c.Id)
	}
}

func TestStateSetCategory(t *testing.T) {
	s := &State{}
	cat := models.CategoryType{}
	cat.Name = "testcat"
	result := s.SetCategory(cat)
	if result != s {
		t.Error("SetCategory() should return same State pointer")
	}
	if s.Categories[0].Name != "testcat" {
		t.Errorf("SetCategory() Name = %q, want %q", s.Categories[0].Name, "testcat")
	}
}

func TestStateEditCategories(t *testing.T) {
	s := &State{}
	cats := s.EditCategories()
	if cats == nil {
		t.Error("EditCategories() returned nil")
	}
}

func TestStateGetCategories(t *testing.T) {
	s := &State{}
	cats := s.GetCategories()
	if len(cats) != 0 {
		t.Errorf("GetCategories() should return empty slice, got %d", len(cats))
	}
}

func TestStateSetCategories(t *testing.T) {
	s := &State{}
	c := models.CategoryType{}
	c.Name = "cat1"
	cats := models.CategoriesType{c}
	result := s.SetCategories(cats)
	if result != s {
		t.Error("SetCategories() should return same State pointer")
	}
	if len(s.Categories) != 1 {
		t.Errorf("SetCategories() len = %d, want 1", len(s.Categories))
	}
}

func TestStateSetEditPost(t *testing.T) {
	s := &State{}
	s.SetEditPost(true)
	if s.EditingPost != true {
		t.Error("SetEditPost(true) did not set EditingPost")
	}
	s.SetEditPost(false)
	if s.EditingPost != false {
		t.Error("SetEditPost(false) did not set EditingPost")
	}
}

func TestStateEditResponse(t *testing.T) {
	s := &State{}
	w := httptest.NewRecorder()
	s.Response = w
	resp := s.EditResponse()
	if resp == nil {
		t.Error("EditResponse() returned nil")
	}
}

func TestStateGetUserLoggedIn(t *testing.T) {
	s := &State{}
	s.User = models.GetGuestUser()
	s.User.LoggedIn = true
	if !s.GetUser().LoggedIn {
		t.Error("GetUser().LoggedIn should be true")
	}
}

func TestGetGuestUser(t *testing.T) {
	u := models.GetGuestUser()
	if u.Username != "guest" {
		t.Errorf("GetGuestUser().Username = %q, want %q", u.Username, "guest")
	}
	if u.LoggedIn {
		t.Error("GetGuestUser().LoggedIn should be false")
	}
}
