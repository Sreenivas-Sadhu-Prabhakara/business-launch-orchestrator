---
sidebar_position: 4
title: Integrations (hybrid)
---

# Integrations — the hybrid model

The system runs end-to-end with **zero credentials**, then lets you turn on real
integrations one at a time. Each step is either **live** (calls a real service) or
**mock** (a deterministic stand-in whose request/response shape mirrors the real
upstream).

## What's live, and with which key

| Step | Service | Key | Notes |
|------|---------|-----|-------|
| Strategy | Anthropic / Claude | `ANTHROPIC_API_KEY` (`ANTHROPIC_MODEL`) | Structured viability assessment, **prompt-cached** system prompt |
| Liabilities | trade.gov CSL | `CSL_API_KEY` | Free key at api.data.gov |
| IP / domain | RDAP | *none* | Public network — **live out of the box** |
| Payments 🇮🇳 | Razorpay | `RAZORPAY_KEY_ID`, `RAZORPAY_KEY_SECRET` | Creates a test-mode order |
| Payments 🇺🇸 | Stripe | `STRIPE_SECRET_KEY` | Creates a test-mode customer |
| Payments 🇵🇭 | PayMongo | `PAYMONGO_SECRET_KEY` | Creates a test-mode payment link |

Everything else (government registries, KYC, banking, licenses) is a documented
mock. Set `FORCE_MOCK=true` to disable **all** live calls (offline demos / CI).

## Why hybrid

- **Demo instantly.** No accounts needed to walk the whole flow.
- **De-risk incrementally.** Wire one real provider at a time and watch its step
  flip from `mock` to `live`.
- **Realistic shapes.** Mock responses match the real API's structure, so swapping
  in a live call is a localized change — the orchestrator, DB and UI don't move.

## Going live on a mock step

Each mock step documents its real upstream in code. To make it real, replace the
mock body with an HTTP call (the payment clients are the reference pattern):

```go
// before: deterministic mock
return domain.StepResult{ExternalRef: mockEIN(b), ...}, nil

// after: real call to the upstream API
resp, err := irs.ApplySS4(ctx, b)   // your live client
if err != nil { return domain.StepResult{}, err }
return domain.StepResult{ExternalRef: resp.EIN, Data: resp.Map()}, nil
```

The return contract (`external_ref` + `data`) is identical, so nothing downstream
changes.
