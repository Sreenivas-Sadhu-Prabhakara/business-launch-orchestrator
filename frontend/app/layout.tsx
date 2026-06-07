import type { Metadata } from "next";
import "./globals.css";
import { Nav } from "@/components/Nav";

export const metadata: Metadata = {
  title: "Business Launch Orchestrator",
  description:
    "End-to-end flow to launch a business in India, the Philippines or the US — KYC, registration, tax, banking, payments and compliance, orchestrated.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <Nav />
        {children}
      </body>
    </html>
  );
}
