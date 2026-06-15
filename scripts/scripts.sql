-- ==========================================
-- Users
-- ==========================================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ==========================================
-- Roles
-- ==========================================

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO roles(role)
VALUES
    ('admin'),
    ('manager'),
    ('user'),
    ('readonly');

-- ==========================================
-- User Roles (Many-to-Many)
-- ==========================================

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
);

-- ==========================================
-- Posts
-- ==========================================

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_posts_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ==========================================
-- Comments
-- ==========================================

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_comments_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_post
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE
);

-- ==========================================
-- Refresh Tokens
-- ==========================================

CREATE TABLE refresh_tokens (
    token_id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- ==========================================
-- Indexes
-- ==========================================

CREATE INDEX idx_posts_author_id
    ON posts(author_id);

CREATE INDEX idx_comments_author_id
    ON comments(author_id);

CREATE INDEX idx_comments_post_id
    ON comments(post_id);

CREATE INDEX idx_refresh_tokens_user_id
    ON refresh_tokens(user_id);

CREATE INDEX idx_roles_role
    ON roles(role);