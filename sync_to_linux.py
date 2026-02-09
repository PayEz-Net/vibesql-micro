import subprocess, os, glob

REPO = r"E:\Repos\VIBESQL-LOCAL\repo"
REMOTE = "zenflow@10.0.0.93"
REMOTE_BASE = "/opt/vibesql"
SSH_KEY = r"C:\Users\jon-local\.ssh\zenflow_93"

dirs_to_sync = [
    "cmd/vibe",
    "internal/postgres",
    "internal/query",
    "internal/server",
    "internal/version",
    "scripts",
    "build",
]

files_to_sync = [
    "go.mod",
    "Makefile",
    "README.md",
]

def scp(local, remote_path):
    cmd = [
        "scp", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no",
        local, f"{REMOTE}:{remote_path}"
    ]
    r = subprocess.run(cmd, capture_output=True, text=True)
    if r.returncode != 0:
        print(f"  FAIL: {r.stderr.strip()}")
    else:
        print(f"  OK: {remote_path}")

def ssh(command):
    cmd = ["ssh", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no", REMOTE] + command.split()
    r = subprocess.run(cmd, capture_output=True, text=True)
    return r.stdout.strip()

def fix_crlf(content):
    return content.replace(b"\r\n", b"\n")

def sync_file(local_path, remote_path):
    with open(local_path, "rb") as f:
        data = fix_crlf(f.read())
    tmp = local_path + ".lf"
    with open(tmp, "wb") as f:
        f.write(data)
    scp(tmp, remote_path)
    os.remove(tmp)

print("Syncing files to Linux box...")

for d in dirs_to_sync:
    local_dir = os.path.join(REPO, d)
    if not os.path.isdir(local_dir):
        print(f"  SKIP dir: {d}")
        continue
    remote_dir = f"{REMOTE_BASE}/{d}"
    ssh(f"mkdir -p {remote_dir}")
    for root, dirs, files in os.walk(local_dir):
        for fn in files:
            local_file = os.path.join(root, fn)
            rel = os.path.relpath(local_file, REPO).replace("\\", "/")
            remote_file = f"{REMOTE_BASE}/{rel}"
            remote_parent = os.path.dirname(remote_file).replace("\\", "/")
            ssh(f"mkdir -p {remote_parent}")
            if fn.endswith((".go", ".sh", ".md", ".mod", ".json", ".yml", ".yaml", ".cmd", ".txt")):
                sync_file(local_file, remote_file)
            else:
                scp(local_file, remote_file)

for fn in files_to_sync:
    local_file = os.path.join(REPO, fn)
    if os.path.isfile(local_file):
        remote_file = f"{REMOTE_BASE}/{fn}"
        if fn.endswith((".go", ".sh", ".md", ".mod", ".json", ".yml", ".yaml", ".cmd", ".txt")):
            sync_file(local_file, remote_file)
        else:
            scp(local_file, remote_file)

print("\nDone! Synced all source files.")
