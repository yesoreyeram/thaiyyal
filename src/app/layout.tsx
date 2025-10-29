import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Thaiyyal - Workflow Builder",
  description: "Visual workflow builder with comprehensive node types",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
