package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/longlnOff/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=3,max=100"`
	Content string   `json:"content" validate:"required,min=3,max=1000"`
	Tags    []string `json:"tags" validate:"required"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=100"`
	Content *string `json:"content" validate:"omitempty,min=3,max=1000"`
}

type CreateCommentForPostPayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,min=3,max=1000"`
}

type CreateCommentForPostResponse struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type PostCTX string

var Postctx PostCTX = "post"

// createPostHandler godoc
//
//	@Summary		Create a new post
//	@Description	Creates a new post with the provided title, content, and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreatePostPayload	true	"Post creation data"
//	@Success		201		{object}	store.Post			"Created post"
//	@Failure		400		{object}	string				"Invalid request payload"
//	@Failure		500		{object}	string				"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID: user.ID,
		Tags:   payload.Tags,
	}
	if err := app.store.Post.Create(r.Context(), &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPostsHandler godoc
//
//	@Summary		Get post details
//	@Description	Retrieves a post by ID including its comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int			true	"Post ID"
//	@Success		200		{object}	store.Post	"Post with comments"
//	@Failure		404		{object}	string		"Post not found"
//	@Failure		500		{object}	string		"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	// Get comments
	comments, err := app.store.Comment.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// deletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	Deletes a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Success		204		{object}	nil		"No content"
//	@Failure		404		{object}	string	"Post not found"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	err = app.store.Post.Delete(r.Context(), postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// updatePostHandler godoc
//
//	@Summary		Update a post
//	@Description	Updates a post's title and/or content
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			request	body		UpdatePostPayload	true	"Post update data"
//	@Success		200		{object}	store.Post			"Updated post"
//	@Failure		400		{object}	string				"Invalid request payload"
//	@Failure		404		{object}	string				"Post not found"
//	@Failure		500		{object}	string				"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Post.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// createCommentHandler godoc
//
//	@Summary		Create a comment
//	@Description	Creates a new comment for a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int								true	"Post ID"
//	@Param			request	body		CreateCommentForPostPayload		true	"Comment creation data"
//	@Success		201		{object}	CreateCommentForPostResponse	"Created comment"
//	@Failure		400		{object}	string							"Invalid request payload"
//	@Failure		404		{object}	string							"Post not found"
//	@Failure		500		{object}	string							"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentForPostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	comment := store.Comment{PostID: post.ID, UserID: payload.UserID, Content: payload.Content}
	if err := app.store.Comment.Create(r.Context(), &comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	res := CreateCommentForPostResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	if err := app.jsonResponse(w, http.StatusCreated, res); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		post, err := app.store.Post.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, Postctx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	return r.Context().Value(Postctx).(*store.Post)
}
