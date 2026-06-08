"use client";

import { useEffect, useState } from "react";
import { api, type AdminUser, type AuthUser } from "@/lib/api";

export default function AdminUsers() {
  const [me, setMe] = useState<AuthUser | null>(null);
  const [ready, setReady] = useState(false);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [listErr, setListErr] = useState("");

  const [form, setForm] = useState<{ username: string; password: string; role: "user" | "admin" }>({
    username: "",
    password: "",
    role: "user",
  });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [ok, setOk] = useState("");

  const loadUsers = () =>
    api.listUsers().then(setUsers).catch((e) => setListErr((e as Error).message));

  useEffect(() => {
    api
      .me()
      .then((u) => {
        setMe(u);
        setReady(true);
        if (u.role === "admin") loadUsers();
      })
      .catch(() => setReady(true));
  }, []);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError("");
    setOk("");
    try {
      await api.createUser(form.username.trim(), form.password, form.role);
      setOk(`Created “${form.username.trim()}”.`);
      setForm({ username: "", password: "", role: "user" });
      await loadUsers();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  };

  const [rowBusy, setRowBusy] = useState("");
  const [rowErr, setRowErr] = useState("");

  const changeRole = async (u: AdminUser) => {
    const next = u.role === "admin" ? "user" : "admin";
    setRowBusy(u.id);
    setRowErr("");
    try {
      await api.updateUserRole(u.id, next);
      await loadUsers();
    } catch (e) {
      setRowErr((e as Error).message);
    } finally {
      setRowBusy("");
    }
  };

  const removeUser = async (u: AdminUser) => {
    if (!window.confirm(`Delete “${u.username}”? This cannot be undone.`)) return;
    setRowBusy(u.id);
    setRowErr("");
    try {
      await api.deleteUser(u.id);
      await loadUsers();
    } catch (e) {
      setRowErr((e as Error).message);
    } finally {
      setRowBusy("");
    }
  };

  if (!ready) {
    return (
      <main className="container">
        <p className="muted">Loading…</p>
      </main>
    );
  }

  if (me?.role !== "admin") {
    return (
      <main className="container">
        <div className="card">
          <h3 style={{ marginTop: 0 }}>Admins only</h3>
          <p className="muted">You need an administrator account to manage users.</p>
        </div>
      </main>
    );
  }

  return (
    <main className="container">
      <div className="brand">
        <div className="logo">◆</div>
        <h1>Users</h1>
      </div>
      <p className="subtitle">
        Create and review accounts. Admins see every launch; users see only their own.
      </p>

      <form className="card" onSubmit={submit}>
        <h3 style={{ marginTop: 0 }}>Create an account</h3>
        <div className="row">
          <div>
            <label>Username</label>
            <input
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              placeholder="jane"
              autoComplete="off"
            />
          </div>
          <div>
            <label>Password</label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              placeholder="••••••••"
              autoComplete="new-password"
            />
          </div>
        </div>
        <label>Role</label>
        <select
          value={form.role}
          onChange={(e) => setForm({ ...form, role: e.target.value as "user" | "admin" })}
        >
          <option value="user">user — sees only their own launches</option>
          <option value="admin">admin — sees all launches &amp; manages users</option>
        </select>

        {error && <div className="error">{error}</div>}
        {ok && (
          <div className="muted" style={{ marginTop: 12, color: "var(--ok)" }}>
            {ok}
          </div>
        )}

        <div className="actions">
          <span />
          <button className="btn" type="submit" disabled={busy || !form.username.trim() || !form.password}>
            {busy ? <span className="spin" /> : null} Create user
          </button>
        </div>
      </form>

      <div className="section-title">All accounts ({users.length})</div>
      {rowErr && <div className="error">{rowErr}</div>}
      {listErr ? (
        <div className="error">{listErr}</div>
      ) : (
        <table className="tbl">
          <thead>
            <tr>
              <th>Username</th>
              <th>Role</th>
              <th>Created</th>
              <th style={{ textAlign: "right" }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => {
              const self = u.id === me.id;
              const working = rowBusy === u.id;
              return (
                <tr key={u.id}>
                  <td>
                    {u.username}
                    {self ? <span className="muted"> · you</span> : null}
                  </td>
                  <td>
                    <span className="nav-role">{u.role}</span>
                  </td>
                  <td className="muted">{new Date(u.created_at).toLocaleString()}</td>
                  <td>
                    {self ? (
                      <span className="muted" style={{ display: "block", textAlign: "right" }}>—</span>
                    ) : (
                      <div className="row-actions">
                        <button
                          className="act"
                          disabled={working}
                          onClick={() => changeRole(u)}
                          title={`Change role to ${u.role === "admin" ? "user" : "admin"}`}
                        >
                          {working ? <span className="spin" /> : null}
                          {u.role === "admin" ? "Make user" : "Make admin"}
                        </button>
                        <button
                          className="act danger"
                          disabled={working}
                          onClick={() => removeUser(u)}
                        >
                          Delete
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </main>
  );
}
