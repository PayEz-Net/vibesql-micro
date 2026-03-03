export default function Schedule() {
  const events = [
    {
      date: 'Mar 20',
      time: '6:00 PM PST',
      title: 'Hackathon Kickoff',
      description: 'Opening ceremony, rules overview, and Q&A session',
      highlight: true,
    },
    {
      date: 'Mar 20',
      time: '6:30 PM PST',
      title: 'Hacking Begins',
      description: 'Start building! Access to Discord support channel opens',
      highlight: false,
    },
    {
      date: 'Mar 21',
      time: '12:00 PM PST',
      title: 'Mid-Hackathon Check-in',
      description: 'Optional: Share progress, get feedback, ask questions',
      highlight: false,
    },
    {
      date: 'Mar 22',
      time: '6:00 PM PST',
      title: 'Final Hours',
      description: '6 hours remaining! Reminder to prepare submissions',
      highlight: false,
    },
    {
      date: 'Mar 22',
      time: '11:59 PM PST',
      title: 'Submissions Close',
      description: 'All app and bug submissions must be in',
      highlight: true,
    },
    {
      date: 'Mar 23',
      time: 'All Day',
      title: 'Judging Period',
      description: 'Judges review all submissions',
      highlight: false,
    },
    {
      date: 'Mar 24',
      time: 'TBA',
      title: 'Winners Announced',
      description: 'Results revealed and prizes distributed',
      highlight: true,
    },
  ];

  return (
    <section id="schedule" className="section-container">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Event <span className="text-vibe-accent">Schedule</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          54 hours of hacking, learning, and building. All times in Pacific Standard Time (PST).
        </p>
      </div>

      <div className="max-w-3xl mx-auto">
        {events.map((event, index) => (
          <div key={index} className="timeline-item">
            <div className={`p-6 rounded-xl transition-all duration-300 ${
              event.highlight
                ? 'bg-vibe-accent/10 border border-vibe-accent/30'
                : 'bg-slate-800/30 border border-slate-700/50 hover:border-slate-600'
            }`}>
              <div className="flex flex-wrap items-center gap-4 mb-2">
                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                  event.highlight
                    ? 'bg-vibe-accent/20 text-vibe-accent'
                    : 'bg-slate-700 text-slate-300'
                }`}>
                  {event.date}
                </span>
                <span className="text-sm text-slate-400 font-mono">
                  {event.time}
                </span>
              </div>
              <h3 className={`text-xl font-bold mb-2 ${
                event.highlight ? 'text-vibe-accent' : 'text-white'
              }`}>
                {event.title}
              </h3>
              <p className="text-slate-400">
                {event.description}
              </p>
            </div>
          </div>
        ))}
      </div>

      {/* Time zone note */}
      <div className="mt-12 text-center">
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-slate-800/50 border border-slate-700">
          <svg className="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span className="text-slate-400 text-sm">
            All times are in Pacific Standard Time (UTC-8)
          </span>
        </div>
      </div>
    </section>
  );
}
