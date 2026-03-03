export default function Footer() {
  return (
    <footer className="border-t border-slate-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-3 mb-4">
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
              <span className="text-xl font-bold">VibeSQL Hackathon</span>
            </div>
            <p className="text-slate-400 text-sm max-w-md">
              A 54-hour hackathon to build amazing applications with VibeSQL and hunt bugs for bounties.
              Join developers worldwide for an exciting weekend of coding.
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="font-semibold mb-4">Quick Links</h4>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="#about" className="text-slate-400 hover:text-vibe-accent transition-colors">About</a>
              </li>
              <li>
                <a href="#tracks" className="text-slate-400 hover:text-vibe-accent transition-colors">Tracks</a>
              </li>
              <li>
                <a href="#prizes" className="text-slate-400 hover:text-vibe-accent transition-colors">Prizes</a>
              </li>
              <li>
                <a href="#schedule" className="text-slate-400 hover:text-vibe-accent transition-colors">Schedule</a>
              </li>
              <li>
                <a href="#rules" className="text-slate-400 hover:text-vibe-accent transition-colors">Rules</a>
              </li>
              <li>
                <a href="#faq" className="text-slate-400 hover:text-vibe-accent transition-colors">FAQ</a>
              </li>
            </ul>
          </div>

          {/* Resources */}
          <div>
            <h4 className="font-semibold mb-4">Resources</h4>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://vibesql.online/docs.html"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-vibe-accent transition-colors inline-flex items-center gap-1"
                >
                  Documentation
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/PayEz-Net/vibesql-skills"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-vibe-accent transition-colors inline-flex items-center gap-1"
                >
                  Starter Kit
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/PayEz-Net/vibesql-micro"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-vibe-accent transition-colors inline-flex items-center gap-1"
                >
                  GitHub Repo
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/PayEz-Net/vibesql-micro/discussions"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-vibe-accent transition-colors inline-flex items-center gap-1"
                >
                  Discussions
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="mt-12 pt-8 border-t border-slate-800 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-slate-500 text-sm">
            &copy; {new Date().getFullYear()} VibeSQL. All rights reserved.
          </p>
          <div className="flex items-center gap-6">
            <a
              href="https://github.com/PayEz-Net"
              target="_blank"
              rel="noopener noreferrer"
              className="text-slate-400 hover:text-white transition-colors"
            >
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
            </a>
            <a
              href="https://vibesql.online"
              target="_blank"
              rel="noopener noreferrer"
              className="text-slate-400 hover:text-white transition-colors"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
