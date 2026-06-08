// Per-market end-to-end flow. Mirrors the backend country adapters so the
// landing's flow matches what the orchestrator actually runs.

export type MarketCode = "IN" | "PH" | "US";

export interface FlowStep {
  n: number;
  title: string;
  provider: string;
  output?: string; // the identifier / artifact this step produces
  live?: boolean; // calls a real service when keyed
}

export interface Phase {
  label: string; // "01 — Decide"
  steps: FlowStep[];
}

export interface Market {
  code: MarketCode;
  name: string;
  short: string;
  entity: string;
  phases: Phase[];
}

export const MARKETS: Market[] = [
  {
    code: "IN",
    name: "India",
    short: "India",
    entity: "Private Limited Company",
    phases: [
      {
        label: "01 — Decide",
        steps: [
          { n: 1, title: "Strategy & viability", provider: "Claude — AI strategist", output: "viability score", live: true },
        ],
      },
      {
        label: "02 — Verify",
        steps: [
          { n: 2, title: "Founder KYC", provider: "SurePass — PAN + Aadhaar" },
          { n: 3, title: "Liabilities & sanctions", provider: "trade.gov CSL + MCA / GST / CIBIL", live: true },
        ],
      },
      {
        label: "03 — Secure the brand",
        steps: [
          { n: 4, title: "Name reservation", provider: "MCA RUN" },
          { n: 5, title: "Trademark & domain", provider: "RDAP + IP India", live: true },
        ],
      },
      {
        label: "04 — Form the entity",
        steps: [
          { n: 6, title: "Incorporation", provider: "MCA SPICe+ (INC-32)", output: "CIN" },
          { n: 7, title: "Tax registration", provider: "Income Tax + GSTN", output: "PAN · TAN · GSTIN" },
          { n: 8, title: "Licenses & registrations", provider: "Udyam + DGFT IEC + Shops & Est." },
        ],
      },
      {
        label: "05 — Go operational",
        steps: [
          { n: 9, title: "Business banking", provider: "RazorpayX", output: "current account" },
          { n: 10, title: "Payment gateway", provider: "Razorpay", live: true },
          { n: 11, title: "Compliance", provider: "EPFO + ESIC", output: "PF · ESI" },
        ],
      },
    ],
  },
  {
    code: "PH",
    name: "Philippines",
    short: "Philippines",
    entity: "Domestic Corporation",
    phases: [
      {
        label: "01 — Decide",
        steps: [
          { n: 1, title: "Strategy & viability", provider: "Claude — AI strategist", output: "viability score", live: true },
        ],
      },
      {
        label: "02 — Verify",
        steps: [
          { n: 2, title: "Founder KYC", provider: "HyperVerge — PhilID" },
          { n: 3, title: "Liabilities & sanctions", provider: "trade.gov CSL + SEC / BIR / CIC", live: true },
        ],
      },
      {
        label: "03 — Secure the brand",
        steps: [
          { n: 4, title: "Name verification", provider: "SEC name verification" },
          { n: 5, title: "Trademark & domain", provider: "RDAP + IPOPHL", live: true },
        ],
      },
      {
        label: "04 — Form the entity",
        steps: [
          { n: 6, title: "Incorporation", provider: "SEC eSPARC / OneSEC", output: "SEC reg. no." },
          { n: 7, title: "Tax registration", provider: "BIR (Form 2303)", output: "TIN" },
          { n: 8, title: "Permits & registrations", provider: "Mayor's permit + Barangay + DTI" },
        ],
      },
      {
        label: "05 — Go operational",
        steps: [
          { n: 9, title: "Business banking", provider: "UnionBank", output: "corporate account" },
          { n: 10, title: "Payment gateway", provider: "PayMongo", live: true },
          { n: 11, title: "Compliance", provider: "SSS + PhilHealth + Pag-IBIG" },
        ],
      },
    ],
  },
  {
    code: "US",
    name: "United States",
    short: "United States",
    entity: "LLC",
    phases: [
      {
        label: "01 — Decide",
        steps: [
          { n: 1, title: "Strategy & viability", provider: "Claude — AI strategist", output: "viability score", live: true },
        ],
      },
      {
        label: "02 — Verify",
        steps: [
          { n: 2, title: "Founder KYC", provider: "Persona / Middesk" },
          { n: 3, title: "Liabilities & sanctions", provider: "trade.gov CSL + OFAC / UCC / PACER", live: true },
        ],
      },
      {
        label: "03 — Secure the brand",
        steps: [
          { n: 4, title: "Name availability", provider: "Secretary of State" },
          { n: 5, title: "Trademark & domain", provider: "RDAP + USPTO", live: true },
        ],
      },
      {
        label: "04 — Form the entity",
        steps: [
          { n: 6, title: "Incorporation", provider: "State filing (Middesk Agents)", output: "filing no." },
          { n: 7, title: "Tax registration", provider: "IRS (Form SS-4)", output: "EIN" },
          { n: 8, title: "Licenses & registrations", provider: "State license + FinCEN BOI" },
        ],
      },
      {
        label: "05 — Go operational",
        steps: [
          { n: 9, title: "Business banking", provider: "Mercury", output: "checking account" },
          { n: 10, title: "Payment gateway", provider: "Stripe", live: true },
          { n: 11, title: "Compliance", provider: "Registered agent + state tax" },
        ],
      },
    ],
  },
];

export function marketByCode(code: string | null | undefined): Market {
  return MARKETS.find((m) => m.code === code) ?? MARKETS[0];
}
