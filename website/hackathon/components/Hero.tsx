import Countdown from './Countdown';

export default function Hero() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Background effects */}
      <div className="absolute inset-0 bg-grid-pattern bg-[size:40px_40px]" />
      <div className="absolute inset-0 bg-gradient-to-b from-transparent via-slate-900/50 to-slate-900" />

      {/* Floating orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-vibe-accent/20 rounded-full blur-[100px] animate-pulse-slow" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-vibe-primary/20 rounded-full blur-[100px] animate-pulse-slow" style={{ animationDelay: '1s' }} />

      <div className="relative z-10 text-center px-4 max-w-5xl mx-auto">
        {/* Logo */}
        <div className="mb-8 flex justify-center">
          <div className="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-slate-800/50 border border-slate-700">
            <svg width="32" height="32" viewBox="0 0 64 64" className="text-vibe-primary">
              <g transform="translate(32, 32)">
                <ellipse cx="0" cy="-16" rx="20" ry="7" fill="currentColor"/>
                <rect x="-20" y="-16" width="40" height="24" fill="currentColor"/>
                <ellipse cx="0" cy="8" rx="20" ry="7" fill="currentColor"/>
                <line x1="-20" y1="-4" x2="20" y2="-4" stroke="#2563EB" strokeWidth="2" opacity="0.5"/>
                <line x1="-20" y1="4" x2="20" y2="4" stroke="#2563EB" strokeWidth="2" opacity="0.5"/>
                <ellipse cx="0" cy="-16" rx="20" ry="7" fill="none" stroke="#2563EB" strokeWidth="2" opacity="0.5"/>
              </g>
            </svg>
            <span className="text-lg font-semibold text-slate-200">VibeSQL</span>
          </div>
        </div>

        {/* Main headline */}
        <h1 className="text-5xl md:text-7xl lg:text-8xl font-black mb-6">
          <span className="bg-gradient-to-r from-white via-vibe-light to-vibe-accent bg-clip-text text-transparent">
            HACKATHON
          </span>
          <br />
          <span className="text-3xl md:text-5xl lg:text-6xl text-slate-400 font-bold">
            2026
          </span>
        </h1>

        {/* Tagline */}
        <p className="text-xl md:text-2xl text-slate-300 mb-4 max-w-2xl mx-auto">
          54 hours to build amazing apps or hunt bugs
        </p>
        <p className="text-lg text-vibe-accent font-mono mb-8">
          March 20-22, 2026 | Online
        </p>

        {/* Countdown */}
        <div className="mb-10 p-6 rounded-2xl bg-slate-800/30 backdrop-blur-sm border border-slate-700/50">
          <Countdown />
        </div>

        {/* CTA Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <a
            href="/register"
            className="btn-primary text-lg px-10 py-4 inline-flex items-center gap-2"
          >
            Register Now
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </a>
          <a
            href="#tracks"
            className="btn-secondary text-lg px-10 py-4 inline-flex items-center gap-2"
          >
            Learn More
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </a>
        </div>

        {/* Quick links */}
        <div className="mt-8 flex flex-wrap justify-center gap-4 text-sm">
          <a
            href="https://github.com/PayEz-Net/vibesql-skills"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700 hover:border-vibe-accent transition-colors"
          >
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Starter Kit
          </a>
          <a
            href="https://vibesql.online/docs.html"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700 hover:border-vibe-accent transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            Documentation
          </a>
        </div>
      </div>

      {/* Scroll indicator */}
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 animate-bounce">
        <svg className="w-6 h-6 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
        </svg>
      </div>
    </section>
  );
}
