import type { NextConfig } from "next";

const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

const nextConfig: NextConfig = {
  rewrites: async () => [
    { source: "/api/users/:path*", destination: `${apiBase}/v1/users/:path*` },
    { source: "/api/catalog/:path*", destination: `${apiBase}/v1/catalog/:path*` },
  ],
};

export default nextConfig;
