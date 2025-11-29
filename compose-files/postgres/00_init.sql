CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts
(
    id         bigserial PRIMARY KEY,
    version    INT                                  DEFAULT 0,
    title      text                        NOT NULL,
    user_id    uuid                        NOT NULL,
    content    text                        NOT NULL,
    tags       VARCHAR(100)[],
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp(0) with time zone
);

CREATE TABLE IF NOT EXISTS comments
(
    id         bigserial PRIMARY KEY,
    post_id    bigserial                   NOT NULL,
    user_id    uuid                        NOT NULL,
    content    TEXT                        NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp(0) with time zone,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers
(
    user_id     uuid                        NOT NULL,
    follower_id bigint                      NOT NULL,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, follower_id)
);

-- Create the extension and indexes for full-text search
-- Check article: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE TABLE IF NOT EXISTS user_invitations
(
    token   bytea PRIMARY KEY,
    user_id uuid                        NOT NULL,
    expiry  TIMESTAMP(0) WITH TIME ZONE NOT NULL
);
