import subprocess, os
SSH_KEY = r"C:\Users\jon-local\.ssh\zenflow_93"
REMOTE = "zenflow@10.0.0.93"
src = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\full_final.sh"
data = open(src, "rb").read().replace(b"\r\n", b"\n")
tmp = src + ".lf"
open(tmp, "wb").write(data)
subprocess.run(["scp", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no", tmp, f"{REMOTE}:/opt/vibesql/full_final.sh"], capture_output=True)
os.remove(tmp)
print("uploaded")
