package post

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name     string
		code     codes.Code
		wantCode codes.Code
	}{
		{
			name:     "internal error",
			code:     codes.Internal,
			wantCode: codes.Internal,
		},
		{
			name:     "not found error",
			code:     codes.NotFound,
			wantCode: codes.NotFound,
		},
		{
			name:     "unauthenticated error",
			code:     codes.Unauthenticated,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newError(tt.code)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			statusErr, ok := status.FromError(err)
			if !ok {
				t.Fatal("expected status error")
			}

			if statusErr.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, statusErr.Code())
			}

			if statusErr.Message() != tt.code.String() {
				t.Errorf("expected message %s, got %s", tt.code.String(), statusErr.Message())
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantID    int32
		wantError bool
		wantCode  codes.Code
	}{
		{
			name:      "valid user ID",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.MD{"user-id": []string{"123"}}),
			wantID:    123,
			wantError: false,
		},
		{
			name:      "missing metadata",
			ctx:       context.Background(),
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "empty user-id",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.MD{"user-id": []string{}}),
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "invalid user-id format",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.MD{"user-id": []string{"invalid"}}),
			wantError: true,
		},
		{
			name:      "negative user-id",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.MD{"user-id": []string{"-1"}}),
			wantID:    -1,
			wantError: false,
		},
		{
			name:      "zero user-id",
			ctx:       metadata.NewIncomingContext(context.Background(), metadata.MD{"user-id": []string{"0"}}),
			wantID:    0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := getUserID(tt.ctx)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.wantCode != 0 {
					statusErr, ok := status.FromError(err)
					if !ok {
						t.Error("expected status error")
						return
					}
					if statusErr.Code() != tt.wantCode {
						t.Errorf("expected code %v, got %v", tt.wantCode, statusErr.Code())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if id != tt.wantID {
				t.Errorf("expected ID %d, got %d", tt.wantID, id)
			}
		})
	}
}
