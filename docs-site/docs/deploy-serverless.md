---
sidebar_position: 8
title: Deploy serverless
---

# Deploy serverless

The stack is **serverless-first**. The Go API is a plain HTTP server: on AWS it
runs on **Lambda** via the [Lambda Web Adapter](https://github.com/awslabs/aws-lambda-web-adapter)
(no code changes), and on GCP/Azure it runs as a **scale-to-zero container**.

| Layer | AWS | GCP | Azure |
|-------|-----|-----|-------|
| API (Go) | Lambda (container + Web Adapter) → Function URL | Cloud Run | Container Apps |
| Frontend (Next.js) | Amplify / Lambda | Cloud Run (standalone) | Static Web Apps / Container Apps |
| Postgres | Aurora Serverless v2 / Neon | Cloud SQL / Neon | Flexible Server / Neon |

:::tip Serverless DB rule
Every warm instance opens its own pool — keep `DB_MAX_CONNS` small (the Lambda
image defaults to `2`) and connect through a pooled endpoint. **Neon** (free tier,
built-in pooler) is the fastest way to satisfy this on any cloud.
:::

## 0. A serverless Postgres (any cloud)

Create a free database at **neon.tech**, copy the **pooled** connection string:

```bash
export DATABASE_URL="postgres://USER:PASS@ep-xxxx-pooler.REGION.aws.neon.tech/biz_launch?sslmode=require"
```

## AWS — Lambda (SAM)

```bash
cd deploy/aws-sam
sam build
sam deploy --guided \
  --parameter-overrides \
    DatabaseUrl="$DATABASE_URL" \
    AnthropicApiKey="$ANTHROPIC_API_KEY" \
    StripeSecretKey="$STRIPE_SECRET_KEY"
# stack output "ApiUrl" is your public serverless endpoint
```

## GCP — Cloud Run

```bash
gcloud run deploy biz-launch-api \
  --source ./backend --region us-central1 --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=$DATABASE_URL,DB_MAX_CONNS=4,ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY"

gcloud run deploy biz-launch-web \
  --source ./frontend --region us-central1 --allow-unauthenticated \
  --set-build-env-vars "NEXT_PUBLIC_API_URL=https://<API_URL>"
```

## Azure — Container Apps

```bash
az containerapp up \
  --name biz-launch-api --resource-group biz-launch \
  --ingress external --target-port 8080 --source ./backend \
  --env-vars "DATABASE_URL=$DATABASE_URL" "DB_MAX_CONNS=4"
```

## Frontend on Vercel

```bash
cd frontend
npx vercel --prod   # set NEXT_PUBLIC_API_URL in the Vercel project settings
```

## In the repo

- `backend/Dockerfile.lambda` — AWS Lambda image (Web Adapter)
- `backend/Dockerfile` — Cloud Run / Container Apps
- `frontend/Dockerfile` — Next.js standalone
- `deploy/aws-sam/template.yaml` — SAM stack
- `deploy/README.md` — full walkthrough
