import type { ReactNode } from "react";
import Heading from "@theme/Heading";
import styles from "./styles.module.css";

type FeatureItem = {
  n: string;
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    n: "01",
    title: "One flow, three countries",
    description:
      "The same eleven-step pipeline maps to each jurisdiction's real agencies through pluggable adapters.",
  },
  {
    n: "02",
    title: "AI-led strategy",
    description:
      "A Claude assessment opens the flow, so structure, jurisdiction and risk are decided with evidence.",
  },
  {
    n: "03",
    title: "Hybrid integrations",
    description:
      "Payments, AI and domain checks call real services; government and KYC steps are mocks you swap for live APIs.",
  },
  {
    n: "04",
    title: "Resumable by design",
    description:
      "Every step's request, response and reference is persisted. Advance one step, run all, or resume after a failure.",
  },
  {
    n: "05",
    title: "Serverless-first",
    description:
      "Runs on Lambda, Cloud Run or Container Apps with a serverless Postgres. Scales to zero.",
  },
  {
    n: "06",
    title: "A self-updating interface",
    description:
      "The wizard renders whatever the API plan returns — add a step or a country and the UI follows.",
  },
];

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className={styles.inner}>
        <div className={styles.label}>Principles</div>
        <Heading as="h2" className={styles.title}>
          Considered, not assembled
        </Heading>
        <div className={styles.grid}>
          {FeatureList.map((f) => (
            <div className={styles.card} key={f.n}>
              <div className={styles.num}>{f.n}</div>
              <Heading as="h3" className={styles.cardTitle}>
                {f.title}
              </Heading>
              <p className={styles.cardBody}>{f.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
