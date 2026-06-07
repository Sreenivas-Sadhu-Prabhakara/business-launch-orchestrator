"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Stepper } from "@/components/Stepper";
import { PlanPreview } from "@/components/StepList";
import {
  api,
  type CountryCode,
  type CountryInfo,
  type CreateBusinessInput,
} from "@/lib/api";

const FLAGS: Record<CountryCode, string> = { IN: "🇮🇳", PH: "🇵🇭", US: "🇺🇸" };

const ENTITY_TYPES: Record<CountryCode, string[]> = {
  IN: ["Private Limited Company", "LLP", "One Person Company", "Sole Proprietorship"],
  US: ["LLC", "C-Corp", "S-Corp"],
  PH: ["Domestic Corporation", "One Person Corporation", "Partnership", "Sole Proprietorship"],
};

const WIZARD = ["Jurisdiction", "Founder & company", "Review & launch"];

export default function Home() {
  const router = useRouter();
  const [phase, setPhase] = useState(0);
  const [countries, setCountries] = useState<CountryInfo[]>([]);
  const [country, setCountry] = useState<CountryCode | null>(null);
  const [entityType, setEntityType] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    legal_name: "",
    founder_name: "",
    founder_email: "",
    founder_phone: "",
    founder_id_number: "",
    line1: "",
    city: "",
    state: "",
    postal_code: "",
  });

  useEffect(() => {
    api.listCountries().then(setCountries).catch((e) => setError(String(e.message)));
  }, []);

  const selectCountry = (c: CountryCode) => {
    setCountry(c);
    setEntityType(ENTITY_TYPES[c][0]);
  };

  const plan = countries.find((c) => c.code === country)?.plan ?? [];
  const set = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [k]: e.target.value });

  const launch = async () => {
    if (!country) return;
    setSubmitting(true);
    setError("");
    const payload: CreateBusinessInput = {
      country,
      entity_type: entityType,
      legal_name: form.legal_name,
      founder_name: form.founder_name,
      founder_email: form.founder_email,
      founder_phone: form.founder_phone,
      founder_id_number: form.founder_id_number,
      address: {
        line1: form.line1,
        city: form.city,
        state: form.state,
        postal_code: form.postal_code,
        country: country,
      },
    };
    try {
      const detail = await api.createBusiness(payload);
      router.push(`/launch/${detail.business.id}`);
    } catch (e) {
      setError((e as Error).message);
      setSubmitting(false);
    }
  };

  const detailsValid = form.legal_name.trim() && form.founder_name.trim();

  return (
    <main className="container">
      <div className="brand">
        <div className="logo">🚀</div>
        <h1>Business Launch Orchestrator</h1>
      </div>
      <p className="subtitle">
        Incorporate, register for tax, open banking, activate payments and file
        compliance — across India, the Philippines and the US — from one flow.
      </p>

      <Stepper steps={WIZARD} current={phase} />

      {/* Phase 0 — choose jurisdiction */}
      {phase === 0 && (
        <div className="card">
          <h3 style={{ marginTop: 0 }}>Where are you launching?</h3>
          <div className="grid">
            {(["IN", "PH", "US"] as CountryCode[]).map((c) => {
              const info = countries.find((x) => x.code === c);
              return (
                <button
                  key={c}
                  className={`country-card ${country === c ? "selected" : ""}`}
                  onClick={() => selectCountry(c)}
                >
                  <div className="flag">{FLAGS[c]}</div>
                  <div className="name">{info?.name ?? c}</div>
                  <div className="meta">
                    {(info?.plan.length ?? 7)} steps ·{" "}
                    {info?.plan.some((s) => s.mode === "live")
                      ? "live payments"
                      : "sandbox"}
                  </div>
                </button>
              );
            })}
          </div>

          {country && (
            <>
              <label>Entity type</label>
              <select
                value={entityType}
                onChange={(e) => setEntityType(e.target.value)}
              >
                {ENTITY_TYPES[country].map((t) => (
                  <option key={t}>{t}</option>
                ))}
              </select>
            </>
          )}

          {error && <div className="error">{error}</div>}

          <div className="actions">
            <span />
            <button
              className="btn"
              disabled={!country}
              onClick={() => setPhase(1)}
            >
              Continue →
            </button>
          </div>
        </div>
      )}

      {/* Phase 1 — founder & company */}
      {phase === 1 && (
        <div className="card">
          <h3 style={{ marginTop: 0 }}>Founder &amp; company details</h3>
          <label>Legal / proposed company name *</label>
          <input value={form.legal_name} onChange={set("legal_name")} placeholder="Acme Technologies" />

          <div className="row">
            <div>
              <label>Founder full name *</label>
              <input value={form.founder_name} onChange={set("founder_name")} placeholder="Jane Doe" />
            </div>
            <div>
              <label>Founder email</label>
              <input value={form.founder_email} onChange={set("founder_email")} placeholder="jane@acme.com" />
            </div>
          </div>

          <div className="row">
            <div>
              <label>Founder phone</label>
              <input value={form.founder_phone} onChange={set("founder_phone")} placeholder="+1 555 010 0101" />
            </div>
            <div>
              <label>
                Founder tax ID{" "}
                <span className="muted">
                  ({country === "IN" ? "PAN" : country === "US" ? "SSN/ITIN" : "TIN"})
                </span>
              </label>
              <input value={form.founder_id_number} onChange={set("founder_id_number")} />
            </div>
          </div>

          <label>Registered address</label>
          <input value={form.line1} onChange={set("line1")} placeholder="Street address" />
          <div className="row" style={{ marginTop: 14 }}>
            <input value={form.city} onChange={set("city")} placeholder="City" />
            <input value={form.state} onChange={set("state")} placeholder="State / region" />
          </div>
          <div style={{ marginTop: 14 }}>
            <input value={form.postal_code} onChange={set("postal_code")} placeholder="Postal code" />
          </div>

          <div className="actions">
            <button className="btn secondary" onClick={() => setPhase(0)}>← Back</button>
            <button className="btn" disabled={!detailsValid} onClick={() => setPhase(2)}>
              Review plan →
            </button>
          </div>
        </div>
      )}

      {/* Phase 2 — review & launch */}
      {phase === 2 && country && (
        <div className="card">
          <h3 style={{ marginTop: 0 }}>
            {FLAGS[country]} {entityType} — {form.legal_name}
          </h3>
          <p className="muted">
            These {plan.length} API-backed steps will run in order. Steps marked{" "}
            <span className="badge live">live</span> hit a real provider sandbox;{" "}
            <span className="badge mock">mock</span> steps are deterministic stand-ins.
          </p>
          <PlanPreview plan={plan} />

          {error && <div className="error">{error}</div>}

          <div className="actions">
            <button className="btn secondary" onClick={() => setPhase(1)} disabled={submitting}>
              ← Back
            </button>
            <button className="btn" onClick={launch} disabled={submitting}>
              {submitting ? <span className="spin" /> : "🚀"} Create launch
            </button>
          </div>
        </div>
      )}
    </main>
  );
}
