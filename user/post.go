package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Post struct {
	ID          string
	Username    string
	Image       []byte
	Description string
	CreatedAt   time.Time
}

type SortOrder string

var Ascending SortOrder = "ASC"
var Descending SortOrder = "DESC"

type PostDBPGX struct {
	URL string
}

func (db *PostDBPGX) CreatePost(post *Post) (string, error) {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())

	var id string
	err = conn.QueryRow(context.Background(), "INSERT INTO post (username, image, description, created_at) VALUES ($1, $2, $3, $4) RETURNING (id)", post.Username, post.Image, post.Description, post.CreatedAt).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

/// List post with username order by created at
/// Returned post does NOT contain image
func (db *PostDBPGX) ListPost(username string, createdAtSortOrder SortOrder) ([]*Post, error) {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	sortOrder := createdAtSortOrder

	rows, err := conn.Query(context.Background(),
		fmt.Sprintf("SELECT id, username, description, created_at FROM post WHERE username = COALESCE(NULLIF($1, ''), username) ORDER BY created_at %s", sortOrder),
		username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postList := make([]*Post, 0)
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Description, &post.CreatedAt); err != nil {
			return nil, err
		}
		postList = append(postList, &post)
	}

	return postList, nil
}

/// Get post image
/// Returned image
func (db *PostDBPGX) GetPostImage(id string) ([]byte, error) {
	conn, err := pgx.Connect(context.Background(), db.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	b := make([]byte, 0)
	err = conn.QueryRow(context.Background(), "SELECT image FROM post WHERE id = $1", id).Scan(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
