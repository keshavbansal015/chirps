package api

import (
	"testing"

	pb "github.com/keshavbansal015/chirps/src/apigateway/genproto"
)

func TestMapUser(t *testing.T) {
	input := &pb.User{
		Id:        1,
		Name:      "John Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Bio:       "Software developer",
		Posts:     42,
		Likes:     100,
		Following: 50,
		Followers: 200,
		Followed:  true,
		Created:   "2024-01-01T00:00:00Z",
	}

	expected := user{
		ID:        1,
		Name:      "John Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Bio:       "Software developer",
		Posts:     42,
		Likes:     100,
		Following: 50,
		Followers: 200,
		Followed:  true,
		Created:   "2024-01-01T00:00:00Z",
	}

	result := mapUser(input)

	if result != expected {
		t.Errorf("mapUser() = %+v, want %+v", result, expected)
	}
}

func TestMapUser_ZeroValues(t *testing.T) {
	input := &pb.User{
		Id:       0,
		Name:     "",
		Username: "",
		Email:    "",
	}

	result := mapUser(input)

	if result.ID != 0 {
		t.Errorf("mapUser() ID = %v, want 0", result.ID)
	}
	if result.Name != "" {
		t.Errorf("mapUser() Name = %v, want empty string", result.Name)
	}
}

func TestMapPost(t *testing.T) {
	input := &pb.Post{
		Id:       1,
		UserId:   2,
		Content:  "Hello, world!",
		Likes:    10,
		Liked:    true,
		Reposts:  5,
		Reposted: false,
		Created:  "2024-01-01T00:00:00Z",
	}

	expected := post{
		ID:       1,
		UserID:   2,
		Content:  "Hello, world!",
		Likes:    10,
		Liked:    true,
		Reposts:  5,
		Reposted: false,
		Created:  "2024-01-01T00:00:00Z",
	}

	result := mapPost(input)

	if result != expected {
		t.Errorf("mapPost() = %+v, want %+v", result, expected)
	}
}

func TestMapPost_ZeroValues(t *testing.T) {
	input := &pb.Post{
		Id:      0,
		Content: "",
	}

	result := mapPost(input)

	if result.ID != 0 {
		t.Errorf("mapPost() ID = %v, want 0", result.ID)
	}
	if result.Content != "" {
		t.Errorf("mapPost() Content = %v, want empty string", result.Content)
	}
	if result.Liked != false {
		t.Errorf("mapPost() Liked = %v, want false", result.Liked)
	}
}
