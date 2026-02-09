import shutil
src = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\.zenflow\tasks\vibe-sql-micro-server-be11\plan.md"
dst = r"E:\Repos\VIBESQL-LOCAL\repo\.zenflow\tasks\vibe-sql-micro-server-be11\plan.md"
shutil.copy2(src, dst)

src2 = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\internal\postgres\connection.go"
dst2 = r"E:\Repos\VIBESQL-LOCAL\repo\internal\postgres\connection.go"
shutil.copy2(src2, dst2)

src3 = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\internal\postgres\connection_test.go"
dst3 = r"E:\Repos\VIBESQL-LOCAL\repo\internal\postgres\connection_test.go"
shutil.copy2(src3, dst3)

src4 = r"C:\Users\jon-local\.zenflow\worktrees\vibe-sql-micro-server-be11\internal\postgres\manager_mock_test.go"
dst4 = r"E:\Repos\VIBESQL-LOCAL\repo\internal\postgres\manager_mock_test.go"
shutil.copy2(src4, dst4)

print("all copied")
