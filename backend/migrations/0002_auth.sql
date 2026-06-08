-- Authentication & authorization: user accounts (with roles) and per-launch
-- ownership. Applied automatically on startup.

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    username      TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL DEFAULT 'user',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Each launch belongs to the user who created it (NULL for legacy rows).
ALTER TABLE businesses
    ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_businesses_owner ON businesses(owner_id);
