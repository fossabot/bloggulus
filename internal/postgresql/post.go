package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/bloggulus/internal/core"
)

type postStorage struct {
	conn *pgxpool.Pool
}

func NewPostStorage(conn *pgxpool.Pool) core.PostStorage {
	s := postStorage{
		conn: conn,
	}
	return &s
}

func (s *postStorage) Create(ctx context.Context, post *core.Post) error {
	stmt := `
		INSERT INTO post
			(url, title, updated, body, blog_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING post_id`
	args := []interface{}{
		post.URL,
		post.Title,
		post.Updated,
		post.Body,
		post.Blog.BlogID,
	}
	row := s.conn.QueryRow(ctx, stmt, args...)

	err := row.Scan(&post.PostID)
	if err != nil {
		// https://github.com/jackc/pgx/wiki/Error-Handling
		// https://github.com/jackc/pgx/issues/474
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return core.ErrExist
			}
		}
		return err
	}

	return nil
}

func (s *postStorage) ReadAllByBlog(ctx context.Context, blogID int) ([]core.Post, error) {
	stmt := `
		SELECT
			post.post_id,
			post.url,
			post.title,
			post.updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(post.content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			blog.blog_id,
			blog.feed_url,
			blog.site_url,
			blog.title
		FROM post
		INNER JOIN blog
			ON blog.blog_id = post.blog_id
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ post.content_index
		WHERE blog.blog_id = $1
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY post.updated DESC`
	rows, err := s.conn.Query(ctx, stmt, blogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) ReadRecent(ctx context.Context, limit, offset int) ([]core.Post, error) {
	stmt := `
		WITH posts AS (
			SELECT
				post.post_id,
				post.url,
				post.title,
				post.updated,
				post.content_index,
				blog.blog_id AS blog_blog_id,
				blog.feed_url AS blog_feed_url,
				blog.site_url AS blog_site_url,
				blog.title AS blog_title
			FROM post
			INNER JOIN blog
				ON blog.blog_id = post.blog_id
			ORDER BY post.updated DESC
			LIMIT $1
			OFFSET $2
		)
		SELECT
			post_id,
			url,
			title,
			updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			blog_blog_id,
			blog_feed_url,
			blog_site_url,
			blog_title
		FROM posts
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ posts.content_index
		GROUP BY 1,2,3,4,6,7,8,9
		ORDER BY updated DESC`
	rows, err := s.conn.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) ReadSearch(ctx context.Context, query string, limit, offset int) ([]core.Post, error) {
	stmt := `
		WITH posts AS (
			SELECT
				post.post_id,
				post.url,
				post.title,
				post.updated,
				post.content_index,
				blog.blog_id AS blog_blog_id,
				blog.feed_url AS blog_feed_url,
				blog.site_url AS blog_site_url,
				blog.title AS blog_title
			FROM post
			INNER JOIN blog
				ON blog.blog_id = post.blog_id
			WHERE post.content_index @@ websearch_to_tsquery('english',  $1)
			ORDER BY ts_rank_cd(post.content_index, websearch_to_tsquery('english',  $1)) DESC
			LIMIT $2
			OFFSET $3
		)
		SELECT
			post_id,
			url,
			title,
			updated,
			array_remove(array_agg(tag.name ORDER BY ts_rank_cd(content_index, to_tsquery(tag.name)) DESC), NULL) as tags,
			blog_blog_id,
			blog_feed_url,
			blog_site_url,
			blog_title
		FROM posts
		LEFT JOIN tag
			ON to_tsquery(tag.name) @@ content_index
		GROUP BY 1,2,3,4,6,7,8,9,content_index
		ORDER BY ts_rank_cd(content_index, websearch_to_tsquery('english',  $1)) DESC`
	rows, err := s.conn.Query(ctx, stmt, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post
	for rows.Next() {
		var post core.Post
		err := rows.Scan(
			&post.PostID,
			&post.URL,
			&post.Title,
			&post.Updated,
			&post.Tags,
			&post.Blog.BlogID,
			&post.Blog.FeedURL,
			&post.Blog.SiteURL,
			&post.Blog.Title,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *postStorage) CountRecent(ctx context.Context) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post`
	row := s.conn.QueryRow(ctx, stmt)

	var count int
	err := row.Scan(&count)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// retry on stale connections
			if pgErr.Code == pgerrcode.AdminShutdown {
				return s.CountRecent(ctx)
			}
		}
		return 0, err
	}

	return count, nil
}

func (s *postStorage) CountSearch(ctx context.Context, query string) (int, error) {
	stmt := `
		SELECT count(*)
		FROM post
		WHERE content_index @@ websearch_to_tsquery('english',  $1)`
	row := s.conn.QueryRow(ctx, stmt, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// retry on stale connections
			if pgErr.Code == pgerrcode.AdminShutdown {
				return s.CountSearch(ctx, query)
			}
		}
		return 0, err
	}

	return count, nil
}
