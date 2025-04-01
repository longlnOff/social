package main

import (
	"net/http"

	"github.com/longlnOff/social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Retrieves posts for a user's feed with pagination
//	@Tags			feed,posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int			false	"Limit number of results"	default(20)
//	@Param			offset	query		int			false	"Offset for pagination"		default(0)
//	@Param			sort	query		string		false	"Sort order (asc or desc)"	default(desc)
//	@Success		200		{array}		store.Post	"Feed posts"
//	@Failure		400		{object}	string		"Invalid pagination parameters"
//	@Failure		500		{object}	string		"Internal Server Error"
//	@Router			/posts/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	pagination := store.PaginatedFeed{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
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
