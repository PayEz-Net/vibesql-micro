import Navigation from '@/components/Navigation';
import Footer from '@/components/Footer';
import Leaderboard from '@/components/Leaderboard';

export const metadata = {
  title: 'Leaderboard | VibeSQL Hackathon 2026',
  description: 'See all app submissions and top bug hunters in the VibeSQL Hackathon 2026.',
};

export default function LeaderboardPage() {
  return (
    <main className="min-h-screen flex flex-col">
      <Navigation />

      <div className="flex-1 pt-24 pb-16">
        <div className="max-w-4xl mx-auto px-4">
          {/* Header */}
          <div className="text-center mb-10">
            <h1 className="text-3xl md:text-4xl font-bold mb-4">
              <span className="text-vibe-accent">Leaderboard</span>
            </h1>
            <p className="text-slate-400">
              Track submissions and see who&apos;s leading the hackathon.
            </p>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-2 gap-4 mb-8">
            <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50 text-center">
              <div className="text-3xl font-bold text-vibe-accent mb-1">—</div>
              <div className="text-sm text-slate-400">Apps Submitted</div>
            </div>
            <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50 text-center">
              <div className="text-3xl font-bold text-hackathon-bug-high mb-1">—</div>
              <div className="text-sm text-slate-400">Bugs Reported</div>
            </div>
          </div>

          {/* Leaderboard */}
          <Leaderboard />
        </div>
      </div>

      <Footer />
    </main>
  );
}
