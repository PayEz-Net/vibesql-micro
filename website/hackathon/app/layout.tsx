import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "VibeSQL Hackathon 2026 | Build Apps & Find Bugs",
  description: "Join the 54-hour VibeSQL hackathon! Build innovative apps with PostgreSQL for AI agents or hunt bugs for bounties. March 20-22, 2026.",
  keywords: ["hackathon", "vibesql", "postgresql", "database", "ai", "bug bounty", "coding competition"],
  authors: [{ name: "VibeSQL Team" }],
  openGraph: {
    title: "VibeSQL Hackathon 2026",
    description: "54-hour hackathon: Build apps with VibeSQL or hunt bugs for bounties!",
    type: "website",
    url: "https://vibesql.online/hackathon",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
        alt: "VibeSQL Hackathon 2026",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "VibeSQL Hackathon 2026",
    description: "54-hour hackathon: Build apps with VibeSQL or hunt bugs for bounties!",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="/favicon.ico" />
        <link rel="icon" type="image/svg+xml" href="/vibesql-icon.svg" />
      </head>
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
