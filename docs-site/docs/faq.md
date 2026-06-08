---
sidebar_position: 9
title: FAQ
---

# FAQ

### Does this actually register a company?
Not yet. The payments, AI strategy, sanctions and domain steps make **real**
calls; the government-registry, KYC, banking and licensing steps are realistic
**simulations**. The architecture is built so each can be swapped for a live API
without touching the orchestrator or UI. See **[Integrations](./integrations.md)**.

### Why these three countries?
India, the Philippines and the US span very different incorporation regimes
(MCA SPICe+, SEC eSPARC, US state filings) and payment ecosystems (Razorpay,
PayMongo, Stripe) — a good test of whether one flow can generalise. More
jurisdictions are additive via adapters.

### Can it run with no API keys?
Yes. With zero credentials the whole flow runs end-to-end in mock mode (plus the
keyless live RDAP domain check). Add keys to turn individual steps live.

### What does the AI strategy step do?
It calls Claude with a prompt-cached system prompt and returns a structured
assessment: recommended entity, top risks, a 90-day go-to-market, and a 0–100
viability score. Without a key it returns a deterministic mock assessment.

### Is it production-ready?
It's a working reference implementation, not a regulated filing service. Real
incorporation involves legal documents, signatures, fees and often licensed
intermediaries. Treat this as the integration backbone.

### How is data stored?
In Postgres: one row per launch and one per step, including each step's request,
response, external reference and status — which is what makes launches resumable.

### How do I deploy it?
Serverless on AWS Lambda, GCP Cloud Run or Azure Container Apps with a serverless
Postgres. See **[Deploy serverless](./deploy-serverless.md)**.

### How can I follow along or give feedback?
**[Request access / share feedback](https://tally.so/r/REPLACE_WITH_YOUR_FORM_ID)**
— we're concept-testing which markets and steps matter most.
