package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/theandrew168/bloggulus/internal/core"
)

func (app *Application) HandleReadPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	post, err := app.storage.ReadPost(ctx, id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, 200, envelope{"post": post})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleReadPosts(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	q := r.URL.Query().Get("q")

	var posts []core.Post
	if q != "" {
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// search if requested
		var err error
		posts, err = app.storage.SearchPosts(ctx, q, limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
		defer cancel()

		// else just read recent
		var err error
		posts, err = app.storage.ReadPosts(ctx, limit, offset)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = writeJSON(w, 200, envelope{"posts": posts})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
