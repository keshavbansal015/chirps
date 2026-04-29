package post

import (
	"context"
	"errors"
	"testing"

	pb "github.com/keshavbansal015/chirps/src/postservice/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// mockDbClient implements DbClientInterface for testing
type mockDbClient struct {
	createPostFunc     func(content string, userID int32) (int32, error)
	getFeedFunc        func(page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
	getPostsFunc       func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
	getLikedPostsFunc  func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
	getPostFunc        func(id int32, currentUserID int32) (*pb.Post, error)
	deletePostFunc     func(postID int32, userID int32) error
	likePostFunc       func(postID int32, userID int32) error
	unlikePostFunc     func(postID int32, userID int32) error
	repostPostFunc     func(postID int32, userID int32) error
	removeRepostFunc   func(postID int32, userID int32) error
}

func (m *mockDbClient) createPost(content string, userID int32) (int32, error) {
	if m.createPostFunc != nil {
		return m.createPostFunc(content, userID)
	}
	return 0, nil
}

func (m *mockDbClient) getFeed(page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
	if m.getFeedFunc != nil {
		return m.getFeedFunc(page, limit, currentUserID)
	}
	return nil, nil
}

func (m *mockDbClient) getPosts(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
	if m.getPostsFunc != nil {
		return m.getPostsFunc(userID, page, limit, currentUserID)
	}
	return nil, nil
}

func (m *mockDbClient) getLikedPosts(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
	if m.getLikedPostsFunc != nil {
		return m.getLikedPostsFunc(userID, page, limit, currentUserID)
	}
	return nil, nil
}

func (m *mockDbClient) getPost(id int32, currentUserID int32) (*pb.Post, error) {
	if m.getPostFunc != nil {
		return m.getPostFunc(id, currentUserID)
	}
	return nil, nil
}

func (m *mockDbClient) deletePost(postID int32, userID int32) error {
	if m.deletePostFunc != nil {
		return m.deletePostFunc(postID, userID)
	}
	return nil
}

func (m *mockDbClient) likePost(postID int32, userID int32) error {
	if m.likePostFunc != nil {
		return m.likePostFunc(postID, userID)
	}
	return nil
}

func (m *mockDbClient) unlikePost(postID int32, userID int32) error {
	if m.unlikePostFunc != nil {
		return m.unlikePostFunc(postID, userID)
	}
	return nil
}

func (m *mockDbClient) repostPost(postID int32, userID int32) error {
	if m.repostPostFunc != nil {
		return m.repostPostFunc(postID, userID)
	}
	return nil
}

func (m *mockDbClient) removeRepost(postID int32, userID int32) error {
	if m.removeRepostFunc != nil {
		return m.removeRepostFunc(postID, userID)
	}
	return nil
}

// Helper function to create authenticated context
func withUserID(ctx context.Context, userID int32) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD{"user-id": []string{intToString(userID)}})
}

func intToString(n int32) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.CreatePostRequest
		mockFunc  func(content string, userID int32) (int32, error)
		wantID    int32
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful creation",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.CreatePostRequest{Content: "Hello World"},
			mockFunc: func(content string, userID int32) (int32, error) {
				if content != "Hello World" || userID != 1 {
					t.Errorf("unexpected arguments: content=%s, userID=%d", content, userID)
				}
				return 123, nil
			},
			wantID:    123,
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.CreatePostRequest{Content: "Hello"},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.CreatePostRequest{Content: "Hello"},
			mockFunc: func(content string, userID int32) (int32, error) {
				return 0, errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{createPostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.CreatePost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp.Id != tt.wantID {
				t.Errorf("expected ID %d, got %d", tt.wantID, resp.Id)
			}
		})
	}
}

