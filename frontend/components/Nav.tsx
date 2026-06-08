"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { AuthUser } from "@/lib/api";

const LINKS = [
  { href: "/", label: "Launch" },
  { href: "/how-it-works", label: "How it works" },
  { href: "/deploy", label: "Deploy" },
];

export function Nav({
  user,
  onLogout,
}: {
  user?: AuthUser | null;
  onLogout?: () => void;
}) {
  const path = usePathname();
  return (
    <nav className="nav">
      <div className="nav-inner">
        <Link href="/" className="nav-brand">
          <span className="logo">◆</span>
          Launch Orchestrator
        </Link>
        <div className="nav-links">
          {LINKS.map((l) => (
            <Link key={l.href} href={l.href} className={path === l.href ? "active" : ""}>
              {l.label}
            </Link>
          ))}
          {user?.role === "admin" && (
            <Link
              href="/admin/users"
              className={path === "/admin/users" ? "active" : ""}
            >
              Users
            </Link>
          )}
          {user && (
            <span className="nav-user">
              {user.username}
              <span className="nav-role">{user.role}</span>
            </span>
          )}
          {onLogout && (
            <button className="nav-logout" onClick={onLogout}>
              Log out
            </button>
          )}
        </div>
      </div>
    </nav>
  );
}
