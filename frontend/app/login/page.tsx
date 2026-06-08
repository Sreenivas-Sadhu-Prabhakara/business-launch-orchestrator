"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";

export default function Login() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError("");
    try {
      await api.login(username, password);
      router.replace("/");
    } catch (err) {
      setError((err as Error).message);
      setBusy(false);
    }
  };

  return (
    <main className="login-wrap">
      <form className="login-card" onSubmit={submit}>
        <div className="login-brand">◆</div>
        <h1 className="login-title">Sign in</h1>
        <p className="login-sub">Business Launch Orchestrator</p>

        <label>Username</label>
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          autoFocus
          autoComplete="username"
          placeholder="your username"
        />

        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="current-password"
          placeholder="••••••••"
        />

        {error && <div className="error">{error}</div>}

        <button
          className="btn"
          type="submit"
          disabled={busy || !username || !password}
          style={{ width: "100%", justifyContent: "center", marginTop: 20 }}
        >
          {busy ? <span className="spin" /> : null} Sign in
        </button>

        <p className="login-hint">
          Demo access — <code>demo</code> / <code>demo123</code>
        </p>
      </form>
    </main>
  );
}
