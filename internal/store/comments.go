package store

import (
	"context"
	"database/sql"
)


type Comment struct {
	ID 			int64 	`json:"id"`
	PostID		int64 	`json:"post_id"`
	UserID 		int64 	`json:"user_id"`
	Content		string 	`json:"content"`
	CreatedAt	string 	`json:"created_at"`
	User 		User 	`json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func NewComment(db *sql.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (c *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return []Comment{}, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.User = User{}
		err := rows.Scan(&comment.ID, 
						 &comment.PostID,
						 &comment.UserID,
						 &comment.Content,
						 &comment.CreatedAt,
						 &comment.User.Username,	// Note that this auto unmarshals into the User struct
						 &comment.User.ID)	      	// Note that this auto unmarshals into the User struct
		if err != nil {
			return []Comment{}, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return []Comment{}, err
	}
	return comments, nil
}

func (c *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return c.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
}
