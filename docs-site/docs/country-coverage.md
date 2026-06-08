---
sidebar_position: 3
title: Country coverage
---

# Country coverage

The same step means different concrete API calls in each country. The orchestrator
selects the right **adapter** by country and runs that jurisdiction's providers.

| # | Step | 🇮🇳 India | 🇵🇭 Philippines | 🇺🇸 United States |
|---|------|----------|----------------|------------------|
| 1 | Strategy & viability | Claude (AI) 🟢 | Claude (AI) 🟢 | Claude (AI) 🟢 |
| 2 | Founder KYC | SurePass (PAN + Aadhaar) | HyperVerge (PhilID) | Persona / Middesk |
| 3 | Liabilities & sanctions | trade.gov 🟢 + MCA/GST/CIBIL | trade.gov 🟢 + SEC/BIR/CIC | trade.gov 🟢 + OFAC/UCC/PACER |
| 4 | Name check | MCA RUN | SEC name verification | Secretary of State |
| 5 | Trademark & domain | RDAP 🟢 + IP India | RDAP 🟢 + IPOPHL | RDAP 🟢 + USPTO |
| 6 | Incorporation | MCA SPICe+ → **CIN** | SEC eSPARC → **SEC reg no.** | State filing → **filing no.** |
| 7 | Tax registration | Income Tax + GSTN → **PAN/TAN/GSTIN** | BIR → **TIN** | IRS SS-4 → **EIN** |
| 8 | Licenses & registrations | Udyam + DGFT IEC + S&E | Mayor's + Barangay + DTI | State license + FinCEN BOI |
| 9 | Business banking | RazorpayX | UnionBank | Mercury |
| 10 | Payment gateway | Razorpay 🟢 | PayMongo 🟢 | Stripe 🟢 |
| 11 | Compliance | EPFO + ESIC | SSS + PhilHealth + Pag-IBIG | Registered agent + state tax |

🟢 = makes a real call when its key is configured (RDAP needs none); otherwise a
deterministic mock runs.

## Entity types offered

| Country | Entity types |
|---------|--------------|
| 🇮🇳 India | Private Limited Company · LLP · One Person Company · Sole Proprietorship |
| 🇵🇭 Philippines | Domestic Corporation · One Person Corporation · Partnership · Sole Proprietorship |
| 🇺🇸 United States | LLC · C-Corp · S-Corp |

## Adding a country

The system is adapter-based, so a new jurisdiction is additive:

1. Implement the `Adapter` interface (`Country()`, `Plan()`, `Execute()`).
2. Register it in the provider registry.
3. Add it to the country validity check.

The UI renders whatever the API plan returns, so the wizard updates itself. See
**[Architecture](./architecture.md)**.
