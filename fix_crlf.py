import os

files = [
    r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\kill_5173.sh",
]

for path in files:
    if os.path.exists(path):
        with open(path, "rb") as f:
            content = f.read()
        content = content.replace(b"\r\n", b"\n")
        with open(path, "wb") as f:
            f.write(content)
        print(f"Fixed: {os.path.basename(path)}")
