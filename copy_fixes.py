import shutil, subprocess, os

SSH_KEY = r"C:\Users\jon-local\.ssh\zenflow_93"
REMOTE = "zenflow@10.0.0.93"
WORKTREE = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11"
REPO = r"E:\Repos\VIBESQL-LOCAL\repo"

files = [
    "internal/postgres/connection.go",
    "internal/postgres/connection_test.go",
    "internal/postgres/manager_mock_test.go",
]

for f in files:
    src = os.path.join(WORKTREE, f.replace("/", os.sep))
    dst_repo = os.path.join(REPO, f.replace("/", os.sep))
    dst_linux = f"/opt/vibesql/{f}"

    shutil.copy2(src, dst_repo)
    print(f"copied to repo: {f}")

    data = open(src, "rb").read().replace(b"\r\n", b"\n")
    tmp = src + ".lf"
    open(tmp, "wb").write(data)
    subprocess.run(["scp", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no", tmp, f"{REMOTE}:{dst_linux}"], capture_output=True)
    os.remove(tmp)
    print(f"synced to linux: {f}")

print("done")
