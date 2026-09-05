-- +goose Up

CREATE TABLE IF NOT EXISTS shortened_link_visits (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    link_id BIGINT NOT NULL,
    ip TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    status BIGINT NOT NULL,
    referrer TEXT,

    CONSTRAINT fk_shortened_link_visits_link
        FOREIGN KEY (link_id)
        REFERENCES shortened_links (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_shortened_link_visits_deleted_at
    ON shortened_link_visits (deleted_at);

CREATE INDEX IF NOT EXISTS idx_shortened_link_visits_link_id
    ON shortened_link_visits (link_id);


-- +goose Down

DROP TABLE IF EXISTS shortened_link_visits;
