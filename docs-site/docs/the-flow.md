---
sidebar_position: 2
title: The 11-step flow
---

# The 11-step flow

Every launch runs the same eleven logical steps, grouped into five phases. The
orchestrator executes them in order, persisting each result so the launch can be
advanced one step at a time or run straight through.

## Phase 1 — Decide

### 1. Strategy & viability assessment
An AI strategist (Claude) analyses the jurisdiction, the proposed entity type,
the top regulatory & tax risks, and a 90-day go-to-market, then returns a
**0–100 viability score**. This shapes every decision after it.
*Live when an Anthropic key is set; deterministic mock otherwise.*

## Phase 2 — Verify the founder

### 2. Founder KYC
Government-ID identity verification — PAN + Aadhaar (IN), PhilID (PH), or a
Persona/Middesk check (US).

### 3. Liabilities & sanctions screening
Sanctions/PEP match via the **trade.gov Consolidated Screening List**, plus
director-disqualification, litigation, lien and tax-defaulter checks. Surfaces
risk before any money or filing is committed.
*Sanctions list is live when a CSL key is set.*

## Phase 3 — Secure the brand

### 4. Name reservation
Reserve the corporate name with the registry — MCA RUN (IN), SEC (PH),
Secretary of State (US).

### 5. Trademark & domain check
**Live domain availability over public RDAP** (needs no key), plus a trademark
conflict search against IP India / IPOPHL / USPTO. Catches a name clash before it
becomes expensive.

## Phase 4 — Form the entity

### 6. Incorporation
File the company and receive its identifier — **CIN** (IN), **SEC registration
no.** (PH), **state filing no.** (US).

### 7. Tax registration
Obtain tax identifiers — **PAN/TAN + GSTIN** (IN), **TIN** via BIR (PH),
**Federal EIN** (US).

### 8. Licenses & registrations
Sector and statutory registrations beyond incorporation — Udyam / IEC / Shops &
Establishment (IN), Mayor's & Barangay permits (PH), state license + **FinCEN
BOI** report (US).

## Phase 5 — Go operational

### 9. Business bank account
Open a current/checking account — RazorpayX (IN), UnionBank (PH), Mercury (US).

### 10. Payment gateway
Activate online payments against a **real provider sandbox** — Razorpay (IN),
PayMongo (PH), Stripe (US).

### 11. Statutory compliance
Register as an employer / for benefits — EPFO + ESIC (IN), SSS + PhilHealth +
Pag-IBIG (PH), registered agent + state tax (US).

---

Each step records a **status** (`pending → running → completed / failed`), the
provider used, the request, the full response, and an **external reference**
(CIN, EIN, GSTIN, order id, …). See **[Country coverage](./country-coverage.md)**
for the exact provider per country.
