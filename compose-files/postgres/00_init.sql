CREATE
EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id
    bigserial
    PRIMARY
    KEY,
    email
    citext
    UNIQUE
    NOT
    NULL,
    username
    varchar
(
    255
) UNIQUE NOT NULL,
    password bytea NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at timestamp
(
    0
) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp
(
    0
)
  with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp
(
    0
)
  with time zone
      );

CREATE TABLE IF NOT EXISTS posts
(
    id
    bigserial
    PRIMARY
    KEY,
    version
    INT
    DEFAULT
    0,
    title
    text
    NOT
    NULL,
    user_id
    bigint
    NOT
    NULL,
    content
    text
    NOT
    NULL,
    tags
    VARCHAR
(
    100
)[],
    created_at timestamp
(
    0
) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp
(
    0
)
  with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp
(
    0
)
  with time zone,
      FOREIGN KEY (user_id) REFERENCES users
(
    id
)
    );

ALTER TABLE posts
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);

CREATE TABLE IF NOT EXISTS comments
(
    id
    bigserial
    PRIMARY
    KEY,
    post_id
    bigserial
    NOT
    NULL,
    user_id
    bigserial
    NOT
    NULL,
    content
    TEXT
    NOT
    NULL,
    created_at
    timestamp
(
    0
) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp
(
    0
)
  with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp
(
    0
)
  with time zone,
      FOREIGN KEY (post_id) REFERENCES posts
(
    id
)
  ON DELETE CASCADE,
    FOREIGN KEY
(
    user_id
) REFERENCES users
(
    id
)
  ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS followers
(
    user_id
    bigint
    NOT
    NULL,
    follower_id
    bigint
    NOT
    NULL,
    created_at
    timestamp
(
    0
) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY
(
    user_id,
    follower_id
),
    FOREIGN KEY
(
    user_id
) REFERENCES users
(
    id
)
  ON DELETE CASCADE,
    FOREIGN KEY
(
    follower_id
) REFERENCES users
(
    id
)
  ON DELETE CASCADE
    );

-- Create the extension and indexes for full-text search
-- Check article: https://niallburkley.com/blog/index-columns-for-like-in-postgres/
CREATE
EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE TABLE IF NOT EXISTS user_invitations
(
    token
    bytea
    PRIMARY
    KEY,
    user_id
    bigint
    NOT
    NULL,
    expiry
    TIMESTAMP
(
    0
) WITH TIME ZONE NOT NULL
      );

CREATE TABLE IF NOT EXISTS roles
(
    id
    BIGSERIAL
    PRIMARY
    KEY,
    name
    VARCHAR
(
    255
) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
    );

ALTER TABLE
    IF EXISTS users
    ADD
    COLUMN role_id INT REFERENCES roles (id) DEFAULT 1;

INSERT INTO roles (name, description, level)
VALUES ('user',
        'A user can create posts and comments',
        1);

INSERT INTO roles (name, description, level)
VALUES ('moderator',
        'A moderator can update other users posts',
        2);

INSERT INTO roles (name, description, level)
VALUES ('admin',
        'An admin can update and delete other users posts',
        3);

ALTER TABLE
    users
    ALTER COLUMN
        role_id DROP DEFAULT;

ALTER TABLE
    users
    ALTER COLUMN
        role_id
        SET
        NOT NULL;

INSERT INTO users (username, email, password, is_active, role_id)
VALUES ('user1@test.com',
        'user1@test.com',
        '$2a$10$SCaWY2MSkIOfEY98tdYWReuN0aRFpzrZty2iuxEbma32byg7FfQYm',
        true,
        1),
       ('moderator1@test.com',
        'moderator1@test.com',
        '$2a$10$SCaWY2MSkIOfEY98tdYWReuN0aRFpzrZty2iuxEbma32byg7FfQYm',
        true,
        2),
       ('admin1@test.com',
        'admin1@test.com',
        '$2a$10$SCaWY2MSkIOfEY98tdYWReuN0aRFpzrZty2iuxEbma32byg7FfQYm',
        true,
        2)