"use client";

import { useState } from "react";
import Link from "next/link";
import { CodeBlock } from "@/components/CodeBlock";

type Cloud = "aws" | "gcp" | "azure";

const NEON = `# Create a free serverless Postgres at neon.tech, copy the POOLED string:
export DATABASE_URL="postgres://USER:PASS@ep-xxxx-pooler.REGION.aws.neon.tech/biz_launch?sslmode=require"
# (migrations auto-apply on the API's first boot — nothing else to run)`;

const CLOUDS: Record<
  Cloud,
  { label: string; flag: string; compute: string; db: string; steps: { title: string; code: string }[] }
> = {
  aws: {
    label: "AWS",
    flag: "🟧",
    compute: "Lambda (container + Web Adapter) → Function URL",
    db: "Aurora Serverless v2 / Neon",
    steps: [
      {
        title: "Deploy the API as a Lambda (SAM)",
        code: `# prerequisites: AWS SAM CLI + Docker, and \`aws configure\`
cd deploy/aws-sam

sam build
sam deploy --guided \\
  --parameter-overrides \\
    DatabaseUrl="$DATABASE_URL" \\
    AnthropicApiKey="$ANTHROPIC_API_KEY" \\
    StripeSecretKey="$STRIPE_SECRET_KEY"

# stack output "ApiUrl" is your public serverless endpoint`,
      },
      {
        title: "Frontend on Amplify Hosting",
        code: `# point the build at the API URL above
amplify init && amplify add hosting
NEXT_PUBLIC_API_URL="https://<API_URL>" amplify publish`,
      },
    ],
  },
  gcp: {
    label: "Google Cloud",
    flag: "🔵",
    compute: "Cloud Run (scale to zero)",
    db: "Cloud SQL / Neon",
    steps: [
      {
        title: "Deploy the API to Cloud Run",
        code: `# prerequisites: gcloud CLI, \`gcloud auth login\`, project selected
gcloud run deploy biz-launch-api \\
  --source ./backend \\
  --region us-central1 --allow-unauthenticated \\
  --set-env-vars "DATABASE_URL=$DATABASE_URL,DB_MAX_CONNS=4,ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY"`,
      },
      {
        title: "Deploy the frontend to Cloud Run",
        code: `gcloud run deploy biz-launch-web \\
  --source ./frontend --region us-central1 --allow-unauthenticated \\
  --set-build-env-vars "NEXT_PUBLIC_API_URL=https://<API_URL>"`,
      },
    ],
  },
  azure: {
    label: "Azure",
    flag: "🔷",
    compute: "Container Apps (scale to zero)",
    db: "Database for PostgreSQL Flexible / Neon",
    steps: [
      {
        title: "Deploy the API to Container Apps",
        code: `# prerequisites: az CLI, \`az login\`, containerapp extension
az containerapp up \\
  --name biz-launch-api --resource-group biz-launch \\
  --ingress external --target-port 8080 --source ./backend \\
  --env-vars "DATABASE_URL=$DATABASE_URL" "DB_MAX_CONNS=4" "ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY"`,
      },
      {
        title: "Deploy the frontend to Container Apps",
        code: `az containerapp up \\
  --name biz-launch-web --resource-group biz-launch \\
  --ingress external --target-port 3000 --source ./frontend \\
  --env-vars "NEXT_PUBLIC_API_URL=https://<API_URL>"`,
      },
    ],
  },
};

export default function Deploy() {
  const [cloud, setCloud] = useState<Cloud>("aws");
  const c = CLOUDS[cloud];

  return (
    <main className="container">
      <div className="hero">
        <span className="eyebrow">Serverless-first</span>
        <h1>Deploy to any hyperscaler</h1>
        <p>
          The Go API is a plain HTTP server — on AWS it runs on Lambda via the Web
          Adapter (no code changes); on GCP & Azure it runs as a scale-to-zero
          container. Pair it with a serverless Postgres and the whole thing scales to zero.
        </p>
      </div>

      <div className="section-title">Architecture</div>
      <div className="features">
        <div className="feature">
          <div className="ic">🖥️</div>
          <h4>Frontend</h4>
          <p>Next.js (standalone) on Vercel, Amplify, Cloud Run or Container Apps.</p>
        </div>
        <div className="feature">
          <div className="ic">⚙️</div>
          <h4>API</h4>
          <p>Go on Lambda (Web Adapter), Cloud Run or Container Apps. Listens on <code>$PORT</code>.</p>
        </div>
        <div className="feature">
          <div className="ic">🗄️</div>
          <h4>Database</h4>
          <p>Serverless Postgres — Neon, Aurora Serverless v2, Cloud SQL or Azure Flexible.</p>
        </div>
        <div className="feature">
          <div className="ic">🔌</div>
          <h4>Pooling</h4>
          <p>Small per-instance pool (<code>DB_MAX_CONNS=2</code>) behind a pooler (RDS Proxy / Neon).</p>
        </div>
      </div>

      <div className="section-title">Step 1 — a serverless database (any cloud)</div>
      <CodeBlock code={NEON} />

      <div className="section-title">Step 2 — pick your cloud</div>
      <div className="tabs">
        {(Object.keys(CLOUDS) as Cloud[]).map((k) => (
          <button
            key={k}
            className={`tab ${cloud === k ? "active" : ""}`}
            onClick={() => setCloud(k)}
          >
            {CLOUDS[k].flag} {CLOUDS[k].label}
          </button>
        ))}
      </div>

      <div className="card">
        <div className="kv" style={{ marginTop: 0, marginBottom: 8 }}>
          <span>Compute: <code>{c.compute}</code></span>
          <span>Database: <code>{c.db}</code></span>
        </div>
        {c.steps.map((s) => (
          <div key={s.title}>
            <h4 style={{ marginBottom: 0 }}>{s.title}</h4>
            <CodeBlock code={s.code} />
          </div>
        ))}
      </div>

      <div className="section-title">Frontend on Vercel (most serverless)</div>
      <CodeBlock code={`cd frontend
npx vercel --prod   # then set NEXT_PUBLIC_API_URL in the Vercel project settings`} />

      <div className="section-title">In the repo</div>
      <div className="card">
        <div className="kv" style={{ marginTop: 0 }}>
          <span><code>backend/Dockerfile.lambda</code> — AWS Lambda image (Web Adapter)</span>
          <span><code>backend/Dockerfile</code> — Cloud Run / Container Apps</span>
          <span><code>frontend/Dockerfile</code> — Next.js standalone</span>
          <span><code>deploy/aws-sam/template.yaml</code> — SAM stack</span>
          <span><code>deploy/README.md</code> — full guide</span>
        </div>
      </div>

      <div className="hero" style={{ paddingTop: 50 }}>
        <Link href="/" className="btn">🚀 Back to the launcher</Link>
      </div>

      <div className="footer">
        Set <code>FORCE_MOCK=true</code> to run fully offline, or add provider keys to go live.
      </div>
    </main>
  );
}
