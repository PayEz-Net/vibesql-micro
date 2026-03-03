import { NextRequest, NextResponse } from 'next/server';
import { registrations, appSubmissions } from '@/lib/vibesql';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { email, repo_url, demo_video_url, writeup, live_demo_url, team_members } = body;

    // Validate required fields
    if (!email || !email.trim()) {
      return NextResponse.json(
        { error: 'Email is required' },
        { status: 400 }
      );
    }

    if (!repo_url || !repo_url.trim()) {
      return NextResponse.json(
        { error: 'Repository URL is required' },
        { status: 400 }
      );
    }

    if (!demo_video_url || !demo_video_url.trim()) {
      return NextResponse.json(
        { error: 'Demo video URL is required' },
        { status: 400 }
      );
    }

    if (!writeup || !writeup.trim()) {
      return NextResponse.json(
        { error: 'Project writeup is required' },
        { status: 400 }
      );
    }

    // Validate URLs
    const urlRegex = /^https?:\/\/.+/;
    if (!urlRegex.test(repo_url)) {
      return NextResponse.json(
        { error: 'Please enter a valid repository URL' },
        { status: 400 }
      );
    }

    if (!urlRegex.test(demo_video_url)) {
      return NextResponse.json(
        { error: 'Please enter a valid demo video URL' },
        { status: 400 }
      );
    }

    // Find registration by email
    const registration = await registrations.findByEmail(email.toLowerCase().trim());
    if (!registration.success || !registration.data || registration.data.length === 0) {
      return NextResponse.json(
        { error: 'No registration found with this email. Please register first.' },
        { status: 404 }
      );
    }

    const registrationId = registration.data[0].id;

    // Create submission
    const result = await appSubmissions.create({
      registration_id: registrationId,
      repo_url: repo_url.trim(),
      demo_video_url: demo_video_url.trim(),
      writeup: writeup.trim(),
      live_demo_url: live_demo_url?.trim() || '',
      team_members: team_members || [],
    });

    if (!result.success) {
      console.error('App submission failed:', result.error);
      return NextResponse.json(
        { error: 'Failed to submit. Please try again.' },
        { status: 500 }
      );
    }

    return NextResponse.json(
      {
        success: true,
        message: 'App submitted successfully!',
        data: result.data?.[0],
      },
      { status: 201 }
    );
  } catch (error) {
    console.error('App submission error:', error);
    return NextResponse.json(
      { error: 'An unexpected error occurred' },
      { status: 500 }
    );
  }
}

export async function GET() {
  try {
    const result = await appSubmissions.getAll();
    if (!result.success) {
      return NextResponse.json({ submissions: [] });
    }
    return NextResponse.json({ submissions: result.data || [] });
  } catch (error) {
    console.error('Error getting app submissions:', error);
    return NextResponse.json({ submissions: [] });
  }
}
