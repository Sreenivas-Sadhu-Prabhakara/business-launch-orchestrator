"use client";

export function Stepper({
  steps,
  current,
}: {
  steps: string[];
  current: number;
}) {
  return (
    <div className="stepper">
      {steps.map((label, i) => (
        <span
          key={label}
          className={`pill ${i === current ? "active" : ""} ${
            i < current ? "done" : ""
          }`}
        >
          {i < current ? "✓ " : `${i + 1}. `}
          {label}
        </span>
      ))}
    </div>
  );
}
