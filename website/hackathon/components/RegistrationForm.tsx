'use client';

import { useState } from 'react';

const SKILLS = [
  { id: 'ai-ml', label: 'AI/ML' },
  { id: 'backend', label: 'Backend' },
  { id: 'frontend', label: 'Frontend' },
  { id: 'devops', label: 'DevOps' },
  { id: 'other', label: 'Other' },
];

const TRACKS = [
  { id: 'apps', label: 'Build Apps', description: 'Create innovative applications' },
  { id: 'bugs', label: 'Find Bugs', description: 'Hunt for bugs and earn bounties' },
  { id: 'both', label: 'Both Tracks', description: 'Participate in both' },
];

interface FormData {
  name: string;
  email: string;
  github_username: string;
  skills: string[];
  track: string;
}

export default function RegistrationForm() {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    email: '',
    github_username: '',
    skills: [],
    track: 'both',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSkillToggle = (skillId: string) => {
    setFormData((prev) => ({
      ...prev,
      skills: prev.skills.includes(skillId)
        ? prev.skills.filter((s) => s !== skillId)
        : [...prev.skills, skillId],
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitStatus('idle');
    setErrorMessage('');

    try {
      const response = await fetch('/api/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (response.ok) {
        setSubmitStatus('success');
        setFormData({
          name: '',
          email: '',
          github_username: '',
          skills: [],
          track: 'both',
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
        <h3 className="text-2xl font-bold mb-2">You&apos;re Registered!</h3>
        <p className="text-slate-400 mb-6">
          Check your email for confirmation and updates about the hackathon.
        </p>
        <div className="space-y-4">
          <a
            href="https://github.com/PayEz-Net/vibesql-skills"
            target="_blank"
            rel="noopener noreferrer"
            className="btn-primary inline-flex items-center gap-2"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Get the Starter Kit
          </a>
          <p className="text-sm text-slate-500">
            Start exploring VibeSQL before the hackathon begins!
          </p>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Name */}
      <div>
        <label htmlFor="name" className="form-label">
          Full Name <span className="text-red-400">*</span>
        </label>
        <input
          type="text"
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          className="form-input"
          placeholder="Enter your full name"
          required
        />
      </div>

      {/* Email */}
      <div>
        <label htmlFor="email" className="form-label">
          Email Address <span className="text-red-400">*</span>
        </label>
        <input
          type="email"
          id="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          className="form-input"
          placeholder="you@example.com"
          required
        />
      </div>

      {/* GitHub Username */}
      <div>
        <label htmlFor="github" className="form-label">
          GitHub Username <span className="text-red-400">*</span>
        </label>
        <div className="relative">
          <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500">@</span>
          <input
            type="text"
            id="github"
            value={formData.github_username}
            onChange={(e) => setFormData({ ...formData, github_username: e.target.value })}
            className="form-input pl-8"
            placeholder="username"
            required
          />
        </div>
      </div>

      {/* Skills */}
      <div>
        <label className="form-label">Skills & Interests</label>
        <div className="flex flex-wrap gap-2">
          {SKILLS.map((skill) => (
            <button
              key={skill.id}
              type="button"
              onClick={() => handleSkillToggle(skill.id)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                formData.skills.includes(skill.id)
                  ? 'bg-vibe-accent text-white'
                  : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
              }`}
            >
              {skill.label}
            </button>
          ))}
        </div>
      </div>

      {/* Track */}
      <div>
        <label className="form-label">Which track are you interested in?</label>
        <div className="space-y-3">
          {TRACKS.map((track) => (
            <label
              key={track.id}
              className={`flex items-center gap-4 p-4 rounded-lg cursor-pointer transition-all ${
                formData.track === track.id
                  ? 'bg-vibe-accent/10 border-2 border-vibe-accent'
                  : 'bg-slate-800/50 border-2 border-transparent hover:border-slate-600'
              }`}
            >
              <input
                type="radio"
                name="track"
                value={track.id}
                checked={formData.track === track.id}
                onChange={(e) => setFormData({ ...formData, track: e.target.value })}
                className="hidden"
              />
              <div
                className={`w-5 h-5 rounded-full border-2 flex items-center justify-center ${
                  formData.track === track.id ? 'border-vibe-accent' : 'border-slate-500'
                }`}
              >
                {formData.track === track.id && (
                  <div className="w-3 h-3 rounded-full bg-vibe-accent" />
                )}
              </div>
              <div>
                <div className="font-semibold">{track.label}</div>
                <div className="text-sm text-slate-400">{track.description}</div>
              </div>
            </label>
          ))}
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
            Registering...
          </span>
        ) : (
          'Register for Hackathon'
        )}
      </button>

      <p className="text-sm text-slate-500 text-center">
        By registering, you agree to the hackathon rules and code of conduct.
      </p>
    </form>
  );
}
