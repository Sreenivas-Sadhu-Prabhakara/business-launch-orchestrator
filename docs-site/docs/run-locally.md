---
sidebar_position: 7
title: Run it locally
---

# Run it locally

## Quick start (Docker)

You need **Docker** + **Docker Compose**. Nothing else.

```bash
git clone https://github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator.git
cd business-launch-orchestrator

# optional: add sandbox keys (skip to run fully in mock mode)
cp .env.example .env

docker compose up --build
```

Then open:

- **UI** → http://localhost:3000
- **API** → http://localhost:8080
- **Postgres** → localhost:5432 (`postgres` / `postgres`, db `biz_launch`)

:::tip Port already in use?
The compose file exposes overridable host ports. If `5432` / `3000` / `8080` are
taken, start with overrides:

```bash
DB_PORT=5434 WEB_PORT=3001 API_PORT=8080 docker compose up --build
```
:::

To stop and wipe the database volume:

```bash
docker compose down -v
```

## Without Docker

### 1. Postgres

```bash
docker run --name biz-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=biz_launch \
  -p 5432:5432 -d postgres:16-alpine
```

### 2. Go backend (migrations auto-apply on boot)

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/biz_launch?sslmode=disable"
# optional: export STRIPE_SECRET_KEY=sk_test_xxx  / ANTHROPIC_API_KEY=...
go run ./cmd/server
```

### 3. Next.js frontend

```bash
cd frontend
npm install
npm run dev
```

Open **http://localhost:3000**.

## Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `DATABASE_URL` | local Postgres | Connection string |
| `DB_MAX_CONNS` | `10` | Pool size (set `2` on serverless) |
| `PORT` | `8080` | API port |
| `CORS_ORIGIN` | `*` | Allowed origin |
| `ANTHROPIC_API_KEY` / `ANTHROPIC_MODEL` | — / `claude-sonnet-4-6` | Live AI strategy |
| `CSL_API_KEY` | — | Live sanctions screening |
| `RAZORPAY_KEY_ID` / `RAZORPAY_KEY_SECRET` | — | Live payments (IN) |
| `STRIPE_SECRET_KEY` | — | Live payments (US) |
| `PAYMONGO_SECRET_KEY` | — | Live payments (PH) |
| `FORCE_MOCK` | `false` | Force every step into mock mode |
