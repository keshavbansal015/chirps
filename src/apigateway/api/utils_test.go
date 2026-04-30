package api

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		code     codes.Code
		expected int
	}{
		{"InvalidArgument", codes.InvalidArgument, http.StatusBadRequest},
		{"Unauthenticated", codes.Unauthenticated, http.StatusUnauthorized},
		{"PermissionDenied", codes.PermissionDenied, http.StatusForbidden},
		{"NotFound", codes.NotFound, http.StatusNotFound},
		{"Internal", codes.Internal, http.StatusInternalServerError},
		{"Unknown", codes.Unknown, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := status.New(tt.code, "test error")
			got := getStatusCode(s)
			if got != tt.expected {
				t.Errorf("getStatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		grpcCode    codes.Code
		grpcMessage string
		wantCode    int
		wantMessage string
	}{
		{"BadRequest", codes.InvalidArgument, "invalid input", http.StatusBadRequest, "invalid input"},
		{"Unauthorized", codes.Unauthenticated, "unauthorized", http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", codes.PermissionDenied, "forbidden", http.StatusForbidden, "forbidden"},
		{"NotFound", codes.NotFound, "not found", http.StatusNotFound, "not found"},
		{"InternalError", codes.Internal, "server error", http.StatusInternalServerError, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := status.Error(tt.grpcCode, tt.grpcMessage)
			httpErr := newHTTPError(err)

			if httpErr.Code != tt.wantCode {
				t.Errorf("newHTTPError() code = %v, want %v", httpErr.Code, tt.wantCode)
			}
			if httpErr.Message != tt.wantMessage {
				t.Errorf("newHTTPError() message = %v, want %v", httpErr.Message, tt.wantMessage)
			}
		})
	}
}

func TestNewHTTPError_NilError(t *testing.T) {
	httpErr := newHTTPError(errors.New("generic error"))
	if httpErr.Code != http.StatusInternalServerError {
		t.Errorf("newHTTPError() with generic error code = %v, want %v", httpErr.Code, http.StatusInternalServerError)
	}
}

func TestInsecureCredentials(t *testing.T) {
	opt := insecureCredentials()
	if opt == nil {
		t.Error("insecureCredentials() returned nil")
	}
}
