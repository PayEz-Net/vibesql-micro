'use client';

import { useState } from 'react';

interface FormData {
  email: string;
  repo_url: string;
  demo_video_url: string;
  writeup: string;
  live_demo_url: string;
  team_members: string;
}

export default function AppSubmissionForm() {
  const [formData, setFormData] = useState<FormData>({
    email: '',
    repo_url: '',
    demo_video_url: '',
    writeup: '',
    live_demo_url: '',
    team_members: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitStatus('idle');
    setErrorMessage('');

    try {
      const teamMembers = formData.team_members
        .split(',')
        .map(m => m.trim().replace('@', ''))
        .filter(m => m.length > 0);

      const response = await fetch('/api/submit-app', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ...formData,
          team_members: teamMembers,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        setSubmitStatus('success');
      } else {
        setSubmitStatus('error');
        setErrorMessage(data.error || 'Something went wrong. Please try again.');
      }
    } catch {
      setSubmitStatus('error');
      setErrorMessage('Failed to connect. Please try again later.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (submitStatus === 'success') {
    return (
      <div className="text-center py-12">
        <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-green-500/20 flex items-center justify-center">
          <svg className="w-10 h-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 className="text-2xl font-bold mb-2">Submission Received!</h3>
        <p className="text-slate-400 mb-6">
          Your app has been submitted for judging. Good luck!
        </p>
        <a href="/leaderboard" className="btn-primary inline-flex items-center gap-2">
          View Submissions
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
          </svg>
        </a>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Email for lookup */}
      <div>
        <label htmlFor="email" className="form-label">
          Registered Email <span className="text-red-400">*</span>
        </label>
        <input
          type="email"
          id="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          className="form-input"
          placeholder="The email you registered with"
          required
        />
        <p className="text-xs text-slate-500 mt-1">
          Use the same email address you used to register
        </p>
      </div>

      {/* GitHub Repo */}
      <div>
        <label htmlFor="repo_url" className="form-label">
          GitHub Repository URL <span className="text-red-400">*</span>
        </label>
        <input
          type="url"
          id="repo_url"
          value={formData.repo_url}
          onChange={(e) => setFormData({ ...formData, repo_url: e.target.value })}
          className="form-input"
          placeholder="https://github.com/username/project"
          required
        />
      </div>

      {/* Demo Video */}
      <div>
        <label htmlFor="demo_video_url" className="form-label">
          Demo Video URL <span className="text-red-400">*</span>
        </label>
        <input
          type="url"
          id="demo_video_url"
          value={formData.demo_video_url}
          onChange={(e) => setFormData({ ...formData, demo_video_url: e.target.value })}
          className="form-input"
          placeholder="https://youtube.com/watch?v=... or https://loom.com/..."
          required
        />
        <p className="text-xs text-slate-500 mt-1">
          YouTube or Loom link (max 5 minutes)
        </p>
      </div>

      {/* Live Demo URL */}
      <div>
        <label htmlFor="live_demo_url" className="form-label">
          Live Demo URL <span className="text-slate-500">(optional)</span>
        </label>
        <input
          type="url"
          id="live_demo_url"
          value={formData.live_demo_url}
          onChange={(e) => setFormData({ ...formData, live_demo_url: e.target.value })}
          className="form-input"
          placeholder="https://your-app.vercel.app"
        />
      </div>

      {/* Team Members */}
      <div>
        <label htmlFor="team_members" className="form-label">
          Team Members <span className="text-slate-500">(optional)</span>
        </label>
        <input
          type="text"
          id="team_members"
          value={formData.team_members}
          onChange={(e) => setFormData({ ...formData, team_members: e.target.value })}
          className="form-input"
          placeholder="@teammate1, @teammate2"
        />
        <p className="text-xs text-slate-500 mt-1">
          Comma-separated GitHub usernames of team members
        </p>
      </div>

      {/* Writeup */}
      <div>
        <label htmlFor="writeup" className="form-label">
          Project Writeup <span className="text-red-400">*</span>
        </label>
        <textarea
          id="writeup"
          value={formData.writeup}
          onChange={(e) => setFormData({ ...formData, writeup: e.target.value })}
          className="form-input min-h-[200px] resize-y"
          placeholder="Describe your project:
- What does it do?
- How does it use VibeSQL?
- What problem does it solve?
- Any challenges you overcame?"
          required
        />
        <p className="text-xs text-slate-500 mt-1">
          Markdown supported. Be clear about how you used VibeSQL.
        </p>
      </div>

      {/* Error message */}
      {submitStatus === 'error' && (
        <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400">
          {errorMessage}
        </div>
      )}

      {/* Submit */}
      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full btn-primary py-4 text-lg disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <span className="inline-flex items-center gap-2">
            <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Submitting...
          </span>
        ) : (
          'Submit App'
        )}
      </button>
    </form>
  );
}
