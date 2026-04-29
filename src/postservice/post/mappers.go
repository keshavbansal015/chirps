package post

import (
	"github.com/jackc/pgx/v5"

	pb "github.com/keshavbansal015/chirps/src/postservice/genproto"
)

type row interface {
	Scan(dest ...any) error
}

func mapPost(r row) (*pb.Post, error) {
	post := pb.Post{}

	err := r.Scan(&post.Id, &post.UserId, &post.Content, &post.Likes,
		&post.Liked, &post.Reposts, &post.Reposted, &post.Created)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func mapPosts(r pgx.Rows) ([]*pb.Post, error) {
	posts := []*pb.Post{}
	for r.Next() {
		post, err := mapPost(r)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
