import sys
import os

scripts = [
    r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\check_git.sh",
    r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\fix_and_build.sh",
    r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\check_linux.sh",
]

for path in scripts:
    if os.path.exists(path):
        with open(path, "rb") as f:
            content = f.read()
        content = content.replace(b"\r\n", b"\n")
        with open(path, "wb") as f:
            f.write(content)
        print(f"Fixed: {path}")
    else:
        print(f"Skip: {path}")
