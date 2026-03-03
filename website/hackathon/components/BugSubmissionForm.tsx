'use client';

import { useState } from 'react';

const SEVERITIES = [
  {
    id: 'critical',
    label: 'Critical',
    description: 'Security vulnerabilities, data loss, complete breakage',
    color: 'text-hackathon-bug-critical border-hackathon-bug-critical',
    bg: 'bg-hackathon-bug-critical/10',
  },
  {
    id: 'high',
    label: 'High',
    description: 'Major functionality issues, significant impact',
    color: 'text-hackathon-bug-high border-hackathon-bug-high',
    bg: 'bg-hackathon-bug-high/10',
  },
  {
    id: 'medium',
    label: 'Medium',
    description: 'Performance issues, edge cases, moderate impact',
    color: 'text-hackathon-bug-medium border-hackathon-bug-medium',
    bg: 'bg-hackathon-bug-medium/10',
  },
  {
    id: 'low',
    label: 'Low',
    description: 'Minor issues, documentation errors, cosmetic bugs',
    color: 'text-hackathon-bug-low border-hackathon-bug-low',
    bg: 'bg-hackathon-bug-low/10',
  },
];

interface FormData {
  email: string;
  issue_pr_url: string;
  severity: string;
  description: string;
}

export default function BugSubmissionForm() {
  const [formData, setFormData] = useState<FormData>({
    email: '',
    issue_pr_url: '',
    severity: '',
    description: '',
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
      const response = await fetch('/api/submit-bug', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (response.ok) {
        setSubmitStatus('success');
        // Reset form for another submission
        setFormData({
          email: formData.email, // Keep email
          issue_pr_url: '',
          severity: '',
          description: '',
        });
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
        <h3 className="text-2xl font-bold mb-2">Bug Report Submitted!</h3>
        <p className="text-slate-400 mb-6">
          Your bug has been logged. Once verified, points will be added to the leaderboard.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <button
            onClick={() => setSubmitStatus('idle')}
            className="btn-primary inline-flex items-center gap-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            Submit Another Bug
          </button>
          <a href="/leaderboard" className="btn-secondary inline-flex items-center gap-2">
            View Leaderboard
          </a>
        </div>
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
      </div>

      {/* Issue/PR URL */}
      <div>
        <label htmlFor="issue_pr_url" className="form-label">
          GitHub Issue or PR URL <span className="text-red-400">*</span>
        </label>
        <input
          type="url"
          id="issue_pr_url"
          value={formData.issue_pr_url}
          onChange={(e) => setFormData({ ...formData, issue_pr_url: e.target.value })}
          className="form-input"
          placeholder="https://github.com/PayEz-Net/vibesql-micro/issues/123"
          required
        />
        <p className="text-xs text-slate-500 mt-1">
          Link to the GitHub issue or pull request you created
        </p>
      </div>

      {/* Severity */}
      <div>
        <label className="form-label">
          Bug Severity <span className="text-red-400">*</span>
        </label>
        <div className="grid grid-cols-2 gap-3">
          {SEVERITIES.map((severity) => (
            <label
              key={severity.id}
              className={`flex flex-col p-4 rounded-lg cursor-pointer transition-all border-2 ${
                formData.severity === severity.id
                  ? `${severity.bg} ${severity.color}`
                  : 'bg-slate-800/50 border-transparent hover:border-slate-600'
              }`}
            >
              <input
                type="radio"
                name="severity"
                value={severity.id}
                checked={formData.severity === severity.id}
                onChange={(e) => setFormData({ ...formData, severity: e.target.value })}
                className="hidden"
                required
              />
              <span className="font-semibold">{severity.label}</span>
              <span className="text-xs text-slate-400 mt-1">{severity.description}</span>
            </label>
          ))}
        </div>
      </div>

      {/* Description */}
      <div>
        <label htmlFor="description" className="form-label">
          Bug Description <span className="text-slate-500">(optional but recommended)</span>
        </label>
        <textarea
          id="description"
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          className="form-input min-h-[120px] resize-y"
          placeholder="Brief description of the bug and its impact..."
        />
        <p className="text-xs text-slate-500 mt-1">
          Include any context not in the GitHub issue
        </p>
      </div>

      {/* Points info */}
      <div className="p-4 rounded-lg bg-slate-800/50 border border-slate-700">
        <h4 className="font-semibold mb-2 text-sm">Point Values</h4>
        <div className="grid grid-cols-4 gap-2 text-sm">
          <div className="text-center">
            <span className="text-hackathon-bug-critical font-bold">100</span>
            <p className="text-xs text-slate-500">Critical</p>
          </div>
          <div className="text-center">
            <span className="text-hackathon-bug-high font-bold">50</span>
            <p className="text-xs text-slate-500">High</p>
          </div>
          <div className="text-center">
            <span className="text-hackathon-bug-medium font-bold">25</span>
            <p className="text-xs text-slate-500">Medium</p>
          </div>
          <div className="text-center">
            <span className="text-hackathon-bug-low font-bold">10</span>
            <p className="text-xs text-slate-500">Low</p>
          </div>
        </div>
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
          'Submit Bug Report'
        )}
      </button>
    </form>
  );
}
