CREATE TABLE post (
    post_id SERIAL PRIMARY KEY,
    blog_id INTEGER NOT NULL REFERENCES blog,
    url TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    updated TIMESTAMPTZ NOT NULL
)
