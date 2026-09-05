-- +goose Up

CREATE TABLE shortened_links (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    original_url TEXT NOT NULL,
    short_name TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_shortened_links_short_name
    ON shortened_links (short_name);

CREATE INDEX idx_shortened_links_deleted_at
    ON shortened_links (deleted_at);


-- +goose Down

DROP TABLE IF EXISTS shortened_links;