import shutil
src = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\internal\postgres\connection.go"
dst = r"E:\Repos\VIBESQL-LOCAL\repo\internal\postgres\connection.go"
shutil.copy2(src, dst)
print("copied")
