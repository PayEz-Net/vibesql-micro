/**
 * VibeSQL HTTP Client
 *
 * A lightweight client for interacting with VibeSQL Micro's HTTP API.
 * Documentation: https://vibesql.online/docs.html
 * Starter Kit: https://github.com/PayEz-Net/vibesql-skills
 */

const VIBESQL_URL = process.env.VIBESQL_URL || 'http://localhost:5432';
const VIBESQL_TOKEN = process.env.VIBESQL_TOKEN || '';

interface QueryResult<T = Record<string, unknown>> {
  success: boolean;
  data?: T[];
  error?: string;
  rowCount?: number;
}

interface ExecuteResult {
  success: boolean;
  rowCount?: number;
  error?: string;
}

/**
 * Execute a SQL query and return results
 */
export async function query<T = Record<string, unknown>>(
  sql: string,
  params?: unknown[]
): Promise<QueryResult<T>> {
  try {
    const response = await fetch(`${VIBESQL_URL}/query`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(VIBESQL_TOKEN && { 'Authorization': `Bearer ${VIBESQL_TOKEN}` }),
      },
      body: JSON.stringify({ sql, params }),
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const data = await response.json();
    return { success: true, data: data.rows, rowCount: data.rowCount };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

/**
 * Execute a SQL statement (INSERT, UPDATE, DELETE)
 */
export async function execute(
  sql: string,
  params?: unknown[]
): Promise<ExecuteResult> {
  try {
    const response = await fetch(`${VIBESQL_URL}/execute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(VIBESQL_TOKEN && { 'Authorization': `Bearer ${VIBESQL_TOKEN}` }),
      },
      body: JSON.stringify({ sql, params }),
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const data = await response.json();
    return { success: true, rowCount: data.rowCount };
  } catch (error) {
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

/**
 * Initialize database tables for hackathon
 */
export async function initHackathonTables(): Promise<ExecuteResult> {
  const sql = `
    -- Registrations
    CREATE TABLE IF NOT EXISTS registrations (
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      email TEXT UNIQUE NOT NULL,
      github_username TEXT NOT NULL,
      skills JSONB DEFAULT '[]',
      track TEXT DEFAULT 'both',
      registered_at TIMESTAMP DEFAULT NOW()
    );

    -- App Submissions
    CREATE TABLE IF NOT EXISTS app_submissions (
      id SERIAL PRIMARY KEY,
      registration_id INTEGER REFERENCES registrations(id),
      repo_url TEXT NOT NULL,
      demo_video_url TEXT NOT NULL,
      writeup TEXT,
      live_demo_url TEXT,
      team_members JSONB DEFAULT '[]',
      submitted_at TIMESTAMP DEFAULT NOW()
    );

    -- Bug Submissions
    CREATE TABLE IF NOT EXISTS bug_submissions (
      id SERIAL PRIMARY KEY,
      registration_id INTEGER REFERENCES registrations(id),
      issue_pr_url TEXT NOT NULL,
      severity TEXT NOT NULL,
      description TEXT,
      points INTEGER DEFAULT 0,
      verified BOOLEAN DEFAULT FALSE,
      submitted_at TIMESTAMP DEFAULT NOW()
    );
  `;

  return execute(sql);
}

// Type definitions for hackathon data
export interface Registration {
  id: number;
  name: string;
  email: string;
  github_username: string;
  skills: string[];
  track: 'apps' | 'bugs' | 'both';
  registered_at: string;
}

export interface AppSubmission {
  id: number;
  registration_id: number;
  repo_url: string;
  demo_video_url: string;
  writeup?: string;
  live_demo_url?: string;
  team_members: string[];
  submitted_at: string;
}

export interface BugSubmission {
  id: number;
  registration_id: number;
  issue_pr_url: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  description?: string;
  points: number;
  verified: boolean;
  submitted_at: string;
}

// CRUD operations for registrations
export const registrations = {
  async create(data: Omit<Registration, 'id' | 'registered_at'>): Promise<QueryResult<Registration>> {
    return query<Registration>(
      `INSERT INTO registrations (name, email, github_username, skills, track)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING *`,
      [data.name, data.email, data.github_username, JSON.stringify(data.skills), data.track]
    );
  },

  async findByEmail(email: string): Promise<QueryResult<Registration>> {
    return query<Registration>(
      'SELECT * FROM registrations WHERE email = $1',
      [email]
    );
  },

  async getAll(): Promise<QueryResult<Registration>> {
    return query<Registration>('SELECT * FROM registrations ORDER BY registered_at DESC');
  },

  async count(): Promise<number> {
    const result = await query<{ count: string }>('SELECT COUNT(*) as count FROM registrations');
    return result.success && result.data ? parseInt(result.data[0].count) : 0;
  },
};

// CRUD operations for app submissions
export const appSubmissions = {
  async create(data: Omit<AppSubmission, 'id' | 'submitted_at'>): Promise<QueryResult<AppSubmission>> {
    return query<AppSubmission>(
      `INSERT INTO app_submissions (registration_id, repo_url, demo_video_url, writeup, live_demo_url, team_members)
       VALUES ($1, $2, $3, $4, $5, $6)
       RETURNING *`,
      [
        data.registration_id,
        data.repo_url,
        data.demo_video_url,
        data.writeup || null,
        data.live_demo_url || null,
        JSON.stringify(data.team_members),
      ]
    );
  },

  async getAll(): Promise<QueryResult<AppSubmission & { name: string; github_username: string }>> {
    return query(
      `SELECT a.*, r.name, r.github_username
       FROM app_submissions a
       JOIN registrations r ON a.registration_id = r.id
       ORDER BY a.submitted_at DESC`
    );
  },

  async count(): Promise<number> {
    const result = await query<{ count: string }>('SELECT COUNT(*) as count FROM app_submissions');
    return result.success && result.data ? parseInt(result.data[0].count) : 0;
  },
};

// CRUD operations for bug submissions
export const bugSubmissions = {
  async create(data: Omit<BugSubmission, 'id' | 'points' | 'verified' | 'submitted_at'>): Promise<QueryResult<BugSubmission>> {
    const points = {
      critical: 100,
      high: 50,
      medium: 25,
      low: 10,
    }[data.severity] || 0;

    return query<BugSubmission>(
      `INSERT INTO bug_submissions (registration_id, issue_pr_url, severity, description, points)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING *`,
      [data.registration_id, data.issue_pr_url, data.severity, data.description || null, points]
    );
  },

  async getAll(): Promise<QueryResult<BugSubmission & { name: string; github_username: string }>> {
    return query(
      `SELECT b.*, r.name, r.github_username
       FROM bug_submissions b
       JOIN registrations r ON b.registration_id = r.id
       ORDER BY b.submitted_at DESC`
    );
  },

  async getLeaderboard(): Promise<QueryResult<{ github_username: string; name: string; total_points: number; bug_count: number }>> {
    return query(
      `SELECT r.github_username, r.name,
              COALESCE(SUM(b.points), 0) as total_points,
              COUNT(b.id) as bug_count
       FROM registrations r
       LEFT JOIN bug_submissions b ON r.id = b.registration_id AND b.verified = true
       GROUP BY r.id, r.github_username, r.name
       HAVING COUNT(b.id) > 0
       ORDER BY total_points DESC, bug_count DESC
       LIMIT 20`
    );
  },

  async count(): Promise<number> {
    const result = await query<{ count: string }>('SELECT COUNT(*) as count FROM bug_submissions');
    return result.success && result.data ? parseInt(result.data[0].count) : 0;
  },
};
