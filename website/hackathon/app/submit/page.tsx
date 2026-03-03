'use client';

import { useState } from 'react';
import Navigation from '@/components/Navigation';
import Footer from '@/components/Footer';
import AppSubmissionForm from '@/components/AppSubmissionForm';
import BugSubmissionForm from '@/components/BugSubmissionForm';

export default function SubmitPage() {
  const [activeTab, setActiveTab] = useState<'app' | 'bug'>('app');

  return (
    <main className="min-h-screen flex flex-col">
      <Navigation />

      <div className="flex-1 pt-24 pb-16">
        <div className="max-w-2xl mx-auto px-4">
          {/* Header */}
          <div className="text-center mb-10">
            <h1 className="text-3xl md:text-4xl font-bold mb-4">
              Submit Your <span className="text-vibe-accent">Work</span>
            </h1>
            <p className="text-slate-400">
              Submit your app project or report bugs you&apos;ve found.
            </p>
          </div>

          {/* Tab switcher */}
          <div className="flex gap-2 mb-8 p-1 bg-slate-800/50 rounded-lg">
            <button
              onClick={() => setActiveTab('app')}
              className={`flex-1 py-3 px-6 rounded-md font-semibold transition-all flex items-center justify-center gap-2 ${
                activeTab === 'app'
                  ? 'bg-vibe-accent text-white'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
              App Submission
            </button>
            <button
              onClick={() => setActiveTab('bug')}
              className={`flex-1 py-3 px-6 rounded-md font-semibold transition-all flex items-center justify-center gap-2 ${
                activeTab === 'bug'
                  ? 'bg-hackathon-bug-high text-white'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              Bug Report
            </button>
          </div>

          {/* Form card */}
          <div className="rounded-2xl bg-slate-800/30 border border-slate-700/50 p-8">
            {activeTab === 'app' ? <AppSubmissionForm /> : <BugSubmissionForm />}
          </div>

          {/* Guidelines */}
          <div className="mt-8">
            {activeTab === 'app' ? (
              <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
                <h3 className="font-semibold mb-4 flex items-center gap-2">
                  <svg className="w-5 h-5 text-vibe-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Submission Checklist
                </h3>
                <ul className="space-y-2 text-sm text-slate-400">
                  <li className="flex items-start gap-2">
                    <span className="text-vibe-accent">•</span>
                    Repository is public and contains all source code
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-vibe-accent">•</span>
                    README explains how to run the project
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-vibe-accent">•</span>
                    Demo video shows the app in action (max 5 min)
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-vibe-accent">•</span>
                    Project clearly uses VibeSQL
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-vibe-accent">•</span>
                    All code was written during the hackathon
                  </li>
                </ul>
              </div>
            ) : (
              <div className="p-6 rounded-xl bg-slate-800/30 border border-slate-700/50">
                <h3 className="font-semibold mb-4 flex items-center gap-2">
                  <svg className="w-5 h-5 text-hackathon-bug-high" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Bug Report Tips
                </h3>
                <ul className="space-y-2 text-sm text-slate-400">
                  <li className="flex items-start gap-2">
                    <span className="text-hackathon-bug-high">•</span>
                    Include clear reproduction steps
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-hackathon-bug-high">•</span>
                    Provide minimal examples when possible
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-hackathon-bug-high">•</span>
                    Describe expected vs actual behavior
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-hackathon-bug-high">•</span>
                    Include environment details (OS, Node version)
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-hackathon-bug-high">•</span>
                    PRs with fixes earn extra consideration
                  </li>
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
