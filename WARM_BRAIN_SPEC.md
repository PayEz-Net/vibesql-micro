# Warm Brain Spec: Sunday Morning Easy

## The Feeling

Not "oh my god this is amazing" (too hot).  
Not "ugh why won't this work" (too cold).  
Just... "ah, that was nice" (warm).

Like:
- Pouring coffee that tastes right on first sip
- A door that opens smoothly without sticking
- A chair that fits just right
- Sunday morning, no alarms, soft light

## Implementation Rules

### 1. No Surprises

**Good:**
```bash
$ vsql-micro ./app.vsql
# (database opens, ready for queries)
```

**Bad:**
```bash
$ vsql-micro ./app.vsql
ERROR: Port 5432 in use
ERROR: Could not lock file
ERROR: PostgreSQL version mismatch
```

**Solution:**
- Auto-find available port
- Handle concurrent access gracefully
- Just work around problems quietly

### 2. Obvious Progress

**Good:**
```bash
$ vsql-micro ./newapp.vsql
Creating database... done
Ready.

vsql> 
```

**Bad:**
```bash
$ vsql-micro ./newapp.vsql
# (30 seconds of silence)
vsql> 
```

**Solution:**
- Show spinner for > 500ms operations
- Clear "done" confirmation
- Never leave user wondering

### 3. Gentle Errors

**Good:**
```bash
$ vsql-micro ./locked.vsql
That database is busy (another program is using it).
Try: vsql-micro ./other.vsql
```

**Bad:**
```bash
$ vsql-micro ./locked.vsql
FATAL: could not create lock file "postmaster.pid": File exists
HINT:  Is another postmaster (PID 1234) running?
```

**Solution:**
- Human-readable errors
- Suggest next steps
- No stack traces unless --verbose

### 4. Progressive Disclosure

**Simple first:**
```go
db, _ := vsql.Open("./app.vsql")
```

**Power available when needed:**
```go
db, _ := vsql.Open("./app.vsql",
    vsql.WithPort(0),        // Auto (default)
    vsql.WithCache("32MB"),  // Bigger cache
)
```

**Solution:**
- Defaults work for 90%
- Options for power users
- Don't clutter simple use case

---

## The Sunday Morning Test

Imagine it's Sunday, 9am. You're in pajamas, half-awake, making a side project. You want to add a database.

**Passes the test:**
```bash
# Install while coffee brews
go install github.com/vibesql/vibesql-micro/cmd/vsql-micro@latest

# Try it
cd ~/Projects/side-thing
vsql-micro ./data.vsql

# Works immediately
vsql> CREATE TABLE notes (id SERIAL, text TEXT);
vsql> INSERT INTO notes VALUES (1, 'buy milk');
vsql> SELECT * FROM notes;
 id |  text   
----+---------
  1 | buy milk

vsql> \q

# Back to coding
```

**Total cognitive load:** Near zero  
**Emotional state:** Calm, pleased  
**Confidence:** High

---

## Warm Brain UX Checklist

### Installation
- [ ] `go install` in < 20 seconds
- [ ] No other steps needed
- [ ] First run feels familiar

### CLI
- [ ] Commands are obvious (`init`, `shell`, `query`)
- [ ] Output is readable (tables aligned)
- [ ] Exit cleanly (`\q` or Ctrl+D)
- [ ] Help is actually helpful

### Go API
- [ ] Import path is short
- [ ] `Open()` does the right thing
- [ ] Errors explain what to do
- [ ] Docs show 5-line example first

### Documentation
- [ ] README has copy-paste example
- [ ] No required reading before using
- [ ] Advanced topics are clearly marked "optional"

---

## Anti-Patterns to Avoid

### ❌ Configuration Files
```yaml
# vibesql.yaml - NO!
port: 5432
max_connections: 10
shared_buffers: 12MB
```

### ❌ Environment Variables Required
```bash
export VIBESQL_PORT=5432
export VIBESQL_DATA_DIR=/var/lib/vibesql
# NO!
```

### ❌ Complex Setup
```bash
# First download postgres...
# Then extract...
# Then set permissions...
# Then configure...
# NO!
```

