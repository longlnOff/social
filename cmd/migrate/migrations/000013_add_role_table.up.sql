CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description) VALUES
    ('user', 0, 'A user can create posts and comments'),
    ('moderator', 1, 'A moderator can update other posts and comments'),
    ('admin', 2, 'An admin can update and delete any post or comment');
