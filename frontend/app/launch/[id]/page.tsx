"use client";

import { useCallback, useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { StepList } from "@/components/StepList";
import { api, type BusinessDetail } from "@/lib/api";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export default function LaunchPage() {
  const params = useParams<{ id: string }>();
  const id = params.id;

  const [detail, setDetail] = useState<BusinessDetail | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    try {
      setDetail(await api.getBusiness(id));
    } catch (e) {
      setError((e as Error).message);
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  const advanceOne = async () => {
    setBusy(true);
    setError("");
    try {
      await api.advance(id);
      await load();
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setBusy(false);
    }
  };

  // Step through every remaining step client-side so each one animates.
  const runAll = async () => {
    setBusy(true);
    setError("");
    try {
      for (let i = 0; i < 30; i++) {
        const res = await api.advance(id);
        await load();
        if (res.done) break;
        if (res.step && res.step.status === "failed") break;
        await sleep(450);
      }
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setBusy(false);
    }
  };

  if (!detail) {
    return (
      <main className="container">
        {error ? <div className="error">{error}</div> : <p className="muted">Loading…</p>}
        <p style={{ marginTop: 20 }}>
          <Link href="/">← New launch</Link>
        </p>
      </main>
    );
  }

  const { business, steps } = detail;
  const done = steps.filter((s) => s.status === "completed").length;
  const pct = steps.length ? Math.round((done / steps.length) * 100) : 0;
  const allDone = done === steps.length && steps.length > 0;
  const completedRefs = steps.filter((s) => s.external_ref && s.status === "completed");

  return (
    <main className="container">
      <div className="brand">
        <div className="logo">🚀</div>
        <h1>{business.legal_name}</h1>
      </div>
      <p className="subtitle">
        {business.entity_type} · {business.country} ·{" "}
        <span className={`status-tag ${business.status}`}>{business.status.replace("_", " ")}</span>
      </p>

      <div className="progressbar">
        <span style={{ width: `${pct}%` }} />
      </div>
      <p className="muted" style={{ marginTop: -14, marginBottom: 22 }}>
        {done} / {steps.length} steps complete ({pct}%)
      </p>

      {allDone && (
        <div className="card" style={{ marginBottom: 18, borderColor: "var(--green)" }}>
          <h3 style={{ marginTop: 0 }}>🎉 Business launched</h3>
          <p className="muted">All government, banking and payment steps completed. Key references:</p>
          <div className="kv">
            {completedRefs.map((s) => (
              <span key={s.id}>
                {s.title}: <code className="ref">{s.external_ref}</code>
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="card">
        <StepList steps={steps} />
        {error && <div className="error">{error}</div>}
        <div className="actions">
          <Link href="/" className="btn secondary">＋ New launch</Link>
          <div style={{ display: "flex", gap: 12 }}>
            <button className="btn secondary" onClick={advanceOne} disabled={busy || allDone}>
              {busy ? <span className="spin" /> : "▶"} Run next step
            </button>
            <button className="btn" onClick={runAll} disabled={busy || allDone}>
              {busy ? <span className="spin" /> : "⏩"} Run all steps
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
