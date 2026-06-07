"use client";

import { useState } from "react";

export function CodeBlock({ code }: { code: string }) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* clipboard unavailable */
    }
  };
  return (
    <div className="code">
      <button className="copy" onClick={copy}>
        {copied ? "✓ Copied" : "Copy"}
      </button>
      <pre>{code}</pre>
    </div>
  );
}
