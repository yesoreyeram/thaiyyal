import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  devIndicators: false,
  output: 'export',
  images: {
    unoptimized: true,
  },
  // For GitHub Pages with dynamic routes, we'll use a SPA fallback approach
  // The 404.html will be used as a fallback for all routes
  trailingSlash: true,
};

export default nextConfig;
