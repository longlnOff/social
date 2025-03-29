CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
