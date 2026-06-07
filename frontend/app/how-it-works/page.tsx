import Link from "next/link";

export const metadata = {
  title: "How it works — Business Launch Orchestrator",
  description:
    "The end-to-end journey from idea to operating company across India, the Philippines and the US.",
};

type Step = { icon: string; title: string; desc: string };
type Phase = { icon: string; name: string; steps: Step[] };

const PHASES: Phase[] = [
  {
    icon: "🧭",
    name: "1 · Decide",
    steps: [
      {
        icon: "🤖",
        title: "Strategy & viability assessment",
        desc: "Claude analyses the jurisdiction, the proposed entity type, the top regulatory & tax risks, a 90-day go-to-market, and returns a 0–100 viability score. This shapes every decision after it.",
      },
    ],
  },
  {
    icon: "🛡️",
    name: "2 · Verify the founder",
    steps: [
      {
        icon: "🪪",
        title: "Founder KYC",
        desc: "Government-ID identity verification — PAN + Aadhaar (IN), PhilID (PH), or a Persona/Middesk check (US).",
      },
      {
        icon: "⚖️",
        title: "Liabilities & sanctions screening",
        desc: "Sanctions/PEP match via the trade.gov Consolidated Screening List, plus director-disqualification, litigation, liens and tax-defaulter checks. Surfaces risk before any money or filing is committed.",
      },
    ],
  },
  {
    icon: "™️",
    name: "3 · Secure the brand",
    steps: [
      {
        icon: "📝",
        title: "Name reservation",
        desc: "Reserve the corporate name with the registry — MCA RUN (IN), SEC (PH), Secretary of State (US).",
      },
      {
        icon: "🔎",
        title: "Trademark & domain check",
        desc: "Live domain availability over public RDAP, plus a trademark conflict search against IP India / IPOPHL / USPTO. Catches a name clash before it becomes expensive.",
      },
    ],
  },
  {
    icon: "🏛️",
    name: "4 · Form the entity",
    steps: [
      {
        icon: "🏢",
        title: "Incorporation",
        desc: "File the company and receive its identifier — CIN (IN), SEC registration no. (PH), state filing no. (US).",
      },
      {
        icon: "🧾",
        title: "Tax registration",
        desc: "Obtain tax identifiers — PAN/TAN + GSTIN (IN), TIN via BIR (PH), Federal EIN (US).",
      },
      {
        icon: "📜",
        title: "Licenses & registrations",
        desc: "Sector and statutory registrations beyond incorporation — Udyam/IEC/Shops & Est (IN), Mayor's & Barangay permits (PH), state license + FinCEN BOI report (US).",
      },
    ],
  },
  {
    icon: "💳",
    name: "5 · Go operational",
    steps: [
      {
        icon: "🏦",
        title: "Business bank account",
        desc: "Open a current/checking account — RazorpayX (IN), UnionBank (PH), Mercury (US).",
      },
      {
        icon: "💸",
        title: "Payment gateway",
        desc: "Activate online payments against a real provider sandbox — Razorpay (IN), PayMongo (PH), Stripe (US).",
      },
      {
        icon: "👥",
        title: "Statutory compliance",
        desc: "Register as an employer / for benefits — EPFO + ESIC (IN), SSS + PhilHealth + Pag-IBIG (PH), registered agent + state tax (US).",
      },
    ],
  },
];

const MATRIX: [string, string, string, string][] = [
  ["Strategy", "Claude (AI)", "Claude (AI)", "Claude (AI)"],
  ["KYC", "SurePass", "HyperVerge", "Persona / Middesk"],
  ["Liabilities", "trade.gov + MCA/GST", "trade.gov + SEC/BIR", "trade.gov + OFAC/UCC"],
  ["Name", "MCA RUN", "SEC", "Secretary of State"],
  ["IP", "RDAP + IP India", "RDAP + IPOPHL", "RDAP + USPTO"],
  ["Incorporation", "MCA SPICe+ → CIN", "SEC eSPARC", "State filing"],
  ["Tax", "PAN/TAN/GSTIN", "BIR → TIN", "IRS → EIN"],
  ["Registrations", "Udyam + IEC", "Mayor's permit", "License + BOI"],
  ["Banking", "RazorpayX", "UnionBank", "Mercury"],
  ["Payments", "Razorpay 🟢", "PayMongo 🟢", "Stripe 🟢"],
  ["Compliance", "EPFO + ESIC", "SSS + PhilHealth", "Agent + state tax"],
];

