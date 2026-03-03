# VibeSQL Hackathon 2026

A 54-hour online hackathon to build innovative applications with VibeSQL and hunt bugs for bounties.

**Event Dates:** March 20-22, 2026
**Duration:** 54 hours
**Format:** Online, worldwide

## Quick Links

- **Landing Page:** [vibesql.online/hackathon](https://vibesql.online/hackathon)
- **Starter Kit:** [github.com/PayEz-Net/vibesql-skills](https://github.com/PayEz-Net/vibesql-skills)
- **Documentation:** [vibesql.online/docs.html](https://vibesql.online/docs.html)
- **Source Code:** [github.com/PayEz-Net/vibesql-micro](https://github.com/PayEz-Net/vibesql-micro)

## Two Prize Tracks

### Track 1: Build Apps

Create innovative applications powered by VibeSQL. Show what's possible when databases are this easy!

**Ideas:**
- AI agents with persistent memory
- Local-first apps with sync
- Developer tools & CLI apps
- Data visualization dashboards
- Edge computing solutions

**Requirements:**
- Public GitHub repository
- Demo video (max 5 minutes)
- Project writeup
- Live demo URL (optional)

**Judging Criteria:**
- Innovation (25%) - Creative use of VibeSQL
- Technical Execution (25%) - Code quality, architecture
- Usefulness (25%) - Solves a real problem
- Presentation (25%) - Demo quality, documentation

### Track 2: Find Bugs

Hunt for bugs in VibeSQL and earn bounties. Critical bugs = bigger rewards!

**Bug Categories:**
| Severity | Description | Points |
|----------|-------------|--------|
| Critical | Security vulnerabilities, data loss | 100 |
| High | Major functionality issues | 50 |
| Medium | Performance, edge cases | 25 |
| Low | Minor issues, documentation | 10 |

**Requirements:**
- GitHub Issue with clear reproduction steps
- Expected vs actual behavior
- PRs with fixes earn extra consideration

## Schedule (PST)

| Date | Time | Event |
|------|------|-------|
| Mar 20 | 6:00 PM | Hackathon Kickoff |
| Mar 20 | 6:30 PM | Hacking Begins |
| Mar 21 | 12:00 PM | Mid-Hackathon Check-in |
| Mar 22 | 6:00 PM | Final Hours (6h remaining) |
| Mar 22 | 11:59 PM | Submissions Close |
| Mar 23 | All Day | Judging Period |
| Mar 24 | TBA | Winners Announced |

## Eligibility

- Open to developers worldwide, 18 years or older
- Individual or team participation (max 4 members)
- Must register before hackathon starts
- Projects must use VibeSQL as primary database
- All code must be written during the hackathon period

## Getting Started

1. **Register** at [vibesql.online/hackathon/register](https://vibesql.online/hackathon/register)

2. **Install VibeSQL:**
   ```bash
   npx vibesql-micro
   ```

3. **Get the Starter Kit:**
   ```bash
   git clone https://github.com/PayEz-Net/vibesql-skills
   ```

4. **Read the Docs:**
   [vibesql.online/docs.html](https://vibesql.online/docs.html)

## Development

### Running the Hackathon Website Locally

```bash
cd website/hackathon
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

### Environment Variables

Copy `.env.example` to `.env.local`:

```bash
cp .env.example .env.local
```

Configure your VibeSQL connection:

```
VIBESQL_URL=http://localhost:5432
VIBESQL_TOKEN=
```

### Database Schema

The hackathon uses three tables:

```sql
-- Registrations
CREATE TABLE registrations (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  github_username TEXT NOT NULL,
  skills JSONB DEFAULT '[]',
  track TEXT DEFAULT 'both',
  registered_at TIMESTAMP DEFAULT NOW()
);

-- App Submissions
CREATE TABLE app_submissions (
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
CREATE TABLE bug_submissions (
  id SERIAL PRIMARY KEY,
  registration_id INTEGER REFERENCES registrations(id),
  issue_pr_url TEXT NOT NULL,
  severity TEXT NOT NULL,
  description TEXT,
  points INTEGER DEFAULT 0,
  verified BOOLEAN DEFAULT FALSE,
  submitted_at TIMESTAMP DEFAULT NOW()
);
```

## Tech Stack

- **Frontend:** Next.js 15, React 19, TypeScript, TailwindCSS
- **Backend:** VibeSQL Micro (PostgreSQL via HTTP API)
- **Deployment:** Vercel or static export

## Project Structure

```
website/hackathon/
├── app/
│   ├── page.tsx              # Landing page
│   ├── register/page.tsx     # Registration form
│   ├── submit/page.tsx       # Submission portal
│   ├── leaderboard/page.tsx  # Leaderboard
│   └── api/
│       ├── register/route.ts
│       ├── submit-app/route.ts
│       ├── submit-bug/route.ts
│       └── leaderboard/route.ts
├── components/
│   ├── Hero.tsx
│   ├── Countdown.tsx
│   ├── Navigation.tsx
│   └── ...
├── lib/
│   └── vibesql.ts            # VibeSQL HTTP client
└── package.json
```

## FAQ

**Can I participate in both tracks?**
Yes! Many participants build apps and hunt bugs.

**Do I need a team?**
No, solo participation is welcome. Teams can have up to 4 members.

**What technologies can I use?**
Any language or framework. The only requirement is using VibeSQL as your primary database.

**When will prizes be paid out?**
Within 2 weeks of winners announcement.

## Support

- **GitHub Discussions:** [github.com/PayEz-Net/vibesql-micro/discussions](https://github.com/PayEz-Net/vibesql-micro/discussions)
- **Documentation:** [vibesql.online/docs.html](https://vibesql.online/docs.html)

---

Good luck and happy hacking!
