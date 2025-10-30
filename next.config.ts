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
  // Security headers
  // Note: These headers will not work with static export (output: 'export')
  // For production deployments with a server, remove 'output: export' to enable headers
  // For static deployments (like GitHub Pages), security headers should be configured
  // at the web server/CDN level (e.g., GitHub Pages doesn't support custom headers,
  // but services like Netlify, Vercel, or Cloudflare Pages do)
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on'
          },
          {
            key: 'Strict-Transport-Security',
            value: 'max-age=63072000; includeSubDomains; preload'
          },
          {
            key: 'X-Frame-Options',
            value: 'SAMEORIGIN'
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff'
          },
          {
            key: 'X-XSS-Protection',
            value: '1; mode=block'
          },
          {
            key: 'Referrer-Policy',
            value: 'strict-origin-when-cross-origin'
          },
          {
            key: 'Permissions-Policy',
            value: 'camera=(), microphone=(), geolocation=()'
          },
          {
            key: 'Content-Security-Policy',
            value: "default-src 'self'; script-src 'self' 'unsafe-eval' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'self';"
          }
        ],
      },
    ]
  },
};

export default nextConfig;
