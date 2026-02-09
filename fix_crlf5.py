import subprocess, os
SSH_KEY = r"C:\Users\jon-local\.ssh\zenflow_93"
REMOTE = "zenflow@10.0.0.93"

for name in ["run_unit_tests.sh"]:
    src = os.path.join(r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11", name)
    data = open(src, "rb").read().replace(b"\r\n", b"\n")
    tmp = src + ".lf"
    open(tmp, "wb").write(data)
    subprocess.run(["scp", "-i", SSH_KEY, "-o", "StrictHostKeyChecking=no", tmp, f"{REMOTE}:/opt/vibesql/{name}"], capture_output=True)
    os.remove(tmp)
    print(f"uploaded {name}")
