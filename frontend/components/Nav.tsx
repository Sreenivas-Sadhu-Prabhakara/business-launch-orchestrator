"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const LINKS = [
  { href: "/", label: "Launch" },
  { href: "/how-it-works", label: "How it works" },
  { href: "/deploy", label: "Deploy" },
];

export function Nav() {
  const path = usePathname();
  return (
    <nav className="nav">
      <div className="nav-inner">
        <Link href="/" className="nav-brand">
          <span className="logo" style={{ display: "grid", placeItems: "center", background: "linear-gradient(135deg,#6d8bff,#8b5cf6)" }}>🚀</span>
          Launch Orchestrator
        </Link>
        <div className="nav-links">
          {LINKS.map((l) => (
            <Link
              key={l.href}
              href={l.href}
              className={path === l.href ? "active" : ""}
            >
              {l.label}
            </Link>
          ))}
        </div>
      </div>
    </nav>
  );
}
