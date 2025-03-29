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


func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
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
