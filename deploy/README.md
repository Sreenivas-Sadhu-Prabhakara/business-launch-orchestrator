# ☁️ Serverless deployment

Everything here is built to run **serverless-first** — scale-to-zero compute and
a serverless Postgres. The Go API is a normal HTTP server; on AWS it runs on
Lambda via the [Lambda Web Adapter](https://github.com/awslabs/aws-lambda-web-adapter)
(no code changes), and on GCP/Azure it runs as a scale-to-zero container.

| Layer | AWS | GCP | Azure |
|-------|-----|-----|-------|
| API (Go) | **Lambda** (container + Web Adapter) → Function URL | **Cloud Run** | **Container Apps** |
| Frontend (Next.js) | Amplify Hosting / Lambda | Cloud Run (standalone) | Static Web Apps / Container Apps |
| Postgres | **Aurora Serverless v2** / Neon | **Cloud SQL** / Neon | **Flexible Server** / Neon |
| Pooling | RDS Proxy | Neon pooler | PgBouncer / Neon |

> **Serverless DB rule:** every warm instance opens its own pool, so set
> `DB_MAX_CONNS=2` (already the default in the Lambda image) and connect through a
> pooled endpoint. A single shared serverless Postgres like **Neon** (free tier,
> built-in pooler) is the fastest way to satisfy this on any cloud.

---

## 0. A serverless Postgres in 1 minute (works for all clouds)

Create a free database at **neon.tech**, copy the **pooled** connection string,
and use it as `DATABASE_URL` below:

```bash
export DATABASE_URL="postgres://USER:PASSWORD@ep-xxxx-pooler.REGION.aws.neon.tech/biz_launch?sslmode=require"
```

The API auto-applies migrations on first boot — nothing else to run.

---

## 1. AWS — Lambda (SAM)

```bash
# prerequisites: AWS SAM CLI + Docker, and `aws configure` done
cd deploy/aws-sam

sam build
sam deploy --guided \
  --parameter-overrides \
    DatabaseUrl="$DATABASE_URL" \
    AnthropicApiKey="$ANTHROPIC_API_KEY" \
    StripeSecretKey="$STRIPE_SECRET_KEY"

# The stack output `ApiUrl` is your public, serverless API endpoint.
```

> Building on an Apple-silicon Mac produces an `arm64` image (matches the
> template). On an x86 host, change `Architectures` to `x86_64` in `template.yaml`.

## 2. GCP — Cloud Run (scale to zero)

```bash
# prerequisites: gcloud CLI, `gcloud auth login`, a project selected
gcloud run deploy biz-launch-api \
  --source ./backend \
  --region us-central1 --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=$DATABASE_URL,DB_MAX_CONNS=4,ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY"

# Frontend (point it at the API URL printed above):
gcloud run deploy biz-launch-web \
  --source ./frontend --region us-central1 --allow-unauthenticated \
  --set-build-env-vars "NEXT_PUBLIC_API_URL=https://<API_URL>"
```

## 3. Azure — Container Apps (scale to zero)

```bash
# prerequisites: az CLI, `az login`, containerapp extension
az containerapp up \
  --name biz-launch-api --resource-group biz-launch --ingress external --target-port 8080 \
  --source ./backend \
  --env-vars "DATABASE_URL=$DATABASE_URL" "DB_MAX_CONNS=4" "ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY"

az containerapp up \
  --name biz-launch-web --resource-group biz-launch --ingress external --target-port 3000 \
  --source ./frontend --env-vars "NEXT_PUBLIC_API_URL=https://<API_URL>"
```

## 4. Frontend on Vercel (most serverless)

```bash
cd frontend
npx vercel --prod    # set NEXT_PUBLIC_API_URL in the Vercel project settings
```

---

## Files

- `../backend/Dockerfile.lambda` — AWS Lambda image (Web Adapter)
- `../backend/Dockerfile` — generic container (Cloud Run / Container Apps / App Runner)
- `../frontend/Dockerfile` — Next.js standalone container
- `aws-sam/template.yaml` — AWS SAM stack (Lambda + Function URL)
