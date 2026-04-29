package post

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pb "github.com/keshavbansal015/chirps/src/postservice/genproto"
)

// mockRow implements the row interface for testing
type mockRow struct {
	id       int32
	userID   int32
	content  string
	likes    int32
	liked    bool
	reposts  int32
	reposted bool
	created  string
	err      error
}

func (m *mockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}

	if len(dest) != 8 {
		return errors.New("expected 8 destination arguments")
	}

	*(dest[0].(*int32)) = m.id
	*(dest[1].(*int32)) = m.userID
	*(dest[2].(*string)) = m.content
	*(dest[3].(*int32)) = m.likes
	*(dest[4].(*bool)) = m.liked
	*(dest[5].(*int32)) = m.reposts
	*(dest[6].(*bool)) = m.reposted
	*(dest[7].(*string)) = m.created

	return nil
}

// mockRows implements pgx.Rows interface for testing
type mockRows struct {
	posts   []*pb.Post
	current int
	closed  bool
	err     error
}

func (m *mockRows) Next() bool {
	if m.err != nil {
		return false
	}
	m.current++
	return m.current <= len(m.posts)
}

func (m *mockRows) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}

	if m.current == 0 || m.current > len(m.posts) {
		return errors.New("no current row")
	}

	post := m.posts[m.current-1]

	if len(dest) != 8 {
		return errors.New("expected 8 destination arguments")
	}

	*(dest[0].(*int32)) = post.Id
	*(dest[1].(*int32)) = post.UserId
	*(dest[2].(*string)) = post.Content
	*(dest[3].(*int32)) = post.Likes
	*(dest[4].(*bool)) = post.Liked
	*(dest[5].(*int32)) = post.Reposts
	*(dest[6].(*bool)) = post.Reposted
	*(dest[7].(*string)) = post.Created

	return nil
}

func (m *mockRows) Close() {
	m.closed = true
}

func (m *mockRows) Err() error {
	return m.err
}

func (m *mockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("")
}

func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m *mockRows) Values() ([]any, error) {
	return nil, nil
}

func (m *mockRows) RawValues() [][]byte {
	return nil
}

func (m *mockRows) Conn() *pgx.Conn {
	return nil
}

func TestMapPost(t *testing.T) {
	tests := []struct {
		name      string
		row       *mockRow
		wantPost  *pb.Post
		wantError bool
	}{
		{
			name: "valid post",
			row: &mockRow{
				id:       1,
				userID:   2,
				content:  "Hello World",
				likes:    10,
				liked:    true,
				reposts:  5,
				reposted: false,
				created:  "2024-01-01",
			},
			wantPost: &pb.Post{
				Id:       1,
				UserId:   2,
				Content:  "Hello World",
				Likes:    10,
				Liked:    true,
				Reposts:  5,
				Reposted: false,
				Created:  "2024-01-01",
			},
			wantError: false,
		},
		{
			name:      "scan error",
			row:       &mockRow{err: errors.New("scan failed")},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := mapPost(tt.row)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if post.Id != tt.wantPost.Id {
				t.Errorf("expected ID %d, got %d", tt.wantPost.Id, post.Id)
			}
			if post.UserId != tt.wantPost.UserId {
				t.Errorf("expected UserId %d, got %d", tt.wantPost.UserId, post.UserId)
			}
			if post.Content != tt.wantPost.Content {
				t.Errorf("expected Content %s, got %s", tt.wantPost.Content, post.Content)
			}
			if post.Likes != tt.wantPost.Likes {
				t.Errorf("expected Likes %d, got %d", tt.wantPost.Likes, post.Likes)
			}
			if post.Liked != tt.wantPost.Liked {
				t.Errorf("expected Liked %v, got %v", tt.wantPost.Liked, post.Liked)
			}
			if post.Reposts != tt.wantPost.Reposts {
				t.Errorf("expected Reposts %d, got %d", tt.wantPost.Reposts, post.Reposts)
			}
			if post.Reposted != tt.wantPost.Reposted {
				t.Errorf("expected Reposted %v, got %v", tt.wantPost.Reposted, post.Reposted)
			}
			if post.Created != tt.wantPost.Created {
				t.Errorf("expected Created %s, got %s", tt.wantPost.Created, post.Created)
			}
		})
	}
}

func TestMapPosts(t *testing.T) {
	tests := []struct {
		name       string
		rows       *mockRows
		wantLen    int
		wantError  bool
		wantClosed bool
	}{
		{
			name: "multiple posts",
			rows: &mockRows{
				posts: []*pb.Post{
					{Id: 1, UserId: 1, Content: "Post 1", Likes: 5, Liked: false, Reposts: 2, Reposted: true, Created: "2024-01-01"},
					{Id: 2, UserId: 2, Content: "Post 2", Likes: 10, Liked: true, Reposts: 3, Reposted: false, Created: "2024-01-02"},
					{Id: 3, UserId: 3, Content: "Post 3", Likes: 0, Liked: false, Reposts: 0, Reposted: false, Created: "2024-01-03"},
				},
			},
			wantLen:    3,
			wantError:  false,
			wantClosed: false,
		},
		{
			name:       "empty rows",
			rows:       &mockRows{posts: []*pb.Post{}},
			wantLen:    0,
			wantError:  false,
			wantClosed: false,
		},
		{
			name: "scan error on second row",
			rows: &mockRows{
				posts: []*pb.Post{
					{Id: 1, UserId: 1, Content: "Post 1", Likes: 5, Liked: false, Reposts: 2, Reposted: true, Created: "2024-01-01"},
				},
				err: errors.New("scan failed"),
			},
			wantLen:    0,
			wantError:  true,
			wantClosed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, err := mapPosts(tt.rows)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(posts) != tt.wantLen {
				t.Errorf("expected %d posts, got %d", tt.wantLen, len(posts))
			}

			if tt.rows.closed != tt.wantClosed {
				t.Errorf("expected rows closed=%v, got %v", tt.wantClosed, tt.rows.closed)
			}
		})
	}
}
