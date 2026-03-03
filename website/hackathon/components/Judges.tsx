export default function Judges() {
  const judges = [
    {
      name: 'TBA',
      role: 'Judge',
      bio: 'Judge details coming soon',
      avatar: null,
    },
    {
      name: 'TBA',
      role: 'Judge',
      bio: 'Judge details coming soon',
      avatar: null,
    },
    {
      name: 'TBA',
      role: 'Judge',
      bio: 'Judge details coming soon',
      avatar: null,
    },
  ];

  return (
    <section id="judges" className="section-container">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Meet the <span className="text-vibe-accent">Judges</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          Industry experts who will evaluate your submissions.
        </p>
      </div>

      <div className="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
        {judges.map((judge, index) => (
          <div key={index} className="prize-card text-center group">
            {/* Avatar */}
            <div className="w-24 h-24 mx-auto mb-4 rounded-full bg-gradient-to-br from-vibe-primary to-vibe-accent p-1">
              <div className="w-full h-full rounded-full bg-slate-800 flex items-center justify-center">
                {judge.avatar ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={judge.avatar}
                    alt={judge.name}
                    className="w-full h-full rounded-full object-cover"
                  />
                ) : (
                  <svg className="w-12 h-12 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                )}
              </div>
            </div>

            {/* Info */}
            <h3 className="text-xl font-bold mb-1">{judge.name}</h3>
            <p className="text-vibe-accent text-sm mb-3">{judge.role}</p>
            <p className="text-slate-400 text-sm">{judge.bio}</p>
          </div>
        ))}
      </div>

      {/* Coming soon note */}
      <div className="mt-12 text-center">
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-vibe-accent/10 border border-vibe-accent/30">
          <svg className="w-5 h-5 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span className="text-vibe-accent text-sm">
            Judge announcements coming soon!
          </span>
        </div>
      </div>
    </section>
  );
}
