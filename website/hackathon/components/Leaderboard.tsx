'use client';

import { useEffect, useState } from 'react';

interface AppSubmission {
  id: number;
  name: string;
  github_username: string;
  repo_url: string;
  demo_video_url: string;
  live_demo_url?: string;
  submitted_at: string;
}

interface BugHunter {
  github_username: string;
  name: string;
  total_points: number;
  bug_count: number;
}

interface LeaderboardData {
  apps: AppSubmission[];
  bugHunters: BugHunter[];
}

function AppSubmissionCard({ submission, rank }: { submission: AppSubmission; rank: number }) {
  return (
    <div className="prize-card flex gap-4">
      <div className="flex-shrink-0 w-10 h-10 rounded-full bg-slate-700 flex items-center justify-center font-bold text-slate-300">
        {rank}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className="font-semibold truncate">{submission.name}</span>
          <a
            href={`https://github.com/${submission.github_username}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-slate-400 hover:text-vibe-accent text-sm"
          >
            @{submission.github_username}
          </a>
        </div>
        <div className="flex flex-wrap gap-2 text-sm">
          <a
            href={submission.repo_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-slate-400 hover:text-white"
          >
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Repo
          </a>
          <a
            href={submission.demo_video_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-slate-400 hover:text-white"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Demo
          </a>
          {submission.live_demo_url && (
            <a
              href={submission.live_demo_url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-vibe-accent hover:text-white"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              Live
            </a>
          )}
        </div>
      </div>
    </div>
  );
}

function BugHunterRow({ hunter, rank }: { hunter: BugHunter; rank: number }) {
  const getMedalColor = (r: number) => {
    if (r === 1) return 'text-hackathon-gold';
    if (r === 2) return 'text-hackathon-silver';
    if (r === 3) return 'text-hackathon-bronze';
    return 'text-slate-400';
  };

  return (
    <div className="flex items-center gap-4 p-4 rounded-lg bg-slate-800/30 border border-slate-700/50 hover:border-slate-600 transition-colors">
      <div className={`w-8 text-center font-bold ${getMedalColor(rank)}`}>
        {rank <= 3 ? (
          <svg className="w-6 h-6 mx-auto" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        ) : (
          <span>{rank}</span>
        )}
      </div>
      <div className="flex-1">
        <div className="font-semibold">{hunter.name}</div>
        <a
          href={`https://github.com/${hunter.github_username}`}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm text-slate-400 hover:text-vibe-accent"
        >
          @{hunter.github_username}
        </a>
      </div>
      <div className="text-right">
        <div className="text-xl font-bold text-vibe-accent">{hunter.total_points}</div>
        <div className="text-xs text-slate-500">{hunter.bug_count} bugs</div>
      </div>
    </div>
  );
}

export default function Leaderboard() {
  const [data, setData] = useState<LeaderboardData>({ apps: [], bugHunters: [] });
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'apps' | 'bugs'>('apps');

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch('/api/leaderboard');
        const result = await response.json();
        setData({
          apps: result.apps || [],
          bugHunters: result.bugHunters || [],
        });
      } catch (error) {
        console.error('Failed to fetch leaderboard:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchData();
    // Refresh every 30 seconds
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-16">
        <svg className="animate-spin h-8 w-8 text-vibe-accent" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
      </div>
    );
  }

  return (
    <div>
      {/* Tab switcher */}
      <div className="flex gap-2 mb-8 p-1 bg-slate-800/50 rounded-lg">
        <button
          onClick={() => setActiveTab('apps')}
          className={`flex-1 py-3 px-6 rounded-md font-semibold transition-all flex items-center justify-center gap-2 ${
            activeTab === 'apps'
              ? 'bg-vibe-accent text-white'
              : 'text-slate-400 hover:text-white'
          }`}
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
          App Submissions
          <span className="px-2 py-0.5 text-xs rounded-full bg-slate-700">
            {data.apps.length}
          </span>
        </button>
        <button
          onClick={() => setActiveTab('bugs')}
          className={`flex-1 py-3 px-6 rounded-md font-semibold transition-all flex items-center justify-center gap-2 ${
            activeTab === 'bugs'
              ? 'bg-hackathon-bug-high text-white'
              : 'text-slate-400 hover:text-white'
          }`}
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          Bug Hunters
          <span className="px-2 py-0.5 text-xs rounded-full bg-slate-700">
            {data.bugHunters.length}
          </span>
        </button>
      </div>

      {/* Content */}
      {activeTab === 'apps' ? (
        <div className="space-y-4">
          {data.apps.length === 0 ? (
            <div className="text-center py-12">
              <svg className="w-16 h-16 mx-auto text-slate-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
              <h3 className="text-xl font-semibold mb-2">No Submissions Yet</h3>
              <p className="text-slate-400 mb-6">
                Be the first to submit your app!
              </p>
              <a href="/submit" className="btn-primary">
                Submit Your App
              </a>
            </div>
          ) : (
            data.apps.map((app, i) => (
              <AppSubmissionCard key={app.id} submission={app} rank={i + 1} />
            ))
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {data.bugHunters.length === 0 ? (
            <div className="text-center py-12">
              <svg className="w-16 h-16 mx-auto text-slate-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <h3 className="text-xl font-semibold mb-2">No Bug Reports Yet</h3>
              <p className="text-slate-400 mb-6">
                Start hunting bugs to climb the leaderboard!
              </p>
              <a href="/submit" className="btn-primary">
                Report a Bug
              </a>
            </div>
          ) : (
            data.bugHunters.map((hunter, i) => (
              <BugHunterRow key={hunter.github_username} hunter={hunter} rank={i + 1} />
            ))
          )}
        </div>
      )}

      {/* Refresh indicator */}
      <p className="text-center text-xs text-slate-500 mt-8">
        Auto-refreshes every 30 seconds
      </p>
    </div>
  );
}
