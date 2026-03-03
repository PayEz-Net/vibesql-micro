export default function Tracks() {
  return (
    <section id="tracks" className="section-container">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Two Ways to <span className="text-vibe-accent">Win</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          Choose your path: build innovative applications or hunt bugs for bounties. You can participate in both!
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-8">
        {/* App Track */}
        <div className="prize-card group">
          <div className="flex items-center gap-4 mb-6">
            <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-vibe-accent to-vibe-primary flex items-center justify-center">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
            </div>
            <div>
              <span className="text-xs uppercase tracking-wider text-vibe-accent font-semibold">Track 1</span>
              <h3 className="text-2xl font-bold">Build Apps</h3>
            </div>
          </div>

          <p className="text-slate-300 mb-6">
            Create innovative applications powered by VibeSQL. Show us what&apos;s possible when databases are this easy.
          </p>

          <div className="space-y-4 mb-6">
            <h4 className="font-semibold text-slate-200">What to Build:</h4>
            <ul className="space-y-2 text-slate-400">
              <li className="flex items-start gap-2">
                <svg className="w-5 h-5 text-vibe-accent flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                AI agents with persistent memory
              </li>
              <li className="flex items-start gap-2">
                <svg className="w-5 h-5 text-vibe-accent flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Local-first apps with sync
              </li>
              <li className="flex items-start gap-2">
                <svg className="w-5 h-5 text-vibe-accent flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Developer tools & CLI apps
              </li>
              <li className="flex items-start gap-2">
                <svg className="w-5 h-5 text-vibe-accent flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Data visualization dashboards
              </li>
              <li className="flex items-start gap-2">
                <svg className="w-5 h-5 text-vibe-accent flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Edge computing solutions
              </li>
            </ul>
          </div>

          <div className="pt-4 border-t border-slate-700/50">
            <span className="text-sm text-slate-500">Submission Requirements</span>
            <p className="text-slate-400 text-sm mt-1">
              GitHub repo + Demo video + Project writeup
            </p>
          </div>
        </div>

        {/* Bug Track */}
        <div className="prize-card group">
          <div className="flex items-center gap-4 mb-6">
            <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-hackathon-bug-critical to-orange-600 flex items-center justify-center">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <div>
              <span className="text-xs uppercase tracking-wider text-hackathon-bug-high font-semibold">Track 2</span>
              <h3 className="text-2xl font-bold">Find Bugs</h3>
            </div>
          </div>

          <p className="text-slate-300 mb-6">
            Hunt for bugs in VibeSQL and earn bounties. Critical bugs = bigger rewards. Top hunters win bonus prizes!
          </p>

          <div className="space-y-4 mb-6">
            <h4 className="font-semibold text-slate-200">Bug Categories:</h4>
            <ul className="space-y-2">
              <li className="flex items-center gap-3">
                <span className="px-2 py-1 text-xs font-bold rounded bg-hackathon-bug-critical/20 text-hackathon-bug-critical border border-hackathon-bug-critical/30">
                  CRITICAL
                </span>
                <span className="text-slate-400">Security vulnerabilities, data loss</span>
              </li>
              <li className="flex items-center gap-3">
                <span className="px-2 py-1 text-xs font-bold rounded bg-hackathon-bug-high/20 text-hackathon-bug-high border border-hackathon-bug-high/30">
                  HIGH
                </span>
                <span className="text-slate-400">Major functionality issues</span>
              </li>
              <li className="flex items-center gap-3">
                <span className="px-2 py-1 text-xs font-bold rounded bg-hackathon-bug-medium/20 text-hackathon-bug-medium border border-hackathon-bug-medium/30">
                  MEDIUM
                </span>
                <span className="text-slate-400">Performance, edge cases</span>
              </li>
              <li className="flex items-center gap-3">
                <span className="px-2 py-1 text-xs font-bold rounded bg-hackathon-bug-low/20 text-hackathon-bug-low border border-hackathon-bug-low/30">
                  LOW
                </span>
                <span className="text-slate-400">Minor issues, documentation</span>
              </li>
            </ul>
          </div>

          <div className="pt-4 border-t border-slate-700/50">
            <span className="text-sm text-slate-500">Submission Requirements</span>
            <p className="text-slate-400 text-sm mt-1">
              GitHub Issue or Pull Request with reproduction steps
            </p>
          </div>
        </div>
      </div>

      {/* Resources */}
      <div className="mt-12 p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
        <h3 className="text-xl font-semibold mb-4 text-center">Getting Started Resources</h3>
        <div className="flex flex-wrap justify-center gap-4">
          <a
            href="https://github.com/PayEz-Net/vibesql-skills"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-slate-700/50 hover:bg-slate-700 transition-colors"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Starter Kit
          </a>
          <a
            href="https://vibesql.online/docs.html"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-slate-700/50 hover:bg-slate-700 transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            Documentation
          </a>
          <a
            href="https://github.com/PayEz-Net/vibesql-micro"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-slate-700/50 hover:bg-slate-700 transition-colors"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Source Code
          </a>
        </div>
      </div>
    </section>
  );
}
