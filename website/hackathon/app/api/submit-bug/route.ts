import { NextRequest, NextResponse } from 'next/server';
import { registrations, bugSubmissions } from '@/lib/vibesql';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { email, issue_pr_url, severity, description } = body;

    // Validate required fields
    if (!email || !email.trim()) {
      return NextResponse.json(
        { error: 'Email is required' },
        { status: 400 }
      );
    }

    if (!issue_pr_url || !issue_pr_url.trim()) {
      return NextResponse.json(
        { error: 'GitHub Issue/PR URL is required' },
        { status: 400 }
      );
    }

    if (!severity) {
      return NextResponse.json(
        { error: 'Bug severity is required' },
        { status: 400 }
      );
    }

    // Validate severity
    const validSeverities = ['critical', 'high', 'medium', 'low'];
    if (!validSeverities.includes(severity)) {
      return NextResponse.json(
        { error: 'Invalid severity level' },
        { status: 400 }
      );
    }

    // Validate URL
    const urlRegex = /^https?:\/\/github\.com\/.+/;
    if (!urlRegex.test(issue_pr_url)) {
      return NextResponse.json(
        { error: 'Please enter a valid GitHub issue or PR URL' },
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
    const result = await bugSubmissions.create({
      registration_id: registrationId,
      issue_pr_url: issue_pr_url.trim(),
      severity: severity as 'critical' | 'high' | 'medium' | 'low',
      description: description?.trim(),
    });

    if (!result.success) {
      console.error('Bug submission failed:', result.error);
      return NextResponse.json(
        { error: 'Failed to submit. Please try again.' },
        { status: 500 }
      );
    }

    return NextResponse.json(
      {
        success: true,
        message: 'Bug report submitted successfully!',
        data: result.data?.[0],
      },
      { status: 201 }
    );
  } catch (error) {
    console.error('Bug submission error:', error);
    return NextResponse.json(
      { error: 'An unexpected error occurred' },
      { status: 500 }
    );
  }
}

export async function GET() {
  try {
    const result = await bugSubmissions.getAll();
    if (!result.success) {
      return NextResponse.json({ submissions: [] });
    }
    return NextResponse.json({ submissions: result.data || [] });
  } catch (error) {
    console.error('Error getting bug submissions:', error);
    return NextResponse.json({ submissions: [] });
  }
}
