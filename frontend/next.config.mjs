/** @type {import('next').NextConfig} */

// Where the Next server proxies API calls. In Docker this is the backend
// service over the compose network; for `npm run dev` it's localhost.
const backend = process.env.BACKEND_ORIGIN || "http://localhost:8080";

const nextConfig = {
  reactStrictMode: true,
  // Emit a self-contained server bundle for slim serverless containers
  // (Cloud Run / Azure Container Apps / App Runner).
  output: "standalone",
  // Serve the API under the same origin as the UI, so a single public URL
  // (e.g. a Cloudflare tunnel) works end-to-end with no CORS and no baked
  // localhost dependency in the browser bundle.
  async rewrites() {
    return [
      { source: "/api/:path*", destination: `${backend}/api/:path*` },
      { source: "/healthz", destination: `${backend}/healthz` },
    ];
  },
};

export default nextConfig;
