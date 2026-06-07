# 🚀 Business Launch Orchestrator

**One end-to-end flow that integrates every API call needed to launch a business in 🇮🇳 India, 🇵🇭 the Philippines, or 🇺🇸 the United States** — AI strategy → founder KYC → liabilities/sanctions screening → name + trademark/domain check → incorporation → tax registration → licenses & registrations → business banking → payment gateway → statutory compliance.

A **Go** orchestration engine runs a country-specific 11-step pipeline of provider integrations, persists every step + response in **Postgres**, and a **Next.js** wizard walks you through it visually. The app ships **serverless-first** (Lambda / Cloud Run / Container Apps) and includes an in-app **[How it works](#)** explainer and **[Deploy](#)** guide.

> Integration mode is **hybrid**: payment steps hit **real provider sandboxes** (Razorpay / Stripe / PayMongo), the AI strategy step calls **Claude**, the IP step does a **real RDAP domain lookup**, and sanctions screening uses the **trade.gov** list — each when you supply its key (RDAP needs none). Every other step (government registries, KYC, banking) is a **deterministic mock** whose request/response shape mirrors the real upstream API, so the whole thing runs end-to-end with **zero credentials** and you wire real APIs in incrementally.

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

Every country runs the same 11 logical steps; each maps to different real upstreams:

| # | Step | 🇮🇳 India | 🇵🇭 Philippines | 🇺🇸 United States |
|---|------|----------|----------------|------------------|
| 1 | **Strategy & viability** | **Claude (AI)** 🟢 | **Claude (AI)** 🟢 | **Claude (AI)** 🟢 |
| 2 | Founder KYC | SurePass (PAN + Aadhaar) | HyperVerge (PhilID) | Persona / Middesk |
| 3 | Liabilities & sanctions | trade.gov 🟢 + MCA/GST/CIBIL | trade.gov 🟢 + SEC/BIR/CIC | trade.gov 🟢 + OFAC/UCC/PACER |
| 4 | Name check | MCA RUN | SEC name verification | Secretary of State |
| 5 | **Trademark & domain** | **RDAP** 🟢 + IP India | **RDAP** 🟢 + IPOPHL | **RDAP** 🟢 + USPTO |
| 6 | Entity registration | MCA SPICe+ → **CIN** | SEC eSPARC → **SEC reg no.** | State filing → **filing no.** |
| 7 | Tax registration | Income Tax + GSTN → **PAN/TAN/GSTIN** | BIR → **TIN** | IRS SS-4 → **EIN** |
| 8 | Licenses & registrations | Udyam + DGFT IEC + S&E | Mayor's + Barangay + DTI | State license + FinCEN BOI |
| 9 | Business banking | RazorpayX | UnionBank | Mercury |
| 10 | **Payment gateway** | **Razorpay** 🟢 | **PayMongo** 🟢 | **Stripe** 🟢 |
| 11 | Compliance | EPFO + ESIC | SSS + PhilHealth + Pag-IBIG | Registered agent + state tax |

🟢 = real call when its key is configured (RDAP needs none), otherwise a deterministic mock.

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

Several steps are wired to real services. Add any of these to `.env`
(or export them before `go run`) and that step flips from `mock` to `live`:

| Step | Service | Key(s) | Get them |
|------|---------|--------|----------|
| Strategy | Anthropic / Claude | `ANTHROPIC_API_KEY` (`ANTHROPIC_MODEL`) | console.anthropic.com → API keys |
| Liabilities | trade.gov CSL | `CSL_API_KEY` | Free key at api.data.gov |
| IP / domain | RDAP | — *(none — public)* | Works out of the box |
| Payments 🇮🇳 | Razorpay | `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET` | Razorpay → Settings → API Keys (`rzp_test_…`) |
| Payments 🇺🇸 | Stripe | `STRIPE_SECRET_KEY` | Stripe → Developers → API keys (`sk_test_…`) |
| Payments 🇵🇭 | PayMongo | `PAYMONGO_SECRET_KEY` | PayMongo → Developers → API keys (`sk_test_…`) |

The AI step calls Claude (with a **prompt-cached** system prompt) and returns a
structured viability assessment. Each live payment call creates a real test-mode
object (Razorpay order / Stripe customer / PayMongo link); the IP step does a live
RDAP domain lookup. Returned ids are stored as the step's `external_ref`. Set
`FORCE_MOCK=true` to disable every live call.

> **Wiring a government API for real:** every mock step in
> `backend/internal/providers/{india,ph,us}.go` documents its real upstream in a
> comment. Replace the mock body with an HTTP call (see `payments.go` for the
> pattern) and you're live — the orchestrator, persistence and UI don't change.

---

## ☁️ Deploy serverless

The stack is **serverless-first**. The Go API is a plain HTTP server: on AWS it
runs on **Lambda** via the [Lambda Web Adapter](https://github.com/awslabs/aws-lambda-web-adapter)
(no code changes), and on GCP/Azure it runs as a **scale-to-zero container**.
Pair it with a serverless Postgres (Neon / Aurora Serverless v2 / Cloud SQL / Azure Flexible).

```bash
# 1) serverless Postgres (any cloud) — create at neon.tech, copy the POOLED URL
export DATABASE_URL="postgres://USER:PASS@ep-xxxx-pooler.REGION.aws.neon.tech/biz_launch?sslmode=require"

# 2a) AWS — Lambda via SAM
cd deploy/aws-sam && sam build && sam deploy --guided \
  --parameter-overrides DatabaseUrl="$DATABASE_URL" AnthropicApiKey="$ANTHROPIC_API_KEY"

# 2b) GCP — Cloud Run (scale to zero)
gcloud run deploy biz-launch-api --source ./backend --region us-central1 \
  --allow-unauthenticated --set-env-vars "DATABASE_URL=$DATABASE_URL,DB_MAX_CONNS=4"

# 2c) Azure — Container Apps (scale to zero)
az containerapp up --name biz-launch-api --resource-group biz-launch \
  --ingress external --target-port 8080 --source ./backend --env-vars "DATABASE_URL=$DATABASE_URL"
```

Full walkthrough: **[`deploy/README.md`](deploy/README.md)** and the in-app **Deploy** page (`/deploy`).
On serverless, keep `DB_MAX_CONNS` small (the Lambda image defaults to `2`) and use a pooled DB endpoint.

## 🗂️ Project structure

```
business-launch-orchestrator/
├── docker-compose.yml            # postgres + backend + frontend
├── .env.example                  # all keys (every one optional)
├── deploy/                       # serverless IaC
│   ├── aws-sam/template.yaml     # AWS Lambda (Web Adapter) + Function URL
│   └── README.md                 # AWS / GCP / Azure walkthrough
├── backend/                      # Go orchestration service
│   ├── cmd/server/main.go        # entrypoint, graceful shutdown, auto-migrate
│   ├── Dockerfile                # container (Cloud Run / Container Apps)
│   ├── Dockerfile.lambda         # AWS Lambda image (Web Adapter)
│   ├── migrations/0001_init.sql  # embedded, applied on boot
│   └── internal/
│       ├── domain/               # dependency-free core types
│       ├── store/                # Postgres persistence (pgx)
│       ├── providers/            # per-country adapters + shared clients
│       │   ├── india.go ph.go us.go    # country pipelines (11 steps each)
│       │   ├── payments.go       # Razorpay / Stripe / PayMongo (real sandbox)
│       │   ├── strategy.go       # Claude AI assessment (prompt-cached)
│       │   ├── ip.go             # RDAP domain (live) + trademark search
│       │   └── liabilities.go    # trade.gov sanctions + diligence checks
│       ├── orchestrator/         # the pipeline engine
│       └── api/                  # chi router + handlers
└── frontend/                     # Next.js App Router wizard (standalone build)
    ├── app/page.tsx              # country → details → review wizard
    ├── app/launch/[id]/page.tsx  # live step runner
    ├── app/how-it-works/page.tsx # end-to-end explainer
    ├── app/deploy/page.tsx       # serverless deploy guide
    ├── components/               # Nav, Stepper, StepList, CodeBlock
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