### ❌ Verbose Output
```
[DEBUG] Opening database
[DEBUG] Checking port 5432
[DEBUG] Port available
[DEBUG] Extracting binary
[DEBUG] Binary extracted
[DEBUG] Starting postgres
[DEBUG] Postgres started
# Too much noise!
```

---

## Warm Brain Implementation

### The Golden Path

```go
// File: pkg/vsql/open.go

package vsql

import (
    "fmt"
    "os"
    "path/filepath"
)

// Open creates or opens a database.
// It's designed to never fail on the happy path.
func Open(path string) (*DB, error) {
    // 1. Resolve path
    if path == "" {
        path = "./default.vsql"
    }
    
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("invalid path: %w", err)
    }
    
    // 2. Check if locked (gentle error)
    if isLocked(absPath) {
        return nil, fmt.Errorf(
            "%s is already open in another program\n"+
            "Try a different file, or close the other program",
            filepath.Base(path),
        )
    }
    
    // 3. Extract binary (with progress if first time)
    binPath, firstTime, err := ensureBinary()
    if err != nil {
        return nil, fmt.Errorf("setup failed: %w", err)
    }
    
    if firstTime {
        fmt.Fprintln(os.Stderr, "Setting up vibesql... done")
    }
    
    // 4. Initialize database if new
    dataDir := absPath + ".data"
    if _, err := os.Stat(dataDir); os.IsNotExist(err) {
        fmt.Fprintln(os.Stderr, "Creating database... done")
    }
    
    // 5. Start postgres (auto-find port)
    port, err := findFreePort()
    if err != nil {
        return nil, err
    }
    
    // ... rest of implementation
    
    return &DB{/* ... */}, nil
}
```

### Gentle Errors

```go
// File: pkg/vsql/errors.go

package vsql

import "fmt"

// WarmError wraps errors with helpful context
type WarmError struct {
    What   string
    Why    string
    Fix    string
    Err    error
}

func (e WarmError) Error() string {
    msg := e.What
    if e.Why != "" {
        msg += "\n  Why: " + e.Why
    }
    if e.Fix != "" {
        msg += "\n  Fix: " + e.Fix
    }
    return msg
}

// Example usage:
func openDB(path string) (*DB, error) {
    if isLocked(path) {
        return nil, WarmError{
            What: fmt.Sprintf("%s is busy", filepath.Base(path)),
            Why:  "Another program is using it",
            Fix:  "Close the other program, or use a different file",
        }
    }
    // ...
}
```

### CLI Output

```go
// File: cmd/vsql-micro/main.go

package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("vsql-micro - PostgreSQL that feels like SQLite")
        fmt.Println()
        fmt.Println("Usage:")
        fmt.Println("  vsql-micro <file.vsql>           # Open interactive shell")
        fmt.Println("  vsql-micro <file.vsql> <query>   # Run single query")
        fmt.Println()
        fmt.Println("Examples:")
        fmt.Println("  vsql-micro ./app.vsql")
        fmt.Println("  vsql-micro ./app.vsql 'SELECT * FROM users'")
        os.Exit(0)
    }
    
    path := os.Args[1]
    
    // Open with progress
    fmt.Fprint(os.Stderr, "Opening database... ")
    db, err := vsql.Open(path)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error")
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Fprintln(os.Stderr, "done")
    defer db.Close()
    
    // Single query mode
    if len(os.Args) > 2 {
        query := os.Args[2]
        rows, err := db.Query(query)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        printRows(rows)
        return
    }
    
    // Interactive shell
    runShell(db)
}
```

---

## Measuring Warmth

Ask beta testers:

1. "How frustrated did you feel during setup?"
   - 1 = Very frustrated (cold)
   - 5 = Very pleased (warm)

2. "How long until you had your first successful query?"
   - Target: < 2 minutes

3. "How many times did you read documentation?"
   - Target: 0 for basic use

4. "Would you recommend this to a friend making a weekend project?"
   - Target: 90% yes

---

## The Warm Brain Promise

> "It doesn't amaze you. It doesn't frustrate you.  
> It just fits into your Sunday morning like it was always there."

That's the goal.