const FEATURES: Step[] = [
  { icon: "🔀", title: "One flow, three countries", desc: "The same 11-step pipeline maps to each jurisdiction's real agencies through pluggable adapters." },
  { icon: "🧠", title: "AI-led strategy", desc: "A Claude assessment opens the flow so structure, jurisdiction and risk are decided with evidence." },
  { icon: "🔌", title: "Hybrid integrations", desc: "Payments + domain checks hit real provider sandboxes; government/KYC steps are realistic mocks you swap for live APIs." },
  { icon: "⏯️", title: "Resumable state machine", desc: "Every step's request, response and reference is persisted. Advance one step or run all — and resume after a failure." },
  { icon: "☁️", title: "Serverless-first", desc: "Runs on Lambda (Web Adapter), Cloud Run or Container Apps with a serverless Postgres. Scales to zero." },
  { icon: "🧱", title: "Typed & data-driven UI", desc: "The wizard renders whatever the API plan returns — add a step or country and the UI updates itself." },
];

export default function HowItWorks() {
  return (
    <main className="container">
      <div className="hero">
        <span className="eyebrow">From idea to operating company</span>
        <h1>How a business gets launched, end to end</h1>
        <p>
          Launching a company is really a chain of dependent decisions and filings across
          half a dozen agencies and vendors. This orchestrates all of them into one
          resumable pipeline — for 🇮🇳 India, 🇵🇭 the Philippines and 🇺🇸 the US.
        </p>
        <div className="cta">
          <Link href="/" className="btn">🚀 Start a launch</Link>
          <Link href="/deploy" className="btn secondary">☁️ Deploy it</Link>
        </div>
      </div>

      <div className="section-title">The journey — 5 phases, 11 steps</div>
      {PHASES.map((p) => (
        <div className="phase" key={p.name}>
          <div className="phase-head">
            <span className="dot">{p.icon}</span>
            <h3>{p.name}</h3>
          </div>
          <div className="timeline">
            {p.steps.map((s) => (
              <div className="tl-item" key={s.title}>
                <div className="t">{s.icon}&nbsp; {s.title}</div>
                <div className="d">{s.desc}</div>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div className="section-title">Same step, different country</div>
      <div style={{ overflowX: "auto" }}>
        <table className="tbl">
          <thead>
            <tr>
              <th>Step</th>
              <th>🇮🇳 India</th>
              <th>🇵🇭 Philippines</th>
              <th>🇺🇸 United States</th>
            </tr>
          </thead>
          <tbody>
            {MATRIX.map((r) => (
              <tr key={r[0]}>
                <td>{r[0]}</td>
                <td>{r[1]}</td>
                <td>{r[2]}</td>
                <td>{r[3]}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="muted" style={{ marginTop: 10 }}>
        🟢 hits a real provider sandbox when a test key is set; otherwise a deterministic mock runs.
      </p>

      <div className="section-title">Why it's built this way</div>
      <div className="features">
        {FEATURES.map((f) => (
          <div className="feature" key={f.title}>
            <div className="ic">{f.icon}</div>
            <h4>{f.title}</h4>
            <p>{f.desc}</p>
          </div>
        ))}
      </div>

      <div className="section-title">Under the hood</div>
      <div className="card">
        <div className="features">
          <div className="feature">
            <div className="ic">🖥️</div>
            <h4>Next.js wizard</h4>
            <p>Country → details → review → run. Renders the live plan and step status from the API.</p>
          </div>
          <div className="feature">
            <div className="ic">⚙️</div>
            <h4>Go orchestrator</h4>
            <p>A chi API + a pipeline engine that runs each step against the right country adapter and records the result.</p>
          </div>
          <div className="feature">
            <div className="ic">🔗</div>
            <h4>Provider adapters</h4>
            <p>One per country, each calling KYC, registry, tax, banking, payment and compliance providers.</p>
          </div>
          <div className="feature">
            <div className="ic">🗄️</div>
            <h4>Postgres</h4>
            <p>Stores every business and step — request, response, external reference and status — so launches resume.</p>
          </div>
        </div>
      </div>

      <div className="hero" style={{ paddingTop: 50 }}>
        <h1 style={{ fontSize: 26 }}>Ready to try it?</h1>
        <div className="cta">
          <Link href="/" className="btn">🚀 Launch a business</Link>
        </div>
      </div>

      <div className="footer">
        Reference implementation — government/KYC/banking steps are simulations. See the README disclaimer.
      </div>
    </main>
  );
}
