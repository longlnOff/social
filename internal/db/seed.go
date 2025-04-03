package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/longlnOff/social/internal/store"
)

func Seed(store *store.Storage, db *sql.DB) error {
	ctx := context.Background()

	// Create user
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for user := range users {
		if err := store.User.Create(ctx, tx, users[user]); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Create posts
	posts := generatePosts(100, users)
	for post := range posts {
		if err := store.Post.Create(ctx, posts[post]); err != nil {
			return err
		}
	}

	// Create comments
	comments := generateComments(100, posts, users)
	for comment := range comments {
		if err := store.Comment.Create(ctx, comments[comment]); err != nil {
			return err
		}
	}
	tx.Commit()
	
	return nil
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := range count {
		users[i] = &store.User{
			Username: "user" + strconv.Itoa(i),
			Email:    "user" + strconv.Itoa(i) + "@gmail.com",
		}
	}

	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := range count {
		posts[i] = &store.Post{
			Title:   "title" + strconv.Itoa(i),
			Content: "content" + strconv.Itoa(i),
			UserID:  users[i].ID,
			Tags:    []string{"tag1", "tag2", "tag3"},
		}
	}

	log.Println("Seeding completed")

	return posts
}

func generateComments(count int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := range count {
		comments[i] = &store.Comment{
			PostID:  posts[i].ID,
			UserID:  users[i].ID,
			Content: "comment" + strconv.Itoa(i),
		}
	}

	return comments
}
