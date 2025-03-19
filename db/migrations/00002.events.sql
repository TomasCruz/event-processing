BEGIN;

CREATE TABLE IF NOT EXISTS events (
    id bigserial PRIMARY KEY,
    player_id bigserial NOT NULL,
    game_id bigserial NOT NULL,
    typ VARCHAR(12) NOT NULL,
    amount INTEGER,
    currency VARCHAR(3),
    has_won BOOLEAN,
    created_at timestamptz NOT NULL,
    amount_eur INTEGER,
    description VARCHAR(100)
);

CREATE INDEX idx_events_player_id ON events(player_id);

COMMIT;
