package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestGetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userId", "123")

	result := getUserID(c)
	if result != "123" {
		t.Errorf("getUserID() = %v, want 123", result)
	}
}

func TestSetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setUserID(c, "456")

	result := c.Get("userId")
	if result != "456" {
		t.Errorf("setUserID() set %v, want 456", result)
	}
}

func TestCreateCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	createCookie(c, "session-123")

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "session" {
		t.Errorf("cookie.Name = %v, want session", cookie.Name)
	}
	if cookie.Value != "session-123" {
		t.Errorf("cookie.Value = %v, want session-123", cookie.Value)
	}
	if !cookie.HttpOnly {
		t.Error("cookie.HttpOnly should be true")
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("cookie.SameSite = %v, want SameSiteStrictMode", cookie.SameSite)
	}
	if cookie.Expires.IsZero() {
		t.Error("cookie.Expires should not be zero")
	}

	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	if cookie.Expires.Before(expectedExpiry.Add(-time.Hour)) || cookie.Expires.After(expectedExpiry.Add(time.Hour)) {
		t.Errorf("cookie.Expires is not within expected range")
	}
}

func TestClearCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	clearCookie(c)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "session" {
		t.Errorf("cookie.Name = %v, want session", cookie.Name)
	}
	if cookie.Value != "" {
		t.Errorf("cookie.Value = %v, want empty string", cookie.Value)
	}
	if cookie.MaxAge != 0 {
		t.Errorf("cookie.MaxAge = %v, want 0", cookie.MaxAge)
	}
}
