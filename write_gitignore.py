import sys

content = """# VibeSQL Build Artifacts
build/output/
build/postgres_micro*
*.exe
/vibe
/vibe-*

# Go artifacts
*.o
*.a
*.test
go.sum

# Coverage reports
coverage.out
coverage*.out
coverage.html
coverage_summary.txt
*.coverprofile

# IDE and editor files
.vs/
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# PostgreSQL data directory
vibe-data/
*.pid

# Temporary files
tmp/
temp/
*.tmp
*.log

# Docker volumes
postgres-data/

# Binaries for all platforms
vibe-linux-amd64
vibe-linux-arm64
vibe-darwin-amd64
vibe-darwin-arm64
vibe-windows-amd64.exe

# Embedded postgres binaries (downloaded during build, not in git)
internal/postgres/embed/postgres_micro_*
internal/postgres/embed/initdb_linux_amd64
internal/postgres/embed/pg_ctl_linux_amd64
internal/postgres/embed/libpq.so.5
internal/postgres/embed/share.tar.gz

# Test binaries
*.test

# Dependency directories
vendor/
"""

target = sys.argv[1] if len(sys.argv) > 1 else r"E:\Repos\VIBESQL-LOCAL\repo\.gitignore"
with open(target, "w", newline="\n") as f:
    f.write(content)
print(f"Written to {target}")
