package ferror

import (
	"net/http"
	"testing"
)

func TestErrorConsume(t *testing.T) {
	tests := []struct {
		name       string
		input      error
		wantCode   int
		wantHas    bool
		wantMsgNil bool
	}{
		{"not found", ErrorNotFound, http.StatusNotFound, true, false},
		{"content not found", ErrorContentNotFound, http.StatusNotFound, true, false},
		{"unauthorized", ErrorUnauthorizedAction, http.StatusForbidden, true, false},
		{"user permission", ErrorUserPermissionDenied, http.StatusForbidden, true, false},
		{"comment permission", ErrorCommentPermissionDenied, http.StatusForbidden, true, false},
		{"post permission", ErrorPostPermissionDenied, http.StatusForbidden, true, false},
		{"method not allowed", ErrorMethodNotAllowed, http.StatusMethodNotAllowed, true, false},
		{"post empty id", ErrorPostEmptyId, http.StatusBadRequest, true, false},
		{"invalid post id", ErrorInvalidPostId, http.StatusBadRequest, true, false},
		{"invalid comment id", ErrorInvalidCommentId, http.StatusBadRequest, true, false},
		{"invalid category id", ErrorInvalidCategoryId, http.StatusBadRequest, true, false},
		{"post body empty", ErrorPostBodyEmpty, http.StatusBadRequest, true, false},
		{"post title empty", ErrorPostTitleEmpty, http.StatusBadRequest, true, false},
		{"post no category", ErrorPostHasNoCategory, http.StatusBadRequest, true, false},
		{"comment empty", ErrorCommentEmpty, http.StatusBadRequest, true, false},
		{"comment too long", ErrorCommentTooLong, http.StatusBadRequest, true, false},
		{"comment empty id", ErrorCommentEmptyId, http.StatusBadRequest, true, false},
		{"category empty id", ErrorCategoryEmptyId, http.StatusBadRequest, true, false},
		{"category name empty", ErrorCategoryNameEmpty, http.StatusBadRequest, true, false},
		{"category name too long", ErrorCategoryNameTooLong, http.StatusBadRequest, true, false},
		{"email field empty", ErrorEmailFieldEmpty, http.StatusBadRequest, true, false},
		{"password field empty", ErrorPasswordFieldEmpty, http.StatusBadRequest, true, false},
		{"bad request", ErrorBadRequest, http.StatusBadRequest, true, false},
		{"internal server error", ErrorInternalServerError, http.StatusInternalServerError, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{}
			result := e.Consume(tt.input)
			if result != e {
				t.Error("Consume() should return same Error pointer")
			}
			if e.StatusCode != tt.wantCode {
				t.Errorf("Consume() StatusCode = %d, want %d", e.StatusCode, tt.wantCode)
			}
			if e.Has != tt.wantHas {
				t.Errorf("Consume() Has = %v, want %v", e.Has, tt.wantHas)
			}
			if len(e.Message) == 0 {
				t.Error("Consume() Message should not be empty")
			}
			if e.Error == nil {
				t.Error("Consume() Error should not be nil")
			}
		})
	}
}

func TestErrorConsumeUnmapped(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorNotRegistered)
	if e.StatusCode != 0 {
		t.Errorf("Consume() for unmapped error StatusCode = %d, want 0", e.StatusCode)
	}
	if !e.Has {
		t.Error("Consume() for unmapped error Has should be true")
	}
}

func TestErrorConsumeErrorMessage(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorNotFound)
	if e.Message != "Not found" {
		t.Errorf("Consume() Message = %q, want %q", e.Message, "Not found")
	}
}

func TestErrorConsumeWrappedError(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorWrongPassword)
	if e.Message != "Wrong password" {
		t.Errorf("Consume() Message = %q, want %q", e.Message, "Wrong password")
	}
	if e.StatusCode != 0 {
		t.Errorf("Consume() StatusCode for unmapped error = %d, want 0", e.StatusCode)
	}
}

func TestErrorConsumeNewlineMessage(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorInternalServerError)
	if e.StatusCode != http.StatusInternalServerError {
		t.Errorf("Consume() StatusCode = %d, want %d", e.StatusCode, http.StatusInternalServerError)
	}
}

func TestErrorConsumeMultipleCalls(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorNotFound)
	if e.StatusCode != http.StatusNotFound {
		t.Errorf("First Consume() StatusCode = %d, want %d", e.StatusCode, http.StatusNotFound)
	}
	e.Consume(ErrorBadRequest)
	if e.StatusCode != http.StatusBadRequest {
		t.Errorf("Second Consume() StatusCode = %d, want %d", e.StatusCode, http.StatusBadRequest)
	}
}

func TestErrorLogError(t *testing.T) {
	e := &Error{}
	e.Consume(ErrorNotFound)
	e.LogError()
	if e.Message == "" {
		t.Error("LogError() called but Message is empty")
	}
}
