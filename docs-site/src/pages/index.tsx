import type { ReactNode } from "react";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";

import styles from "./index.module.css";

const PHASES = [
  { icon: "🧭", name: "Decide", steps: ["Strategy & viability (AI)"] },
  { icon: "🛡️", name: "Verify", steps: ["Founder KYC", "Liabilities & sanctions"] },
  { icon: "™️", name: "Brand", steps: ["Name reservation", "Trademark & domain"] },
  { icon: "🏛️", name: "Form", steps: ["Incorporation", "Tax registration", "Licenses & registrations"] },
  { icon: "💳", name: "Operate", steps: ["Bank account", "Payment gateway", "Compliance"] },
];

function Hero() {
  const { siteConfig } = useDocusaurusContext();
  const cf = siteConfig.customFields as { requestAccessUrl: string; githubUrl: string };
  return (
    <header className={styles.hero}>
      <div className={styles.heroInner}>
        <span className={styles.eyebrow}>🇮🇳 India · 🇵🇭 Philippines · 🇺🇸 United States</span>
        <h1 className={styles.title}>Launch a company, end&#8209;to&#8209;end</h1>
        <p className={styles.subtitle}>
          One orchestrated flow for everything it takes to stand up a real operating
          business — AI strategy, KYC, liabilities, IP, incorporation, tax, licenses,
          banking, payments and compliance — across three countries.
        </p>
        <div className={styles.ctaRow}>
          <Link className={styles.ctaPrimary} to={cf.requestAccessUrl}>
            Request early access →
          </Link>
          <Link className={styles.ctaSecondary} to="/docs/overview">
            Read the docs
          </Link>
          <Link className={styles.ctaGhost} to={cf.githubUrl}>
            ★ GitHub
          </Link>
        </div>
        <div className={styles.statRow}>
          <div><b>11</b><span>orchestrated steps</span></div>
          <div><b>3</b><span>jurisdictions</span></div>
          <div><b>0</b><span>credentials to demo</span></div>
        </div>
      </div>
    </header>
  );
}

function FlowStrip() {
  return (
    <section className={styles.section}>
      <div className="container">
        <h2 className={styles.sectionTitle}>The journey — 5 phases, 11 steps</h2>
        <div className={styles.flow}>
          {PHASES.map((p, i) => (
            <div className={styles.phaseCard} key={p.name}>
              <div className={styles.phaseHead}>
                <span className={styles.phaseIcon}>{p.icon}</span>
                <span className={styles.phaseNum}>{i + 1}</span>
              </div>
              <h3>{p.name}</h3>
              <ul>
                {p.steps.map((s) => (
                  <li key={s}>{s}</li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <p className={styles.flowNote}>
          Payments, the AI strategy step and the domain check hit real services;
          government &amp; KYC steps are realistic mocks you swap for live APIs.{" "}
          <Link to="/docs/the-flow">See the full flow →</Link>
        </p>
      </div>
    </section>
  );
}

function FinalCTA() {
  const { siteConfig } = useDocusaurusContext();
  const cf = siteConfig.customFields as { requestAccessUrl: string };
  return (
    <section className={styles.ctaBand}>
      <div className="container">
        <h2>Validating the idea — want in?</h2>
        <p>
          We&apos;re concept-testing which markets and steps matter most. Tell us about
          your launch and we&apos;ll keep you posted.
        </p>
        <Link className={styles.ctaPrimary} to={cf.requestAccessUrl}>
          Request access / share feedback →
        </Link>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title="Launch a company end-to-end"
      description="One orchestrated flow to launch a business in India, the Philippines or the US — strategy, KYC, liabilities, IP, incorporation, tax, banking, payments and compliance."
    >
      <Hero />
      <main>
        <FlowStrip />
        <HomepageFeatures />
        <FinalCTA />
      </main>
    </Layout>
  );
}
