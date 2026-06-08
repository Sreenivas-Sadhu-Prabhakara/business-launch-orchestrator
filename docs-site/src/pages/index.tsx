import { useEffect, useState } from "react";
import type { ReactNode, CSSProperties } from "react";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";
import { MARKETS, marketByCode, type Market, type MarketCode } from "@site/src/data/markets";

import styles from "./index.module.css";

function delay(i: number): CSSProperties {
  return { animationDelay: `${0.12 + i * 0.09}s` };
}

function MarketSelect({
  market,
  onChange,
}: {
  market: MarketCode;
  onChange: (c: MarketCode) => void;
}) {
  return (
    <div className={styles.marketSelect} role="tablist" aria-label="Choose a market">
      {MARKETS.map((m) => (
        <button
          key={m.code}
          role="tab"
          aria-selected={market === m.code}
          className={`${styles.marketOpt} ${market === m.code ? styles.marketActive : ""}`}
          onClick={() => onChange(m.code)}
        >
          {m.name}
        </button>
      ))}
    </div>
  );
}

function Hero({
  market,
  onChange,
}: {
  market: MarketCode;
  onChange: (c: MarketCode) => void;
}) {
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
        <div className={`${styles.reveal}`} style={delay(3)}>
          <div className={styles.selectLabel}>Choose your market</div>
          <MarketSelect market={market} onChange={onChange} />
        </div>
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

function Flow({ market }: { market: Market }) {
  return (
    <section className={styles.section} id="flow">
      <div className={styles.sectionInner}>
        <div className={styles.sectionLabel}>The journey · {market.name}</div>
        <h2 className={styles.sectionTitle}>
          Forming {indefinite(market.entity)} <em>{market.entity}</em>
        </h2>
        <p className={styles.flowSub}>
          Eleven steps, end to end — with {market.name}&apos;s actual agencies and the
          identifiers each step produces.
        </p>

        {market.phases.map((phase) => (
          <div className={styles.phaseBlock} key={phase.label}>
            <div className={styles.phaseLabel}>{phase.label}</div>
            {phase.steps.map((s) => (
              <div className={styles.stepRow} key={s.n}>
                <div className={styles.stepNum}>
                  {String(s.n).padStart(2, "0")}
                </div>
                <div className={styles.stepMain}>
                  <h3 className={styles.stepTitle}>{s.title}</h3>
                  <div className={styles.stepProvider}>{s.provider}</div>
                </div>
                <div className={styles.stepMeta}>
                  {s.output ? <span className={styles.stepOut}>→ {s.output}</span> : null}
                  {s.live ? <span className={styles.liveDot}>live</span> : null}
                </div>
              </div>
            ))}
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

function indefinite(word: string): string {
  return /^[aeiou]/i.test(word) ? "an" : "a";
}

export default function Home(): ReactNode {
  const [code, setCode] = useState<MarketCode>("IN");

  // Read ?market= on mount; keep the URL in sync so a market is shareable.
  useEffect(() => {
    const param = new URLSearchParams(window.location.search).get("market");
    if (param && MARKETS.some((m) => m.code === (param.toUpperCase() as MarketCode))) {
      setCode(param.toUpperCase() as MarketCode);
    }
  }, []);

  const onChange = (c: MarketCode) => {
    setCode(c);
    const url = new URL(window.location.href);
    url.searchParams.set("market", c);
    window.history.replaceState({}, "", url);
  };

  const market = marketByCode(code);

  return (
    <Layout
      title="Launch a company end-to-end"
      description="One orchestrated flow to launch a business in India, the Philippines or the United States — strategy, KYC, liabilities, IP, incorporation, tax, banking, payments and compliance."
    >
      <Hero market={code} onChange={onChange} />
      <main>
        <Flow market={market} />
        <HomepageFeatures />
        <Closing />
      </main>
    </Layout>
  );
}
