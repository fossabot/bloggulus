package test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
)

func CreatePost(storage core.Storage, t *testing.T) {
	post := CreateMockPost(storage, t)

	if post.ID == 0 {
		t.Fatal("post id after creation should be nonzero")
	}
}

func CreatePostAlreadyExists(storage core.Storage, t *testing.T) {
	post := CreateMockPost(storage, t)

	// attempt to create the same post again
	err := storage.CreatePost(context.Background(), &post)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate post should return an error")
	}
}

func ReadPost(storage core.Storage, t *testing.T) {
	post := CreateMockPost(storage, t)

	got, err := storage.ReadPost(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != post.ID {
		t.Fatalf("want %v, got %v", post.ID, got.ID)
	}
}

func ReadPosts(storage core.Storage, t *testing.T) {
	CreateMockPost(storage, t)
	CreateMockPost(storage, t)
	CreateMockPost(storage, t)
	CreateMockPost(storage, t)
	post := CreateMockPost(storage, t)

	limit := 3
	offset := 0
	posts, err := storage.ReadPosts(context.Background(), limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != limit {
		t.Fatalf("want %v, got %v", limit, len(posts))
	}

	// most recent post should be the one just added
	if posts[0].ID != post.ID {
		t.Fatalf("want %v, got %v", post.ID, posts[0].ID)
	}
}

func ReadPostsByBlog(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	// create 5 posts leaving the most recent one in "post"
	var post core.Post
	for i := 0; i < 5; i++ {
		post = core.NewPost(RandomURL(32), RandomString(32), RandomTime(), blog)
		err := storage.CreatePost(context.Background(), &post)
		if err != nil {
			t.Fatal(err)
		}
	}

	limit := 3
	offset := 0
	posts, err := storage.ReadPostsByBlog(context.Background(), blog.ID, limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != limit {
		t.Fatalf("want %v, got %v", limit, len(posts))
	}

	// most recent post should be the one just added
	if posts[0].ID != post.ID {
		t.Fatalf("want %v, got %v", post.ID, posts[0].ID)
	}
}

func SearchPosts(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)
	q := "python rust"

	// create 5 searchable posts leaving the most recent one in "post"
	var post core.Post
	for i := 0; i < 5; i++ {
		post = core.NewPost(RandomURL(32), q, RandomTime(), blog)
		err := storage.CreatePost(context.Background(), &post)
		if err != nil {
			t.Fatal(err)
		}
	}

	limit := 3
	offset := 0
	posts, err := storage.SearchPosts(context.Background(), q, limit, offset)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != limit {
		t.Fatalf("want %v, got %v", limit, len(posts))
	}

	// tags will always come back sorted desc
	tags := []string{"Python", "Rust"}
	if !subset(tags, posts[0].Tags) {
		t.Fatalf("want superset of %v, got %v", tags, posts[0].Tags)
	}
}

func CountPosts(storage core.Storage, t *testing.T) {
	CreateMockPost(storage, t)

	count, err := storage.CountPosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v", count)
	}
}

func CountSearchPosts(storage core.Storage, t *testing.T) {
	blog := CreateMockBlog(storage, t)

	// generate some searchable post data
	post := core.NewPost(
		RandomURL(32),
		"python rust",
		RandomTime(),
		blog,
	)

	// create a searchable post
	err := storage.CreatePost(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	count, err := storage.CountSearchPosts(context.Background(), "python rust")
	if err != nil {
		t.Fatal(err)
	}

	// ensure count is at least one
	if count < 1 {
		t.Fatalf("want >= 1, got %v", count)
	}
}

func CreateMockPost(storage core.Storage, t *testing.T) core.Post {
	t.Helper()

	// generate some random blog data
	blog := NewMockBlog()

	// create an example blog
	err := storage.CreateBlog(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// generate some random post data
	post := NewMockPost(blog)

	// create an example post
	err = storage.CreatePost(context.Background(), &post)
	if err != nil {
		t.Fatal(err)
	}

	return post
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
