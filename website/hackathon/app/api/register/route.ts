import { NextRequest, NextResponse } from 'next/server';
import { registrations } from '@/lib/vibesql';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Validate required fields
    const { name, email, github_username, skills, track } = body;

    if (!name || !name.trim()) {
      return NextResponse.json(
        { error: 'Name is required' },
        { status: 400 }
      );
    }

    if (!email || !email.trim()) {
      return NextResponse.json(
        { error: 'Email is required' },
        { status: 400 }
      );
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      return NextResponse.json(
        { error: 'Please enter a valid email address' },
        { status: 400 }
      );
    }

    if (!github_username || !github_username.trim()) {
      return NextResponse.json(
        { error: 'GitHub username is required' },
        { status: 400 }
      );
    }

    // Check if email is already registered
    const existing = await registrations.findByEmail(email.toLowerCase().trim());
    if (existing.success && existing.data && existing.data.length > 0) {
      return NextResponse.json(
        { error: 'This email is already registered' },
        { status: 409 }
      );
    }

    // Create registration
    const result = await registrations.create({
      name: name.trim(),
      email: email.toLowerCase().trim(),
      github_username: github_username.trim().replace('@', ''),
      skills: skills || [],
      track: track || 'both',
    });

    if (!result.success) {
      console.error('Registration failed:', result.error);
      return NextResponse.json(
        { error: 'Failed to register. Please try again.' },
        { status: 500 }
      );
    }

    return NextResponse.json(
      {
        success: true,
        message: 'Registration successful!',
        data: result.data?.[0],
      },
      { status: 201 }
    );
  } catch (error) {
    console.error('Registration error:', error);
    return NextResponse.json(
      { error: 'An unexpected error occurred' },
      { status: 500 }
    );
  }
}

export async function GET() {
  try {
    const count = await registrations.count();
    return NextResponse.json({ count });
  } catch (error) {
    console.error('Error getting registration count:', error);
    return NextResponse.json({ count: 0 });
  }
}
