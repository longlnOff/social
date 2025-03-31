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

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromCtx(r)

	var followedPayload FollowedPayload
	if err := readJSON(w, r, &followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Follow(r.Context(), follower.ID, followedPayload.UserID); err != nil {
		switch {
			case errors.Is(err, store.ErrConflict):
				app.conflictResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromCtx(r)

	var followedPayload FollowedPayload
	if err := readJSON(w, r, &followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(followedPayload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Unfollow(r.Context(), follower.ID, followedPayload.UserID); err != nil {
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
