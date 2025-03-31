package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID 			int64 `json:"id"`
	Content		string `json:"content"`
	Title		string `json:"title"`
	UserID 		int64 `json:"user_id"`
	Tags		[]string `json:"tags"`	
	Version		int64 `json:"version"`
	CreatedAt	string `json:"created_at"`
	UpdatedAt	string `json:"updated_at"`
	Comments 	[]Comment `json:"comments"`
	User 		User 	`json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `
		SELECT 
		p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username,
		COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $1 OR p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	feeds := []PostWithMetadata{}

	for rows.Next() {
		var feed PostWithMetadata
		err := rows.Scan(
			&feed.ID,
			&feed.UserID,
			&feed.Title,
			&feed.Content,
			&feed.CreatedAt,
			&feed.Version,
			pq.Array(&feed.Tags),
			&feed.User.Username,
			&feed.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, updated_at
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrNotFound
			default:
				return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, 
								query, 
								post.Title, 
								post.Content, 
								post.ID, 
								post.Version).Scan(&post.Version)
	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return ErrNotFound
			default:
				return err
		}
	}
	return err
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
