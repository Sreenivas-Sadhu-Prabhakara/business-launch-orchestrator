import type { ReactNode } from "react";
import clsx from "clsx";
import Heading from "@theme/Heading";
import styles from "./styles.module.css";

type FeatureItem = {
  icon: string;
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    icon: "🔀",
    title: "One flow, three countries",
    description: (
      <>The same 11-step pipeline maps to each jurisdiction&apos;s real agencies through pluggable adapters.</>
    ),
  },
  {
    icon: "🧠",
    title: "AI-led strategy",
    description: (
      <>A Claude assessment opens the flow so structure, jurisdiction and risk are decided with evidence.</>
    ),
  },
  {
    icon: "🔌",
    title: "Hybrid integrations",
    description: (
      <>Payments, AI and domain checks hit real services; government/KYC steps are mocks you swap for live APIs.</>
    ),
  },
  {
    icon: "⏯️",
    title: "Resumable state machine",
    description: (
      <>Every step&apos;s request, response and reference is persisted. Advance one step, run all, or resume on failure.</>
    ),
  },
  {
    icon: "☁️",
    title: "Serverless-first",
    description: (
      <>Runs on Lambda (Web Adapter), Cloud Run or Container Apps with a serverless Postgres. Scales to zero.</>
    ),
  },
  {
    icon: "🧱",
    title: "Self-updating UI",
    description: (
      <>The wizard renders whatever the API plan returns — add a step or country and the UI updates itself.</>
    ),
  },
];

function Feature({ icon, title, description }: FeatureItem) {
  return (
    <div className={clsx("col col--4")}>
      <div className={styles.card}>
        <div className={styles.icon}>{icon}</div>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <Heading as="h2" className={styles.sectionTitle}>
          Why it&apos;s built this way
        </Heading>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
