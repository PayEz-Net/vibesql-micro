export default function About() {
  return (
    <section id="about" className="section-container">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          What is <span className="text-vibe-accent">VibeSQL</span>?
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          PostgreSQL + JSONB + HTTP API in one command. Zero-config database for AI agents and modern apps.
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-8 lg:gap-12 items-center">
        {/* Code example */}
        <div className="code-block text-sm overflow-x-auto">
          <div className="flex items-center gap-2 mb-4 text-slate-400">
            <span className="w-3 h-3 rounded-full bg-red-500/80"></span>
            <span className="w-3 h-3 rounded-full bg-yellow-500/80"></span>
            <span className="w-3 h-3 rounded-full bg-green-500/80"></span>
            <span className="ml-2 text-xs">terminal</span>
          </div>
          <pre className="text-green-400">
{`$ npx vibesql-micro

🚀 VibeSQL Micro started!
   PostgreSQL: localhost:5432
   HTTP API:   localhost:8080

# Store data with HTTP
$ curl -X POST localhost:8080/query \\
    -d '{"sql": "INSERT INTO users
         VALUES ($1, $2)",
         "params": ["alice", {"role": "dev"}]}'

# Query with JSON
$ curl localhost:8080/query \\
    -d '{"sql": "SELECT * FROM users
         WHERE data->>'role' = 'dev'"}'`}
          </pre>
        </div>

        {/* Features list */}
        <div className="space-y-6">
          <div className="flex gap-4">
            <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-6 h-6 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <h3 className="text-xl font-semibold mb-1">Instant Setup</h3>
              <p className="text-slate-400">One command to start. No Docker, no config files, no installation headaches.</p>
            </div>
          </div>

          <div className="flex gap-4">
            <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-6 h-6 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
              </svg>
            </div>
            <div>
              <h3 className="text-xl font-semibold mb-1">Full PostgreSQL</h3>
              <p className="text-slate-400">PostgreSQL 16.1 with JSONB, full-text search, and all the features you love.</p>
            </div>
          </div>

          <div className="flex gap-4">
            <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-6 h-6 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h3 className="text-xl font-semibold mb-1">HTTP API Built-in</h3>
              <p className="text-slate-400">RESTful API for queries. Perfect for AI agents, serverless functions, and quick prototypes.</p>
            </div>
          </div>

          <div className="flex gap-4">
            <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-vibe-accent/10 flex items-center justify-center">
              <svg className="w-6 h-6 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h3 className="text-xl font-semibold mb-1">AI-Ready</h3>
              <p className="text-slate-400">Designed for AI agents to store memories, manage state, and persist data without setup friction.</p>
            </div>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="mt-16 grid grid-cols-2 md:grid-cols-4 gap-6">
        {[
          { value: '54', label: 'Hours of Hacking' },
          { value: '2', label: 'Prize Tracks' },
          { value: '$TBD', label: 'Total Prizes' },
          { value: '100%', label: 'Remote' },
        ].map((stat, i) => (
          <div key={i} className="text-center p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
            <div className="text-3xl md:text-4xl font-bold text-vibe-accent mb-1">{stat.value}</div>
            <div className="text-sm text-slate-400">{stat.label}</div>
          </div>
        ))}
      </div>
    </section>
  );
}
