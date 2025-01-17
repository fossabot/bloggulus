package core

import (
	"context"
)

type Blog struct {
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
	Title   string `json:"title"`

	// readonly (from database, after creation)
	ID int `json:"id"`
}

func NewBlog(feedURL, siteURL, title string) Blog {
	blog := Blog{
		FeedURL: feedURL,
		SiteURL: siteURL,
		Title:   title,
	}
	return blog
}

type BlogStorage interface {
	CreateBlog(ctx context.Context, blog *Blog) error
	ReadBlog(ctx context.Context, id int) (Blog, error)
	ReadBlogs(ctx context.Context, limit, offset int) ([]Blog, error)
}
