# 🚀 Business Launch Orchestrator

**One end-to-end flow that integrates every API call needed to launch a business in 🇮🇳 India, 🇵🇭 the Philippines, or 🇺🇸 the United States** — founder KYC → name reservation → incorporation → tax registration → business banking → payment gateway → statutory compliance.

A **Go** orchestration engine runs a country-specific pipeline of provider integrations, persists every step + response in **Postgres**, and a **Next.js** wizard walks you through it visually.

> Integration mode is **hybrid**: payment steps hit **real provider sandboxes** (Razorpay / Stripe / PayMongo) when you supply test keys, and every other step (government registries, KYC, banking) is a **deterministic mock** whose request/response shape mirrors the real upstream API — so the whole thing runs end-to-end with **zero credentials**, and you wire real APIs in incrementally.

---

## 📐 Architecture

```
        ┌─────────────────────┐         ┌──────────────────────────────┐
        │   Next.js wizard     │  HTTP   │      Go orchestrator API      │
        │  (country → details  │ ──────► │   chi router · /api/v1/...    │
        │   → review → run)    │ ◄────── │                               │
        └─────────────────────┘  JSON   │  ┌─────────────────────────┐  │
                                         │  │  orchestrator (engine)  │  │
                                         │  └───────────┬─────────────┘  │
                                         │   ┌──────────┴───────────┐    │
                                         │   ▼          ▼           ▼    │
                                         │  IN         PH          US    │  provider adapters
                                         │  adapter    adapter   adapter │
                                         │   │ Razorpay  │ PayMongo │ Stripe  ← LIVE sandbox
                                         │   └──────────────────────┘    │
                                         └───────────────┬───────────────┘
                                                         │  pgx
                                                         ▼
                                                ┌─────────────────┐
                                                │    Postgres     │  businesses · launch_steps
                                                └─────────────────┘
```

## 🧭 The launch pipeline

Every country runs the same 7 logical steps; each maps to different real upstreams:

| # | Step | 🇮🇳 India | 🇵🇭 Philippines | 🇺🇸 United States |
|---|------|----------|----------------|------------------|
| 1 | Founder KYC | SurePass (PAN + Aadhaar) | HyperVerge (PhilID) | Persona / Middesk |
| 2 | Name check | MCA RUN | SEC name verification | Secretary of State |
| 3 | Entity registration | MCA SPICe+ → **CIN** | SEC eSPARC → **SEC reg no.** | State filing → **filing no.** |
| 4 | Tax registration | Income Tax + GSTN → **PAN/TAN/GSTIN** | BIR → **TIN** | IRS SS-4 → **EIN** |
| 5 | Business banking | RazorpayX | UnionBank | Mercury |
| 6 | **Payment gateway** | **Razorpay** 🟢 | **PayMongo** 🟢 | **Stripe** 🟢 |
| 7 | Compliance | EPFO + ESIC | SSS + PhilHealth + Pag-IBIG | Registered agent + state tax |

🟢 = real sandbox call when a test key is configured, otherwise deterministic mock.

---

## ⚡ Quick start (Docker — recommended)

You need **Docker** + **Docker Compose**. Nothing else.

**1. Clone**

```bash
git clone https://github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator.git
cd business-launch-orchestrator
```

**2. (Optional) add sandbox keys** — skip this and everything runs in mock mode.

```bash
cp .env.example .env
# then edit .env and paste any of: RAZORPAY_KEY_ID/SECRET, STRIPE_SECRET_KEY, PAYMONGO_SECRET_KEY
```

**3. Launch the whole stack**

```bash
docker compose up --build
```

**4. Open the app**

```bash
open http://localhost:3000      # UI wizard  (macOS; use xdg-open on Linux)
```

- UI → http://localhost:3000
- API → http://localhost:8080
- Postgres → localhost:5432 (`postgres` / `postgres`, db `biz_launch`)

To stop and wipe the database volume:

```bash
docker compose down -v
```

---

## 🛠️ Run it locally (without Docker)

### Step 1 — start Postgres

```bash
docker run --name biz-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=biz_launch -p 5432:5432 -d postgres:16-alpine
```

### Step 2 — run the Go backend (migrations apply automatically on boot)

```bash
cd backend
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/biz_launch?sslmode=disable"
# optional sandbox keys:
# export STRIPE_SECRET_KEY=sk_test_xxx
go run ./cmd/server
```

You should see `migrations applied` then `listening on :8080`.

### Step 3 — run the Next.js frontend (in a second terminal)

```bash
cd frontend
npm install
npm run dev
```

Open **http://localhost:3000**.

---

## 🔌 API reference (with copyable `curl`)

The UI is just a client of this API — you can drive the entire launch from the terminal.

### List supported countries and their plans

```bash
curl -s http://localhost:8080/api/v1/countries | jq
```

### Create a launch 🇮🇳

