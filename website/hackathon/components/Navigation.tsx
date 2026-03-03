'use client';

import { useState, useEffect } from 'react';

export default function Navigation() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const navLinks = [
    { href: '#about', label: 'About' },
    { href: '#tracks', label: 'Tracks' },
    { href: '#prizes', label: 'Prizes' },
    { href: '#schedule', label: 'Schedule' },
    { href: '#rules', label: 'Rules' },
    { href: '#faq', label: 'FAQ' },
  ];

  return (
    <nav
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? 'bg-slate-900/95 backdrop-blur-md border-b border-slate-800'
          : 'bg-transparent'
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <a href="/" className="flex items-center gap-2">
            <svg width="28" height="28" viewBox="0 0 64 64" className="text-vibe-primary">
              <g transform="translate(32, 32)">
                <ellipse cx="0" cy="-16" rx="20" ry="7" fill="currentColor"/>
                <rect x="-20" y="-16" width="40" height="24" fill="currentColor"/>
                <ellipse cx="0" cy="8" rx="20" ry="7" fill="currentColor"/>
              </g>
            </svg>
            <span className="font-bold text-white">VibeSQL</span>
          </a>

          {/* Desktop navigation */}
          <div className="hidden md:flex items-center gap-8">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="text-sm text-slate-300 hover:text-white transition-colors"
              >
                {link.label}
              </a>
            ))}
            <a
              href="/register"
              className="bg-vibe-accent hover:bg-vibe-accent/90 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
            >
              Register
            </a>
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden p-2 text-slate-300 hover:text-white"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? (
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            ) : (
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            )}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      <div
        className={`md:hidden overflow-hidden transition-all duration-300 ${
          isMobileMenuOpen ? 'max-h-96' : 'max-h-0'
        }`}
      >
        <div className="bg-slate-900/95 backdrop-blur-md border-b border-slate-800 px-4 py-4 space-y-4">
          {navLinks.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="block text-slate-300 hover:text-white transition-colors"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              {link.label}
            </a>
          ))}
          <a
            href="/register"
            className="block w-full text-center bg-vibe-accent hover:bg-vibe-accent/90 text-white font-semibold px-4 py-2 rounded-lg transition-colors"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Register
          </a>
        </div>
      </div>
    </nav>
  );
}
