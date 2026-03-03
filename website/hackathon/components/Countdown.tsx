'use client';

import { useEffect, useState } from 'react';

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
}

const HACKATHON_START = new Date('2026-03-20T18:00:00-08:00'); // 6:00 PM PST
const HACKATHON_END = new Date('2026-03-22T23:59:00-08:00'); // 11:59 PM PST

function calculateTimeLeft(target: Date): TimeLeft {
  const difference = target.getTime() - new Date().getTime();

  if (difference <= 0) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0 };
  }

  return {
    days: Math.floor(difference / (1000 * 60 * 60 * 24)),
    hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((difference / 1000 / 60) % 60),
    seconds: Math.floor((difference / 1000) % 60),
  };
}

function TimeBlock({ value, label }: { value: number; label: string }) {
  return (
    <div className="flex flex-col items-center">
      <div className="relative">
        <div className="countdown-number text-4xl md:text-6xl font-bold bg-gradient-to-b from-white to-slate-400 bg-clip-text text-transparent">
          {String(value).padStart(2, '0')}
        </div>
        <div className="absolute inset-0 bg-vibe-accent/10 rounded-lg blur-xl -z-10" />
      </div>
      <span className="text-xs md:text-sm text-slate-400 uppercase tracking-wider mt-2">
        {label}
      </span>
    </div>
  );
}

function Separator() {
  return (
    <div className="text-3xl md:text-5xl font-bold text-vibe-accent animate-pulse">
      :
    </div>
  );
}

export default function Countdown() {
  const [mounted, setMounted] = useState(false);
  const [timeLeft, setTimeLeft] = useState<TimeLeft>({ days: 0, hours: 0, minutes: 0, seconds: 0 });
  const [status, setStatus] = useState<'upcoming' | 'live' | 'ended'>('upcoming');

  useEffect(() => {
    setMounted(true);

    const updateCountdown = () => {
      const now = new Date();

      if (now < HACKATHON_START) {
        setStatus('upcoming');
        setTimeLeft(calculateTimeLeft(HACKATHON_START));
      } else if (now < HACKATHON_END) {
        setStatus('live');
        setTimeLeft(calculateTimeLeft(HACKATHON_END));
      } else {
        setStatus('ended');
        setTimeLeft({ days: 0, hours: 0, minutes: 0, seconds: 0 });
      }
    };

    updateCountdown();
    const interval = setInterval(updateCountdown, 1000);

    return () => clearInterval(interval);
  }, []);

  if (!mounted) {
    return (
      <div className="flex items-center justify-center gap-4 md:gap-6 py-8">
        <TimeBlock value={0} label="Days" />
        <Separator />
        <TimeBlock value={0} label="Hours" />
        <Separator />
        <TimeBlock value={0} label="Minutes" />
        <Separator />
        <TimeBlock value={0} label="Seconds" />
      </div>
    );
  }

  return (
    <div className="text-center">
      {status === 'upcoming' && (
        <p className="text-sm text-slate-400 mb-4 uppercase tracking-widest">
          Hackathon starts in
        </p>
      )}
      {status === 'live' && (
        <div className="flex items-center justify-center gap-2 mb-4">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
          </span>
          <span className="text-red-400 font-semibold uppercase tracking-widest text-sm">
            Live Now! Time Remaining
          </span>
        </div>
      )}
      {status === 'ended' && (
        <p className="text-lg text-slate-400 mb-4">
          The hackathon has ended. Thank you for participating!
        </p>
      )}

      {status !== 'ended' && (
        <div className="flex items-center justify-center gap-4 md:gap-6">
          <TimeBlock value={timeLeft.days} label="Days" />
          <Separator />
          <TimeBlock value={timeLeft.hours} label="Hours" />
          <Separator />
          <TimeBlock value={timeLeft.minutes} label="Minutes" />
          <Separator />
          <TimeBlock value={timeLeft.seconds} label="Seconds" />
        </div>
      )}
    </div>
  );
}
