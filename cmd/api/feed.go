package main

import (
	"net/http"

	"github.com/longlnOff/social/internal/store"
)



func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pagination := store.PaginatedFeed{
		Limit: 20,
		Offset: 0,
		Sort: "desc",
	}

	pagination, err := pagination.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(pagination); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	feeds, err := app.store.Post.GetUserFeed(ctx, int64(20), pagination)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feeds); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
