package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theandrew168/bloggulus/internal/core"
	"github.com/theandrew168/bloggulus/internal/postgresql"
	"github.com/theandrew168/bloggulus/internal/test"
)

func TestBlogCreate(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)

	// generate some random blog data
	feedURL := "https://" + test.RandomString(32)
	siteURL := "https://" + test.RandomString(32)
	title := test.RandomString(32)

	blog := core.NewBlog(feedURL, siteURL, title)
	if blog.BlogID != 0 {
		t.Fatal("blog id before creation should be zero")
	}

	// create an example blog
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	if blog.BlogID == 0 {
		t.Fatal("blog id after creation should be nonzero")
	}
}

func TestBlogCreateExists(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	// test connection to ensure all is well
	if err := conn.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)

	// generate some random blog data
	feedURL := "https://" + test.RandomString(32)
	siteURL := "https://" + test.RandomString(32)
	title := test.RandomString(32)

	blog := core.NewBlog(feedURL, siteURL, title)

	// create an example blog
	err := blogStorage.Create(context.Background(), &blog)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same blog again
	err = blogStorage.Create(context.Background(), &blog)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate blog should return an error")
	}
}

func TestBlogReadAll(t *testing.T) {
	conn := test.ConnectDB(t)
	defer conn.Close()

	// instantiate storage interfaces
	blogStorage := postgresql.NewBlogStorage(conn)

	_, err := blogStorage.ReadAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
