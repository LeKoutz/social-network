package router

import (
	"forum/src/db"
	"forum/src/models"
	"forum/src/state"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouterDB(t *testing.T) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
}

func newState(req *http.Request) *state.State {
	w := httptest.NewRecorder()
	s := &state.State{}
	s.Init().SetResponse(w).SetRequest(req).SetUser(models.GetGuestUser())
	return s
}

func TestMatchRouteIndex(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /")
	}
	if route.Path != "/" {
		t.Errorf("matchRoute() Path = %q, want %q", route.Path, "/")
	}
}

func TestMatchRoutePosts(t *testing.T) {
	req := httptest.NewRequest("GET", "/posts", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /posts")
	}
	if route.Path != "/posts" {
		t.Errorf("matchRoute() Path = %q, want %q", route.Path, "/posts")
	}
}

func TestMatchRoutePostView(t *testing.T) {
	req := httptest.NewRequest("GET", "/post/view/1", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /post/view/1")
	}
	if !route.Prefix {
		t.Error("matchRoute() route should have Prefix=true for /post/view/")
	}
}

func TestMatchRoutePostCreate(t *testing.T) {
	req := httptest.NewRequest("GET", "/post/create", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /post/create")
	}
	if !route.NeedsLogin {
		t.Error("matchRoute() /post/create should require login")
	}
}

func TestMatchRouteCategoryView(t *testing.T) {
	req := httptest.NewRequest("GET", "/category/view/1", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /category/view/1")
	}
	if !route.Prefix {
		t.Error("matchRoute() route should have Prefix=true for /category/view/")
	}
}

func TestMatchRouteLogin(t *testing.T) {
	req := httptest.NewRequest("GET", "/user/login", nil)
	s := newState(req)

	route, err := matchRoute(s)
	if err != nil {
		t.Fatalf("matchRoute() error: %v", err)
	}
	if route == nil {
		t.Fatal("matchRoute() returned nil route for /user/login")
	}
	if route.NeedsLogin {
		t.Error("matchRoute() /user/login should not require login")
	}
}

func TestMatchRouteNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/nonexistent/path", nil)
	s := newState(req)

	_, err := matchRoute(s)
	if err == nil {
		t.Error("matchRoute() should return error for unknown path")
	}
}

func TestMatchRouteMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/posts", nil)
	s := newState(req)

	_, err := matchRoute(s)
	if err == nil {
		t.Error("matchRoute() should return error for wrong method")
	}
}

func TestRoutesHandlerIndex(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandler404(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/does/not/exist", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("RoutesHandler() status for 404 = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestRoutesHandlerMethodNotAllowed(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("DELETE", "/posts", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("RoutesHandler() status for wrong method = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestRoutesHandlerPostsPage(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() status for /posts = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandlerCategoriesPage(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() status for /categories = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandlerLoginPage(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/user/login", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() status for /user/login = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandlerRegisterPage(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/user/register", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() status for /user/register = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandlerNeedsLogin(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/user/posts", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("RoutesHandler() status for /user/posts without login = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestRoutesHandlerPostViewNonexistent(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/post/view/999", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("RoutesHandler() status for nonexistent post = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestRoutesHandlerWithInvalidCookie(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "__Host-FRMSessionID",
		Value: "nonexistent-session",
	})
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("RoutesHandler() with invalid cookie status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRoutesHandlerCommentRoutes(t *testing.T) {
	setupTestRouterDB(t)

	routes := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/comment/create", http.StatusMethodNotAllowed},
		{"GET", "/comment/react", http.StatusMethodNotAllowed},
		{"GET", "/comment/edit", http.StatusMethodNotAllowed},
		{"GET", "/comment/delete", http.StatusMethodNotAllowed},
	}

	for _, tt := range routes {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			RoutesHandler(w, req)
			if w.Code != tt.code {
				t.Errorf("RoutesHandler(%s %s) status = %d, want %d", tt.method, tt.path, w.Code, tt.code)
			}
		})
	}
}

func TestRoutesHandlerPostReaction(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/post/react", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("RoutesHandler() status for GET /post/react = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestRoutesHandlerUserActivity(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/user/activity", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("RoutesHandler() status for /user/activity without login = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestRoutesHandlerUserLogout(t *testing.T) {
	setupTestRouterDB(t)

	req := httptest.NewRequest("GET", "/user/logout", nil)
	w := httptest.NewRecorder()

	RoutesHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("RoutesHandler() status for /user/logout without login = %d, want %d", w.Code, http.StatusForbidden)
	}
}
