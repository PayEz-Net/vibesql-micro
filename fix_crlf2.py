import sys
src = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\rebuild.sh"
dst = src + ".lf"
data = open(src, "rb").read().replace(b"\r\n", b"\n")
open(dst, "wb").write(data)
print("done")
