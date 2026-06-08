"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { api, type AuthUser } from "@/lib/api";
import { Nav } from "@/components/Nav";

/**
 * Gates the whole app behind a session. Unauthenticated visitors are redirected
 * to /login; the login route itself renders without the gate.
 */
export function AuthGate({ children }: { children: React.ReactNode }) {
  const path = usePathname();
  const router = useRouter();
  const [state, setState] = useState<"loading" | "authed" | "anon">("loading");
  const [user, setUser] = useState<AuthUser | null>(null);

  useEffect(() => {
    if (path === "/login") {
      setState("anon");
      return;
    }
    let active = true;
    api
      .me()
      .then((u) => {
        if (active) {
          setUser(u);
          setState("authed");
        }
      })
      .catch(() => {
        if (active) {
          setState("anon");
          router.replace("/login");
        }
      });
    return () => {
      active = false;
    };
  }, [path, router]);

  const logout = async () => {
    try {
      await api.logout();
    } catch {
      /* ignore */
    }
    setUser(null);
    router.replace("/login");
  };

  if (path === "/login") return <>{children}</>;

  if (state !== "authed") {
    return <div className="auth-splash">◆</div>;
  }

  return (
    <>
      <Nav user={user} onLogout={logout} />
      {children}
    </>
  );
}
