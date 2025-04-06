package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/longlnOff/social/internal/store"
)

type UserCTX string

var Userctx UserCTX = "user"

type FollowedPayload struct {
	UserID int64 `json:"user_id" validate:"required"`
}

// ActivateUser godoc
//
//	@Summary		Activate a user
//	@Description	Activates a user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if err := app.store.User.ActivateUser(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, "User activated"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getUserHandler godoc
//
//	@Summary		Get user details
//	@Description	Retrieves a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int			true	"User ID"
//	@Success		200		{object}	store.User	"User details"
//	@Failure		404		{object}	string		"User not found"
//	@Failure		500		{object}	string		"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.getUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// followUserHandler godoc
//
//	@Summary		Follow a user
//	@Description	Creates a follow relationship between the authenticated user and target user
//	@Tags			users,follows
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int				true	"Follower User ID"
//	@Param			request	body		FollowedPayload	true	"User to follow"
//	@Success		204		{object}	nil				"No content"
//	@Failure		400		{object}	string			"Invalid request payload"
//	@Failure		409		{object}	string			"Already following this user"
//	@Failure		500		{object}	string			"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromCtx(r)
	followedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var followedPayload FollowedPayload
	if err := readJSON(w, r, &followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Follow(r.Context(), follower.ID, followedUserID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// unfollowUserHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	Removes a follow relationship between the authenticated user and target user
//	@Tags			users,follows
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int				true	"Follower User ID"
//	@Param			request	body		FollowedPayload	true	"User to unfollow"
//	@Success		204		{object}	nil				"No content"
//	@Failure		400		{object}	string			"Invalid request payload"
//	@Failure		404		{object}	string			"Not following this user"
//	@Failure		500		{object}	string			"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	follweruser := getUserFromCtx(r)
	followedUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var followedPayload FollowedPayload
	if err := readJSON(w, r, &followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Unfollow(r.Context(), follweruser.ID, followedUserID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.User.GetByUserID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, Userctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	return r.Context().Value(Userctx).(*store.User)
}
