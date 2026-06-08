---
sidebar_position: 1
title: Overview
slug: /overview
---

# Overview

**Business Launch Orchestrator** turns "I want to start a company" into a single,
guided, end-to-end flow — across 🇮🇳 **India**, 🇵🇭 **the Philippines** and 🇺🇸 **the United States**.

Launching a company is really a chain of dependent decisions and filings spread
across half a dozen government agencies and private vendors: decide the structure,
verify the founder, clear the brand, incorporate, register for tax, get licenses,
open banking, switch on payments, and stay compliant. This product runs all of
them as **one resumable pipeline** with a clear status for every step.

## What you get

- **One flow, three countries.** The same 11-step pipeline maps to each
  jurisdiction's real agencies through pluggable adapters.
- **AI-led strategy.** A Claude assessment opens the flow, so entity type,
  jurisdiction and risk are decided with evidence — not guesswork.
- **Real diligence built in.** Sanctions/liabilities screening and an IP
  (trademark + domain) check happen *before* you spend money on filings.
- **Hybrid integrations.** Payments, the AI step, and domain checks hit real
  services; government/KYC/banking steps are realistic mocks you swap for live
  APIs one at a time.
- **Resumable.** Every step's request, response and reference is persisted —
  advance one step, run them all, or resume after a failure.
- **Serverless-first.** Runs on AWS Lambda, GCP Cloud Run or Azure Container
  Apps with a serverless Postgres. Scales to zero.

## The pipeline at a glance

```
🧭 Strategy → 🪪 KYC → ⚖️ Liabilities → 📝 Name → 🔎 IP →
🏢 Incorporation → 🧾 Tax → 📜 Registrations → 🏦 Banking → 💸 Payments → 👥 Compliance
```

See **[The 11-step flow](./the-flow.md)** for what each step does, and
**[Country coverage](./country-coverage.md)** for how it maps to each jurisdiction.

## Who it's for

Founders, incorporation/fintech platforms, accelerators and "company-in-a-box"
services that need a single API + UI to stand up a compliant operating entity in
multiple markets.

:::note Concept testing
This site documents the concept and the working reference implementation. If the
problem resonates, **[request access](https://tally.so/r/REPLACE_WITH_YOUR_FORM_ID)**
— we're validating which markets and steps matter most.
:::

:::warning Reference implementation
The government-registry, KYC and banking steps are **simulations** — generated
identifiers are realistically formatted but are not real registrations. See the
**[Disclaimer](./disclaimer.md)**.
:::
