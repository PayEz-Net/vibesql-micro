export default function Rules() {
  return (
    <section id="rules" className="section-container bg-slate-900/50">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Rules & <span className="text-vibe-accent">Judging</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          Everything you need to know to participate and win.
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-12">
        {/* Eligibility */}
        <div>
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-5 h-5 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-bold">Eligibility</h3>
          </div>
          <ul className="space-y-3 text-slate-400">
            <li className="flex items-start gap-3">
              <span className="text-vibe-accent mt-1">•</span>
              Open to developers worldwide, 18 years or older
            </li>
            <li className="flex items-start gap-3">
              <span className="text-vibe-accent mt-1">•</span>
              Individual or team participation (max 4 members per team)
            </li>
            <li className="flex items-start gap-3">
              <span className="text-vibe-accent mt-1">•</span>
              Must register before hackathon starts
            </li>
            <li className="flex items-start gap-3">
              <span className="text-vibe-accent mt-1">•</span>
              Projects must use VibeSQL as the primary database
            </li>
            <li className="flex items-start gap-3">
              <span className="text-vibe-accent mt-1">•</span>
              Code must be written during the hackathon period
            </li>
          </ul>
        </div>

        {/* Submission Requirements */}
        <div>
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-5 h-5 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 className="text-xl font-bold">Submission Requirements</h3>
          </div>

          <div className="space-y-6">
            <div>
              <h4 className="font-semibold text-white mb-2">App Track:</h4>
              <ul className="space-y-2 text-slate-400 text-sm">
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  Public GitHub repository with source code
                </li>
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  Demo video (max 5 minutes, YouTube or Loom)
                </li>
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  Project writeup explaining what you built
                </li>
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-slate-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Live demo URL (optional but recommended)
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-white mb-2">Bug Track:</h4>
              <ul className="space-y-2 text-slate-400 text-sm">
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  GitHub Issue with clear reproduction steps
                </li>
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  Expected vs actual behavior
                </li>
                <li className="flex items-start gap-2">
                  <svg className="w-4 h-4 text-slate-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Pull Request with fix (earns extra points)
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Judging Criteria */}
      <div className="mt-16">
        <div className="text-center mb-8">
          <h3 className="text-2xl font-bold">Judging Criteria</h3>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {/* App Track Criteria */}
          <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
            <h4 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <svg className="w-5 h-5 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
              App Track
            </h4>
            <div className="space-y-4">
              {[
                { label: 'Innovation', percent: 25, desc: 'Creative use of VibeSQL' },
                { label: 'Technical Execution', percent: 25, desc: 'Code quality & architecture' },
                { label: 'Usefulness', percent: 25, desc: 'Solves a real problem' },
                { label: 'Presentation', percent: 25, desc: 'Demo & documentation quality' },
              ].map((item, i) => (
                <div key={i}>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-medium">{item.label}</span>
                    <span className="text-sm text-vibe-accent">{item.percent}%</span>
                  </div>
                  <div className="h-2 rounded-full bg-slate-700 overflow-hidden">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-vibe-primary to-vibe-accent"
                      style={{ width: `${item.percent}%` }}
                    />
                  </div>
                  <p className="text-xs text-slate-500 mt-1">{item.desc}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Bug Track Criteria */}
          <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
            <h4 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <svg className="w-5 h-5 text-hackathon-bug-high" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              Bug Track
            </h4>
            <div className="space-y-3">
              <div className="p-3 rounded-lg bg-hackathon-bug-critical/10 border border-hackathon-bug-critical/20">
                <div className="font-semibold text-hackathon-bug-critical">Severity</div>
                <p className="text-sm text-slate-400">Critical &gt; High &gt; Medium &gt; Low</p>
              </div>
              <div className="p-3 rounded-lg bg-slate-700/50 border border-slate-600">
                <div className="font-semibold">Impact</div>
                <p className="text-sm text-slate-400">How many users are affected by the bug</p>
              </div>
              <div className="p-3 rounded-lg bg-slate-700/50 border border-slate-600">
                <div className="font-semibold">Fix Quality</div>
                <p className="text-sm text-slate-400">For PRs: quality and completeness of the fix</p>
              </div>
              <div className="p-3 rounded-lg bg-slate-700/50 border border-slate-600">
                <div className="font-semibold">Report Quality</div>
                <p className="text-sm text-slate-400">Clear reproduction steps, minimal examples</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
