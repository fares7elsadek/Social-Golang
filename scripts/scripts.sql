CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) ,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_roles(
    userId INT NOT NULL,
    roleId INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (userId,roleId)

    CONSTRAINT fk_users_userId
        FOREIGN KEY (userId)
        REFERENCES users(id)
        ON DELETE CASCADE
    
    CONSTRAINT fk_roles_roleId
        FOREIGN KEY (roleId)
        REFERENCES roles(id)
        ON DELETE CASCADE
)

CREATE TABLE roles(
    id SERIAL PRIMARY KEY,
    role TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
)


CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_posts_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    post_id INT NOT NULL,
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


CREATE TABLE refresh_token(
    id SERIAL PRIMARY KEY,
    tokenId TEXT NOT NULL,
    userId INT NOT NULL,
    ttl TIMESTAMP NOT NULL

    CONSTRAINT fk_refreshtoken_user
        FOREIGN KEY (userId)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_posts_author_id
    ON posts(author_id);

CREATE INDEX idx_comments_author_id
    ON comments(author_id);

CREATE INDEX idx_comments_post_id
    ON comments(post_id);

CREATE INDEX idx_refreshtoken_userid
    ON refresh_token(userId);

CREATE INDEX idx_refreshtoken_tokenid
    ON refresh_token(tokenId);


INSERT INTO roles(role)
Values("admin")
INSERT INTO roles(role)
Values("manager")
INSERT INTO roles(role)
Values("user")
INSERT INTO roles(role)
Values("readonly")

