package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func PostCreate(storage core.Storage, t *testing.T) {
	_, post := createMockBlogAndPost(storage, t)

	if post.ID == 0 {
		t.Fatal("post id after creation should be nonzero")
	}
}

func PostCreateExists(storage core.Storage, t *testing.T) {
	_, post := createMockBlogAndPost(storage, t)

	// attempt to create the same post again
	err := storage.PostCreate(context.Background(), &post)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate post should return an error")
	}
}

func PostReadAllByBlog(storage core.Storage, t *testing.T) {
	blog, _ := createMockBlogAndPost(storage, t)

	posts, err := storage.PostReadAllByBlog(context.Background(), blog.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 1 {
		t.Fatal("expected one post linked to blog")
	}
}

func PostReadRecent(storage core.Storage, t *testing.T) {
	_, post := createMockBlogAndPost(storage, t)

	posts, err := storage.PostReadRecent(context.Background(), 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	// most recent post should be the one just added
	if posts[0].ID != post.ID {
		t.Fatalf("want %v, got %v\n", post.ID, posts[0].ID)
	}
}

func PostReadSearch(storage core.Storage, t *testing.T) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some searchable post data
	post := core.NewPost(
		RandomURL(32),
		"python rust",
		RandomTime(),
		blog,
	)

	// create a searchable post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := storage.PostReadSearch(context.Background(), "python rust", 20, 0)
	if err != nil {
		t.Fatal(err)
	}

	// tags will always come back sorted desc
	tags := []string{"Python", "Rust"}
	if !subset(tags, posts[0].Tags) {
		t.Fatalf("want superset of %v, got %v\n", tags, posts[0].Tags)
	}
}

func PostCountRecent(storage core.Storage, t *testing.T) {
	createMockBlogAndPost(storage, t)

	count, err := storage.PostCountRecent(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v\n", count)
	}
}

func PostCountSearch(storage core.Storage, t *testing.T) {
	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some searchable post data
	post := core.NewPost(
		RandomURL(32),
		"python rust",
		RandomTime(),
		blog,
	)

	// create a searchable post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	count, err := storage.PostCountSearch(context.Background(), "python rust")
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v\n", count)
	}
}

func createMockBlogAndPost(storage core.Storage, t *testing.T) (core.Blog, core.Post) {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.BlogCreate(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = storage.PostCreate(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	return blog, post
}

func subset(a, b []string) bool {
	bset := make(map[string]bool)
	for _, s := range b {
		bset[s] = true
	}

	for _, s := range a {
		if _, ok := bset[s]; !ok {
			return false
		}
	}

	return true
}