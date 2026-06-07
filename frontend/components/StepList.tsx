"use client";

import type { LaunchStep, PlannedStep } from "@/lib/api";

function ModeBadge({ mode }: { mode: "live" | "mock" }) {
  return <span className={`badge ${mode}`}>{mode}</span>;
}

/** Renders the read-only plan preview (before a launch is created). */
export function PlanPreview({ plan }: { plan: PlannedStep[] }) {
  return (
    <div>
      {plan.map((s) => (
        <div className="plan-step" key={s.seq}>
          <div className="seq">{s.seq}</div>
          <div className="grow">
            <div className="title">{s.title}</div>
            <div className="provider">{s.provider}</div>
          </div>
          <ModeBadge mode={s.mode} />
        </div>
      ))}
    </div>
  );
}

function StatusIcon({ status }: { status: LaunchStep["status"] }) {
  if (status === "completed") return <>✓</>;
  if (status === "failed") return <>✕</>;
  if (status === "running") return <span className="spin" />;
  return <>•</>;
}

/** Renders live launch steps with status, refs and any error. */
export function StepList({ steps }: { steps: LaunchStep[] }) {
  return (
    <div>
      {steps.map((s) => (
        <div className={`plan-step ${s.status}`} key={s.id}>
          <div className="seq">
            {s.status === "running" ? <span className="spin" /> : <StatusIcon status={s.status} />}
          </div>
          <div className="grow">
            <div className="title">{s.title}</div>
            <div className="provider">{s.provider}</div>
            {s.error ? <div className="error">{s.error}</div> : null}
          </div>
          <div style={{ textAlign: "right" }}>
            {s.external_ref ? <span className="ref">{s.external_ref}</span> : null}
            <div style={{ marginTop: 6 }}>
              <ModeBadge mode={s.mode} />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