```bash
curl -s -X POST http://localhost:8080/api/v1/businesses \
  -H 'Content-Type: application/json' \
  -d '{
    "country": "IN",
    "entity_type": "Private Limited Company",
    "legal_name": "Acme Technologies",
    "founder_name": "Jane Doe",
    "founder_email": "jane@acme.com",
    "founder_phone": "+91 98000 12345",
    "founder_id_number": "ABCDE1234F",
    "address": { "line1": "12 MG Road", "city": "Bengaluru", "state": "Karnataka", "postal_code": "560001", "country": "IN" }
  }' | jq
```

Grab the returned `business.id`:

```bash
BIZ=$(curl -s -X POST http://localhost:8080/api/v1/businesses \
  -H 'Content-Type: application/json' \
  -d '{"country":"US","entity_type":"LLC","legal_name":"Globex LLC","founder_name":"John Roe","founder_email":"john@globex.com","address":{"state":"Delaware","country":"US"}}' \
  | jq -r '.business.id')
echo "business id: $BIZ"
```

### Run the next step (one at a time)

```bash
curl -s -X POST http://localhost:8080/api/v1/businesses/$BIZ/advance | jq
```

### Run the whole pipeline to completion

```bash
curl -s -X POST http://localhost:8080/api/v1/businesses/$BIZ/run | jq '.business.status, .steps[] | {seq, title, status, external_ref}'
```

### Fetch current state

```bash
curl -s http://localhost:8080/api/v1/businesses/$BIZ | jq
```

### Full endpoint list

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/healthz` | Liveness check |
| `GET`  | `/api/v1/countries` | All countries + their step plans |
| `GET`  | `/api/v1/countries/{code}/plan` | Plan for `IN` / `PH` / `US` |
| `POST` | `/api/v1/businesses` | Create a launch (returns `{business, steps}`) |
| `GET`  | `/api/v1/businesses` | List recent launches |
| `GET`  | `/api/v1/businesses/{id}` | Get one launch + its steps |
| `GET`  | `/api/v1/businesses/{id}/steps` | Steps only |
| `POST` | `/api/v1/businesses/{id}/advance` | Execute the next pending step |
| `POST` | `/api/v1/businesses/{id}/run` | Execute all remaining steps |

---

## 🟢 Going live (real sandbox calls)

The payment step is wired to real provider sandboxes. Add any of these to `.env`
(or export them before `go run`) and that step flips from `mock` to `live`:

| Country | Provider | Key(s) | Get them |
|---------|----------|--------|----------|
| 🇮🇳 IN | Razorpay | `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET` | Razorpay Dashboard → Settings → API Keys (`rzp_test_…`) |
| 🇺🇸 US | Stripe | `STRIPE_SECRET_KEY` | Stripe Dashboard → Developers → API keys (`sk_test_…`) |
| 🇵🇭 PH | PayMongo | `PAYMONGO_SECRET_KEY` | PayMongo Dashboard → Developers → API keys (`sk_test_…`) |

Each live call creates a real test-mode object (a Razorpay order / Stripe
customer / PayMongo payment link) and stores the returned id as the step's
`external_ref`. Set `FORCE_MOCK=true` to disable all live calls.

> **Wiring a government API for real:** every mock step in
> `backend/internal/providers/{india,ph,us}.go` documents its real upstream in a
> comment. Replace the mock body with an HTTP call (see `payments.go` for the
> pattern) and you're live — the orchestrator, persistence and UI don't change.

---

## 🗂️ Project structure

```
business-launch-orchestrator/
├── docker-compose.yml            # postgres + backend + frontend
├── .env.example                  # sandbox keys (all optional)
├── backend/                      # Go orchestration service
│   ├── cmd/server/main.go        # entrypoint, graceful shutdown, auto-migrate
│   ├── migrations/0001_init.sql  # embedded, applied on boot
│   └── internal/
│       ├── domain/               # dependency-free core types
│       ├── store/                # Postgres persistence (pgx)
│       ├── providers/            # per-country adapters + live payment clients
│       │   ├── india.go ph.go us.go
│       │   └── payments.go       # Razorpay / Stripe / PayMongo (real sandbox)
│       ├── orchestrator/         # the pipeline engine
│       └── api/                  # chi router + handlers
└── frontend/                     # Next.js App Router wizard
    ├── app/page.tsx              # country → details → review wizard
    ├── app/launch/[id]/page.tsx  # live step runner
    ├── components/               # Stepper, StepList
    └── lib/api.ts                # typed API client
```

## 🧩 Extending it

- **Add a country:** implement the `providers.Adapter` interface (`Country()`,
  `Plan()`, `Execute()`) in a new file, register it in `providers/providers.go`,
  add the code to `domain.Country.Valid()`. The UI picks it up automatically from
  `GET /api/v1/countries`.
- **Add a step type:** add a `StepType` constant in `domain`, include it in each
  country's `Plan()`, and handle it in each adapter's `Execute()` switch.

---

## ⚠️ Disclaimer

This is a reference implementation / scaffold. The government-registry, KYC and
banking steps are **simulations** — generated identifiers (CIN, EIN, GSTIN, TIN…)
are realistically formatted but **not real registrations**. Real incorporation
involves legal documents, signatures, fees and (often) licensed intermediaries.
Use this as the integration backbone; consult a professional before filing.

## 📄 License

MIT
