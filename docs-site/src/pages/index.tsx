import type { ReactNode, CSSProperties } from "react";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";

import styles from "./index.module.css";

const PHASES = [
  { n: "01", name: "Decide", steps: "Strategy & viability assessment, AI-led" },
  { n: "02", name: "Verify", steps: "Founder KYC · Liabilities & sanctions screening" },
  { n: "03", name: "Secure the brand", steps: "Name reservation · Trademark & domain clearance" },
  { n: "04", name: "Form the entity", steps: "Incorporation · Tax registration · Licenses & registrations" },
  { n: "05", name: "Go operational", steps: "Business banking · Payment gateway · Compliance" },
];

function delay(i: number): CSSProperties {
  return { animationDelay: `${0.12 + i * 0.09}s` };
}

function Hero() {
  const { siteConfig } = useDocusaurusContext();
  const cf = siteConfig.customFields as { requestAccessUrl: string };
  return (
    <header className={styles.hero}>
      <div className={styles.heroInner}>
        <span className={`${styles.kicker} ${styles.reveal}`} style={delay(0)}>
          Incorporation · Tax · Banking · Payments
        </span>
        <h1 className={`${styles.title} ${styles.reveal}`} style={delay(1)}>
          Launch a company, <em>end&nbsp;to&nbsp;end</em>
        </h1>
        <p className={`${styles.lede} ${styles.reveal}`} style={delay(2)}>
          One considered flow for everything it takes to stand up a real operating
          business — strategy, diligence, incorporation, tax, banking, payments and
          compliance, orchestrated as a single resumable pipeline.
        </p>
        <p className={`${styles.countries} ${styles.reveal}`} style={delay(3)}>
          India — Philippines — United States
        </p>
        <div className={`${styles.ctaRow} ${styles.reveal}`} style={delay(4)}>
          <Link className={styles.ctaPrimary} to={cf.requestAccessUrl}>
            Request access
          </Link>
          <Link className={styles.ctaText} to="/docs/overview">
            Read the documentation →
          </Link>
        </div>
        <div className={`${styles.stats} ${styles.reveal}`} style={delay(5)}>
          <div className={styles.stat}>
            <div className={styles.statNum}>11</div>
            <span className={styles.statLabel}>Orchestrated steps</span>
          </div>
          <div className={styles.stat}>
            <div className={styles.statNum}>3</div>
            <span className={styles.statLabel}>Jurisdictions</span>
          </div>
          <div className={styles.stat}>
            <div className={styles.statNum}>0</div>
            <span className={styles.statLabel}>Credentials to demo</span>
          </div>
        </div>
      </div>
    </header>
  );
}

function Flow() {
  return (
    <section className={styles.section}>
      <div className={styles.sectionInner}>
        <div className={styles.sectionLabel}>The journey</div>
        <h2 className={styles.sectionTitle}>Five phases, eleven steps</h2>
        {PHASES.map((p) => (
          <div className={styles.phase} key={p.n}>
            <div className={styles.phaseNum}>{p.n}</div>
            <div>
              <h3 className={styles.phaseName}>{p.name}</h3>
              <div className={styles.phaseSteps}>{p.steps}</div>
            </div>
          </div>
        ))}
        <p className={styles.flowNote}>
          Payments, the AI strategy step and the domain check call real services;
          government and KYC steps are faithful mocks you replace with live APIs, one
          at a time. <Link to="/docs/the-flow">Explore the full flow →</Link>
        </p>
      </div>
    </section>
  );
}

function Closing() {
  const { siteConfig } = useDocusaurusContext();
  const cf = siteConfig.customFields as { requestAccessUrl: string };
  return (
    <section className={styles.closing}>
      <h2 className={styles.closingTitle}>We&apos;re validating the idea.</h2>
      <p className={styles.closingText}>
        Tell us about your launch — which market, which steps slow you down most — and
        we&apos;ll keep you close to what we build next.
      </p>
      <Link className={styles.ctaPrimary} to={cf.requestAccessUrl}>
        Request access &amp; share feedback
      </Link>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="Launch a company end-to-end"
      description="One orchestrated flow to launch a business in India, the Philippines or the United States — strategy, KYC, liabilities, IP, incorporation, tax, banking, payments and compliance."
    >
      <Hero />
      <main>
        <Flow />
        <HomepageFeatures />
        <Closing />
      </main>
    </Layout>
  );
}
