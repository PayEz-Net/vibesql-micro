import { NextResponse } from 'next/server';
import { appSubmissions, bugSubmissions } from '@/lib/vibesql';

export async function GET() {
  try {
    // Fetch app submissions
    const apps = await appSubmissions.getAll();

    // Fetch bug hunter leaderboard
    const bugHunters = await bugSubmissions.getLeaderboard();

    return NextResponse.json({
      apps: apps.success ? apps.data || [] : [],
      bugHunters: bugHunters.success ? bugHunters.data || [] : [],
    });
  } catch (error) {
    console.error('Leaderboard error:', error);
    return NextResponse.json({
      apps: [],
      bugHunters: [],
    });
  }
}
