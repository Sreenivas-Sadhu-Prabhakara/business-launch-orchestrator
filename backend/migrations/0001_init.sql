-- Schema for the business launch orchestrator.
-- Applied automatically on backend startup (see internal/store/migrate.go).

CREATE TABLE IF NOT EXISTS businesses (
    id                UUID PRIMARY KEY,
    country           TEXT        NOT NULL,
    entity_type       TEXT        NOT NULL,
    legal_name        TEXT        NOT NULL,
    founder_name      TEXT        NOT NULL,
    founder_email     TEXT        NOT NULL DEFAULT '',
    founder_phone     TEXT        NOT NULL DEFAULT '',
    founder_id_number TEXT        NOT NULL DEFAULT '',
    address           JSONB       NOT NULL DEFAULT '{}'::jsonb,
    status            TEXT        NOT NULL DEFAULT 'draft',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS launch_steps (
    id           UUID PRIMARY KEY,
    business_id  UUID        NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    seq          INT         NOT NULL,
    step_type    TEXT        NOT NULL,
    provider     TEXT        NOT NULL,
    title        TEXT        NOT NULL DEFAULT '',
    mode         TEXT        NOT NULL DEFAULT 'mock',
    status       TEXT        NOT NULL DEFAULT 'pending',
    request      JSONB       NOT NULL DEFAULT '{}'::jsonb,
    response     JSONB       NOT NULL DEFAULT '{}'::jsonb,
    external_ref TEXT        NOT NULL DEFAULT '',
    error        TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    UNIQUE (business_id, seq)
);

CREATE INDEX IF NOT EXISTS idx_launch_steps_business ON launch_steps(business_id);
CREATE INDEX IF NOT EXISTS idx_businesses_status ON businesses(status);
