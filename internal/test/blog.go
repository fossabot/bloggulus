package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func BlogCreate(storage core.Storage, t *testing.T) {
	blog := createMockBlog(storage, t)

	// blog should have an ID after creation
	if blog.ID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func BlogCreateExists(storage core.Storage, t *testing.T) {
	blog := createMockBlog(storage, t)

	// attempt to create the same blog again
	err := storage.BlogCreate(context.Background(), &blog)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func BlogReadAll(storage core.Storage, t *testing.T) {
	createMockBlog(storage, t)

	blogs, err := storage.BlogReadAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(blogs) < 1 {
		t.Fatalf("want >= 1, got %v\n", len(blogs))
	}
}

func createMockBlog(storage core.Storage, t *testing.T) core.Blog {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	return blog
}