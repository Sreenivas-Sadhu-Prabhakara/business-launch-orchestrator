// Typed client for the Go orchestrator API.

export const API_BASE =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type CountryCode = "IN" | "PH" | "US";

export interface PlannedStep {
  seq: number;
  type: string;
  provider: string;
  title: string;
  mode: "live" | "mock";
}

export interface CountryInfo {
  code: CountryCode;
  name: string;
  plan: PlannedStep[];
}

export interface Address {
  line1?: string;
  line2?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  country?: string;
}

export interface Business {
  id: string;
  country: CountryCode;
  entity_type: string;
  legal_name: string;
  founder_name: string;
  founder_email: string;
  founder_phone: string;
  founder_id_number: string;
  address: Address;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface LaunchStep {
  id: string;
  business_id: string;
  seq: number;
  step_type: string;
  provider: string;
  title: string;
  mode: "live" | "mock";
  status: "pending" | "running" | "completed" | "failed";
  response: Record<string, unknown>;
  external_ref: string;
  error: string;
  completed_at?: string | null;
}

export interface BusinessDetail {
  business: Business;
  steps: LaunchStep[];
}

export interface CreateBusinessInput {
  country: CountryCode;
  entity_type: string;
  legal_name: string;
  founder_name: string;
  founder_email: string;
  founder_phone: string;
  founder_id_number: string;
  address: Address;
}

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
    cache: "no-store",
  });
  if (!res.ok) {
    let msg = `Request failed (${res.status})`;
    try {
      const body = await res.json();
      if (body?.error) msg = body.error;
    } catch {
      /* ignore */
    }
    throw new Error(msg);
  }
  return res.json() as Promise<T>;
}

export const api = {
  listCountries: () =>
    req<{ countries: CountryInfo[] }>("/api/v1/countries").then((d) => d.countries),

  getPlan: (code: CountryCode) =>
    req<CountryInfo>(`/api/v1/countries/${code}/plan`),

  createBusiness: (input: CreateBusinessInput) =>
    req<BusinessDetail>("/api/v1/businesses", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getBusiness: (id: string) =>
    req<BusinessDetail>(`/api/v1/businesses/${id}`),

  advance: (id: string) =>
    req<{ step?: LaunchStep; done?: boolean; message?: string }>(
      `/api/v1/businesses/${id}/advance`,
      { method: "POST" }
    ),

  runAll: (id: string) =>
    req<BusinessDetail>(`/api/v1/businesses/${id}/run`, { method: "POST" }),
};
