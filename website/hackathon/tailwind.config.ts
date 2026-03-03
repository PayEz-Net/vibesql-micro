import type { Config } from "tailwindcss";

export default {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // VibeSQL brand colors
        'vibe': {
          'primary': '#336791',      // PostgreSQL blue
          'accent': '#2563EB',       // Bright blue accent
          'dark': '#1e3a5f',         // Darker blue
          'light': '#4a90c2',        // Lighter blue
        },
        'hackathon': {
          'gold': '#FFD700',         // 1st place
          'silver': '#C0C0C0',       // 2nd place
          'bronze': '#CD7F32',       // 3rd place
          'bug-critical': '#DC2626', // Critical bugs
          'bug-high': '#EA580C',     // High severity
          'bug-medium': '#CA8A04',   // Medium severity
          'bug-low': '#16A34A',      // Low severity
        }
      },
      fontFamily: {
        'mono': ['JetBrains Mono', 'Fira Code', 'Consolas', 'monospace'],
        'sans': ['Inter', '-apple-system', 'BlinkMacSystemFont', 'sans-serif'],
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow': 'glow 2s ease-in-out infinite alternate',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px #2563EB, 0 0 10px #2563EB, 0 0 15px #2563EB' },
          '100%': { boxShadow: '0 0 10px #336791, 0 0 20px #336791, 0 0 30px #336791' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        }
      },
      backgroundImage: {
        'grid-pattern': 'linear-gradient(to right, rgba(51, 103, 145, 0.1) 1px, transparent 1px), linear-gradient(to bottom, rgba(51, 103, 145, 0.1) 1px, transparent 1px)',
      },
    },
  },
  plugins: [],
} satisfies Config;
