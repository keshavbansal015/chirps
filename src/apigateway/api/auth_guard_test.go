package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuthGuard_AllowedRoutes(t *testing.T) {
	e := echo.New()

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/sessions"},
		{http.MethodDelete, "/sessions"},
		{http.MethodPost, "/users"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			nextCalled := false
			next := func(c echo.Context) error {
				nextCalled = true
				return c.NoContent(http.StatusOK)
			}

			guard := authGuard(nil)
			err := guard(next)(c)

			if err != nil {
				t.Errorf("authGuard() error = %v", err)
			}
			if !nextCalled {
				t.Error("next handler was not called for allowed route")
			}
		})
	}
}
