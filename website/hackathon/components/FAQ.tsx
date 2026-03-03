'use client';

import { useState } from 'react';

interface FAQItem {
  question: string;
  answer: string;
}

const faqs: FAQItem[] = [
  {
    question: 'What is VibeSQL?',
    answer: 'VibeSQL is PostgreSQL + JSONB + HTTP API in one command. It\'s a zero-config embedded database perfect for AI agents, edge computing, and rapid prototyping. Run `npx vibesql-micro` and you have a full PostgreSQL database with an HTTP API ready to go.',
  },
  {
    question: 'Do I need to be an expert to participate?',
    answer: 'Not at all! The hackathon is open to all skill levels. VibeSQL is designed to be beginner-friendly with instant setup. Whether you\'re a seasoned developer or just starting out, you can participate and learn.',
  },
  {
    question: 'Can I participate in both tracks?',
    answer: 'Yes! You can submit both an app project and bug reports. Many participants do both - build something cool while also helping improve VibeSQL.',
  },
  {
    question: 'Do I need a team?',
    answer: 'No, you can participate solo or with a team of up to 4 people. Teams must be formed before the hackathon starts, and all members must be registered.',
  },
  {
    question: 'What technologies can I use?',
    answer: 'You can use any programming language or framework you like! The only requirement is that your project must use VibeSQL as its primary database. Common choices include Python, JavaScript/TypeScript, Go, and Rust.',
  },
  {
    question: 'Can I start working on my project before the hackathon?',
    answer: 'No, all code must be written during the 54-hour hackathon period. You can plan, design, and research beforehand, but actual coding must happen during the event.',
  },
  {
    question: 'How do bug bounty rewards work?',
    answer: 'Each valid bug report earns points based on severity. At the end of the hackathon, verified bugs are rewarded with cash prizes. The top bug hunters on the leaderboard also receive bonus rewards.',
  },
  {
    question: 'What counts as a valid bug?',
    answer: 'Valid bugs include security vulnerabilities, crashes, data corruption, incorrect behavior, performance issues, and documentation errors. Duplicates of already-reported bugs don\'t count.',
  },
  {
    question: 'How will winners be selected?',
    answer: 'App track submissions are judged on innovation, technical execution, usefulness, and presentation. Bug track rankings are based on the number and severity of verified bugs found.',
  },
  {
    question: 'When will prizes be paid out?',
    answer: 'Prizes will be distributed within 2 weeks of the winners announcement. We\'ll reach out via email to coordinate payment details.',
  },
];

function FAQItem({ item }: { item: FAQItem }) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border-b border-slate-700/50 last:border-0">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full py-5 flex items-center justify-between text-left hover:text-vibe-accent transition-colors"
      >
        <span className="font-semibold pr-8">{item.question}</span>
        <svg
          className={`w-5 h-5 flex-shrink-0 transition-transform duration-200 ${
            isOpen ? 'rotate-180' : ''
          }`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>
      <div
        className={`overflow-hidden transition-all duration-300 ${
          isOpen ? 'max-h-96 pb-5' : 'max-h-0'
        }`}
      >
        <p className="text-slate-400">{item.answer}</p>
      </div>
    </div>
  );
}

export default function FAQ() {
  return (
    <section id="faq" className="section-container bg-slate-900/50">
      <div className="text-center mb-16">
        <h2 className="text-3xl md:text-5xl font-bold mb-4">
          Frequently Asked <span className="text-vibe-accent">Questions</span>
        </h2>
        <p className="text-slate-400 text-lg max-w-2xl mx-auto">
          Everything you need to know about participating in the hackathon.
        </p>
      </div>

      <div className="max-w-3xl mx-auto">
        <div className="rounded-xl bg-slate-800/30 border border-slate-700/50 px-6">
          {faqs.map((faq, index) => (
            <FAQItem key={index} item={faq} />
          ))}
        </div>
      </div>

      {/* Contact */}
      <div className="mt-12 text-center">
        <p className="text-slate-400 mb-4">
          Still have questions?
        </p>
        <a
          href="https://github.com/PayEz-Net/vibesql-micro/discussions"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-slate-700/50 hover:bg-slate-700 transition-colors"
        >
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          Ask on GitHub Discussions
        </a>
      </div>
    </section>
  );
}
