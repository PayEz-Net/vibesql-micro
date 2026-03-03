export default function Prizes() {
  return (
    <section id="prizes" className="section-container bg-slate-900/50">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Prizes & <span className="text-vibe-accent">Rewards</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          Compete for prizes in both tracks. Build apps, find bugs, win rewards!
        </p>
      </div>

      <div className="grid lg:grid-cols-2 gap-12">
        {/* App Track Prizes */}
        <div>
          <div className="flex items-center gap-3 mb-8">
            <svg className="w-8 h-8 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
            </svg>
            <h3 className="text-2xl font-bold">App Track Prizes</h3>
          </div>

          <div className="space-y-4">
            {/* 1st Place */}
            <div className="prize-card flex items-center gap-6 border-l-4 border-hackathon-gold">
              <div className="w-16 h-16 rounded-full bg-gradient-to-br from-hackathon-gold to-yellow-600 flex items-center justify-center text-2xl font-bold text-slate-900">
                1st
              </div>
              <div className="flex-1">
                <div className="text-2xl font-bold text-hackathon-gold">$TBD</div>
                <p className="text-slate-400">Grand Prize Winner</p>
              </div>
              <svg className="w-10 h-10 text-hackathon-gold" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
            </div>

            {/* 2nd Place */}
            <div className="prize-card flex items-center gap-6 border-l-4 border-hackathon-silver">
              <div className="w-16 h-16 rounded-full bg-gradient-to-br from-hackathon-silver to-gray-500 flex items-center justify-center text-2xl font-bold text-slate-900">
                2nd
              </div>
              <div className="flex-1">
                <div className="text-2xl font-bold text-hackathon-silver">$TBD</div>
                <p className="text-slate-400">Runner Up</p>
              </div>
            </div>

            {/* 3rd Place */}
            <div className="prize-card flex items-center gap-6 border-l-4 border-hackathon-bronze">
              <div className="w-16 h-16 rounded-full bg-gradient-to-br from-hackathon-bronze to-orange-800 flex items-center justify-center text-2xl font-bold text-slate-900">
                3rd
              </div>
              <div className="flex-1">
                <div className="text-2xl font-bold text-hackathon-bronze">$TBD</div>
                <p className="text-slate-400">Third Place</p>
              </div>
            </div>
          </div>
        </div>

        {/* Bug Bounty Prizes */}
        <div>
          <div className="flex items-center gap-3 mb-8">
            <svg className="w-8 h-8 text-hackathon-bug-high" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 className="text-2xl font-bold">Bug Bounty Rewards</h3>
          </div>

          {/* Per-Bug Rewards */}
          <div className="mb-8">
            <h4 className="text-lg font-semibold mb-4 text-slate-300">Per-Bug Rewards</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="prize-card text-center border-t-4 border-hackathon-bug-critical">
                <div className="text-xs uppercase tracking-wider text-hackathon-bug-critical mb-2">Critical</div>
                <div className="text-2xl font-bold">$TBD</div>
              </div>
              <div className="prize-card text-center border-t-4 border-hackathon-bug-high">
                <div className="text-xs uppercase tracking-wider text-hackathon-bug-high mb-2">High</div>
                <div className="text-2xl font-bold">$TBD</div>
              </div>
              <div className="prize-card text-center border-t-4 border-hackathon-bug-medium">
                <div className="text-xs uppercase tracking-wider text-hackathon-bug-medium mb-2">Medium</div>
                <div className="text-2xl font-bold">$TBD</div>
              </div>
              <div className="prize-card text-center border-t-4 border-hackathon-bug-low">
                <div className="text-xs uppercase tracking-wider text-hackathon-bug-low mb-2">Low</div>
                <div className="text-2xl font-bold">$TBD</div>
              </div>
            </div>
          </div>

          {/* Leaderboard Bonuses */}
          <div>
            <h4 className="text-lg font-semibold mb-4 text-slate-300">Leaderboard Bonuses</h4>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-4 rounded-lg bg-slate-800/50 border border-slate-700">
                <div className="flex items-center gap-3">
                  <span className="text-hackathon-gold text-lg">
                    <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                    </svg>
                  </span>
                  <span className="font-semibold">Top Bug Hunter</span>
                </div>
                <span className="text-xl font-bold text-hackathon-gold">+$TBD</span>
              </div>
              <div className="flex items-center justify-between p-4 rounded-lg bg-slate-800/50 border border-slate-700">
                <div className="flex items-center gap-3">
                  <span className="text-hackathon-silver text-lg">2nd</span>
                  <span className="font-semibold">Second Place</span>
                </div>
                <span className="text-xl font-bold text-hackathon-silver">+$TBD</span>
              </div>
              <div className="flex items-center justify-between p-4 rounded-lg bg-slate-800/50 border border-slate-700">
                <div className="flex items-center gap-3">
                  <span className="text-hackathon-bronze text-lg">3rd</span>
                  <span className="font-semibold">Third Place</span>
                </div>
                <span className="text-xl font-bold text-hackathon-bronze">+$TBD</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Note */}
      <div className="mt-12 text-center">
        <p className="text-slate-500 text-sm">
          Prize amounts will be announced soon. Stay tuned!
        </p>
      </div>
    </section>
  );
}
