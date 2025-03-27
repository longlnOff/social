CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
