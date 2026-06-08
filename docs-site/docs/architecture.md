---
sidebar_position: 5
title: Architecture
---

# Architecture

A **Go** orchestration engine runs the pipeline against per-country provider
adapters, persists everything in **Postgres**, and a **Next.js** wizard drives it.

```
        ┌─────────────────────┐         ┌──────────────────────────────┐
        │   Next.js wizard     │  HTTP   │      Go orchestrator API      │
        │  country → details   │ ──────► │   chi router · /api/v1/...    │
        │  → review → run      │ ◄────── │                               │
        └─────────────────────┘  JSON   │  ┌─────────────────────────┐  │
                                         │  │  orchestrator (engine)  │  │
                                         │  └───────────┬─────────────┘  │
                                         │   ┌──────────┴───────────┐    │
                                         │   ▼          ▼           ▼    │
                                         │  IN         PH          US    │  adapters
                                         │   │ shared clients: Claude,   │
                                         │   │ RDAP, trade.gov, Razorpay │
                                         │   │ Stripe, PayMongo          │
                                         │   └──────────────────────┘    │
                                         └───────────────┬───────────────┘
                                                         │  pgx
                                                         ▼
                                                ┌─────────────────┐
                                                │    Postgres     │  businesses · launch_steps
                                                └─────────────────┘
```

## Components

| Component | Role |
|-----------|------|
| **Next.js wizard** | Country → details → review → run. Renders the live plan and per-step status from the API. Plus an in-app *How it works* and *Deploy* page. |
| **Go orchestrator** | A `chi` HTTP API + a pipeline engine that runs each step against the right country adapter and records the result. |
| **Provider adapters** | One per country (`india.go`, `ph.go`, `us.go`), each calling KYC, registry, tax, banking, payment and compliance providers. |
| **Shared clients** | Country-agnostic clients: Claude (strategy), RDAP + trademark (IP), trade.gov (liabilities), Razorpay/Stripe/PayMongo (payments). |
| **Postgres** | Stores every business and step — request, response, external reference and status — so launches resume. |

## Data model

- **`businesses`** — one launch application (country, entity type, founder, address, status).
- **`launch_steps`** — the ordered steps for a business, each with `status`,
  `provider`, `mode`, `request`, `response`, `external_ref`, `error`.

Migrations are embedded in the binary and applied automatically on startup.

## Orchestration semantics

- `POST /businesses` creates the business **and** seeds its step plan.
- `POST /businesses/{id}/advance` runs the next pending (or previously failed) step.
- `POST /businesses/{id}/run` runs all remaining steps, stopping at the first failure.
- A launch is **resumable**: re-run `advance` after fixing the cause of a failure.

See the **[API reference](./api-reference.md)** for the full surface.