func TestGetFeed(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.GetFeedRequest
		mockFunc  func(page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
		wantPosts int
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful feed retrieval",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.GetFeedRequest{Page: 0, Limit: 10},
			mockFunc: func(page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
				return []*pb.Post{
					{Id: 1, Content: "Post 1"},
					{Id: 2, Content: "Post 2"},
				}, nil
			},
			wantPosts: 2,
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.GetFeedRequest{},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.GetFeedRequest{},
			mockFunc: func(page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
				return nil, errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{getFeedFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.GetFeed(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if len(resp.Posts) != tt.wantPosts {
				t.Errorf("expected %d posts, got %d", tt.wantPosts, len(resp.Posts))
			}
		})
	}
}

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.GetPostsRequest
		mockFunc  func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
		wantPosts int
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful posts retrieval",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.GetPostsRequest{UserId: 2, Page: 0, Limit: 10},
			mockFunc: func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
				if userID != 2 || currentUserID != 1 {
					t.Errorf("unexpected arguments: userID=%d, currentUserID=%d", userID, currentUserID)
				}
				return []*pb.Post{{Id: 1, Content: "Post 1"}}, nil
			},
			wantPosts: 1,
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.GetPostsRequest{UserId: 1},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.GetPostsRequest{UserId: 2},
			mockFunc: func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
				return nil, errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{getPostsFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.GetPosts(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if len(resp.Posts) != tt.wantPosts {
				t.Errorf("expected %d posts, got %d", tt.wantPosts, len(resp.Posts))
			}
		})
	}
}

func TestGetLikedPosts(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.GetPostsRequest
		mockFunc  func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error)
		wantPosts int
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful liked posts retrieval",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.GetPostsRequest{UserId: 2, Page: 0, Limit: 10},
			mockFunc: func(userID int32, page int32, limit int32, currentUserID int32) ([]*pb.Post, error) {
				return []*pb.Post{{Id: 1, Content: "Liked Post"}}, nil
			},
			wantPosts: 1,
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.GetPostsRequest{UserId: 1},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{getLikedPostsFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.GetLikedPosts(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if len(resp.Posts) != tt.wantPosts {
				t.Errorf("expected %d posts, got %d", tt.wantPosts, len(resp.Posts))
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(id int32, currentUserID int32) (*pb.Post, error)
		wantPost  *pb.Post
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful post retrieval",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(id int32, currentUserID int32) (*pb.Post, error) {
				if id != 123 || currentUserID != 1 {
					t.Errorf("unexpected arguments: id=%d, currentUserID=%d", id, currentUserID)
				}
				return &pb.Post{Id: 123, Content: "Test Post"}, nil
			},
			wantPost:  &pb.Post{Id: 123, Content: "Test Post"},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "post not found",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 999},
			mockFunc: func(id int32, currentUserID int32) (*pb.Post, error) {
				return nil, errors.New("not found")
			},
			wantError: true,
			wantCode:  codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{getPostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.GetPost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp.Id != tt.wantPost.Id || resp.Content != tt.wantPost.Content {
				t.Errorf("expected post %+v, got %+v", tt.wantPost, resp)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(postID int32, userID int32) error
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful deletion",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				if postID != 123 || userID != 1 {
					t.Errorf("unexpected arguments: postID=%d, userID=%d", postID, userID)
				}
				return nil
			},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{deletePostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.DeletePost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp == nil {
				t.Error("expected non-nil response")
			}
		})
	}
}

func TestLikePost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(postID int32, userID int32) error
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful like",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				if postID != 123 || userID != 1 {
					t.Errorf("unexpected arguments: postID=%d, userID=%d", postID, userID)
				}
				return nil
			},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{likePostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.LikePost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp == nil {
				t.Error("expected non-nil response")
			}
		})
	}
}

func TestUnlikePost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(postID int32, userID int32) error
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful unlike",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return nil
			},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{unlikePostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.UnlikePost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp == nil {
				t.Error("expected non-nil response")
			}
		})
	}
}

func TestRepostPost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(postID int32, userID int32) error
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful repost",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return nil
			},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{repostPostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.RepostPost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp == nil {
				t.Error("expected non-nil response")
			}
		})
	}
}

func TestRemoveRepost(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		req       *pb.PostRequest
		mockFunc  func(postID int32, userID int32) error
		wantError bool
		wantCode  codes.Code
	}{
		{
			name: "successful remove repost",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return nil
			},
			wantError: false,
		},
		{
			name:      "unauthenticated",
			ctx:       context.Background(),
			req:       &pb.PostRequest{PostId: 123},
			wantError: true,
			wantCode:  codes.Unauthenticated,
		},
		{
			name: "database error",
			ctx:  withUserID(context.Background(), 1),
			req:  &pb.PostRequest{PostId: 123},
			mockFunc: func(postID int32, userID int32) error {
				return errors.New("db error")
			},
			wantError: true,
			wantCode:  codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDbClient{removeRepostFunc: tt.mockFunc}
			ctrl := newController(mock)

			resp, err := ctrl.RemoveRepost(tt.ctx, tt.req)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
					return
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

			if resp == nil {
				t.Error("expected non-nil response")
			}
		})
	}
}
