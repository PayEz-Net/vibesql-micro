import subprocess, os

REPO = r"E:\Repos\VIBESQL-LOCAL\repo"
REMOTE = "zenflow@10.0.0.93"
REMOTE_BASE = "/opt/vibesql"
SSH_KEY = r"C:\Users\jon-local\.ssh\zenflow_93"

def scp(local, remote_path):
    cmd = ["scp", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no", local, f"{REMOTE}:{remote_path}"]
    r = subprocess.run(cmd, capture_output=True, text=True)
    if r.returncode != 0:
        print(f"  FAIL: {r.stderr.strip()}")
    else:
        print(f"  OK: {remote_path}")

f = os.path.join(REPO, "internal", "postgres", "connection.go")
data = open(f, "rb").read().replace(b"\r\n", b"\n")
tmp = f + ".lf"
open(tmp, "wb").write(data)
scp(tmp, f"{REMOTE_BASE}/internal/postgres/connection.go")
os.remove(tmp)
print("done")
