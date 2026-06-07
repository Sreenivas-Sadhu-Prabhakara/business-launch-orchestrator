/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Emit a self-contained server bundle for slim serverless containers
  // (Cloud Run / Azure Container Apps / App Runner).
  output: "standalone",
};

export default nextConfig;
