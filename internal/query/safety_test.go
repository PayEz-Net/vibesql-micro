package query

import (
	"testing"

	"github.com/vibesql/vibe/internal/postgres"
)

func TestCheckSafety_UpdateWithoutWhere(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "UPDATE without WHERE",
			sql:     "UPDATE users SET name = 'Alice'",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "UPDATE with WHERE",
			sql:     "UPDATE users SET name = 'Alice' WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "UPDATE with WHERE 1=1",
			sql:     "UPDATE users SET name = 'Alice' WHERE 1=1",
			wantErr: false,
		},
		{
			name:    "UPDATE with complex WHERE",
			sql:     "UPDATE users SET name = 'Alice' WHERE id > 5 AND status = 'active'",
			wantErr: false,
		},
		{
			name:    "lowercase update without where",
			sql:     "update users set name = 'Alice'",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "lowercase update with where",
			sql:     "update users set name = 'Alice' where id = 1",
			wantErr: false,
		},
		{
			name:    "mixed case UPDATE without WHERE",
			sql:     "UpDaTe users SET name = 'Alice'",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSafety(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, vibeErr.Code)
				}
			}
		})
	}
}

func TestCheckSafety_DeleteWithoutWhere(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "DELETE without WHERE",
			sql:     "DELETE FROM users",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "DELETE with WHERE",
			sql:     "DELETE FROM users WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "DELETE with WHERE 1=1",
			sql:     "DELETE FROM users WHERE 1=1",
			wantErr: false,
		},
		{
			name:    "DELETE with complex WHERE",
			sql:     "DELETE FROM users WHERE id > 5 AND status = 'inactive'",
			wantErr: false,
		},
		{
			name:    "lowercase delete without where",
			sql:     "delete from users",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "lowercase delete with where",
			sql:     "delete from users where id = 1",
			wantErr: false,
		},
		{
			name:    "mixed case DELETE without WHERE",
			sql:     "DeLeTe FROM users",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSafety(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, vibeErr.Code)
				}
			}
		})
	}
}

func TestCheckSafety_SafeQueries(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "SELECT",
			sql:  "SELECT * FROM users",
		},
		{
			name: "INSERT",
			sql:  "INSERT INTO users (name) VALUES ('Alice')",
		},
		{
			name: "CREATE TABLE",
			sql:  "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)",
		},
		{
			name: "DROP TABLE",
			sql:  "DROP TABLE users",
		},
		{
			name: "UPDATE with WHERE",
			sql:  "UPDATE users SET name = 'Bob' WHERE id = 1",
		},
		{
			name: "DELETE with WHERE",
			sql:  "DELETE FROM users WHERE id = 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSafety(tt.sql)
			if err != nil {
				t.Errorf("CheckSafety() unexpected error = %v", err)
			}
		})
	}
}

func TestCheckSafety_CommentsAndWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "UPDATE with WHERE after comment",
			sql:     "UPDATE users SET name = 'Alice' -- comment\nWHERE id = 1",
			wantErr: false,
		},
		{
			name:    "UPDATE without WHERE with comment",
			sql:     "UPDATE users SET name = 'Alice' -- WHERE id = 1",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "UPDATE with WHERE in multi-line comment",
			sql:     "UPDATE users SET name = 'Alice' /* WHERE id = 1 */",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "DELETE with WHERE after multi-line comment",
			sql:     "DELETE FROM users /* some comment */ WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "UPDATE with WHERE and lots of whitespace",
			sql:     "UPDATE users SET name = 'Alice'   \n\n\n   WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "DELETE with WHERE 1=1 and comment",
			sql:     "DELETE FROM users WHERE 1=1 -- delete all",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSafety(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, vibeErr.Code)
				}
			}
		})
	}
}

func TestCheckSafety_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "WHERE in column name without WHERE clause (UPDATE)",
			sql:     "UPDATE users SET somewhere = 'location'",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "WHERE in string literal without WHERE clause (UPDATE)",
			sql:     "UPDATE users SET description = 'WHERE is my data'",
			wantErr: true,
			errCode: postgres.ErrorCodeUnsafeQuery,
		},
		{
			name:    "WHERE in column name WITH actual WHERE clause (UPDATE)",
			sql:     "UPDATE users SET somewhere = 'location' WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "WHERE as table alias",
			sql:     "DELETE FROM users WHERE id IN (SELECT id FROM orders WHERE status = 'deleted')",
			wantErr: false,
		},
		{
			name:    "Multiple WHERE clauses",
			sql:     "UPDATE users SET name = 'Alice' WHERE id = 1 AND status = 'active' WHERE enabled = true",
			wantErr: false,
		},
		{
			name:    "WHERE at start of line",
			sql:     "UPDATE users SET name = 'Alice'\nWHERE id = 1",
			wantErr: false,
		},
		{
			name:    "leading and trailing spaces",
			sql:     "   UPDATE users SET name = 'Alice' WHERE id = 1   ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSafety(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSafety() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != tt.errCode {
					t.Errorf("Expected error code %s, got %s", tt.errCode, vibeErr.Code)
				}
			}
		})
	}
}

func TestHasWhereClause(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want bool
	}{
		{
			name: "simple WHERE",
			sql:  "SELECT * FROM users WHERE id = 1",
			want: true,
		},
		{
			name: "no WHERE",
			sql:  "SELECT * FROM users",
			want: false,
		},
		{
			name: "WHERE in comment",
			sql:  "SELECT * FROM users -- WHERE id = 1",
			want: false,
		},
		{
			name: "WHERE in multi-line comment",
			sql:  "SELECT * FROM users /* WHERE id = 1 */",
			want: false,
		},
		{
			name: "WHERE after comment",
			sql:  "SELECT * FROM users -- comment\nWHERE id = 1",
			want: true,
		},
		{
			name: "SOMEWHERE (not a WHERE clause)",
			sql:  "UPDATE users SET somewhere = 'value'",
			want: false,
		},
		{
			name: "WHERE as word boundary",
			sql:  "UPDATE users SET name = 'Alice' WHERE id = 1",
			want: true,
		},
		{
			name: "lowercase where",
			sql:  "select * from users where id = 1",
			want: true,
		},
		{
			name: "mixed case WhErE",
			sql:  "SELECT * FROM users WhErE id = 1",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasWhereClause(tt.sql)
			if got != tt.want {
				t.Errorf("hasWhereClause() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveComments(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "single-line comment",
			sql:  "SELECT * FROM users -- this is a comment",
			want: "SELECT * FROM users ",
		},
		{
			name: "multi-line comment",
			sql:  "SELECT * FROM users /* this is a comment */",
			want: "SELECT * FROM users ",
		},
		{
			name: "multiple single-line comments",
			sql:  "SELECT * FROM users -- comment 1\nWHERE id = 1 -- comment 2",
			want: "SELECT * FROM users \nWHERE id = 1 ",
		},
		{
			name: "nested comments",
			sql:  "SELECT * FROM users /* outer /* inner */ comment */",
			want: "SELECT * FROM users  comment */",
		},
		{
			name: "no comments",
			sql:  "SELECT * FROM users WHERE id = 1",
			want: "SELECT * FROM users WHERE id = 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeComments(tt.sql)
			if got != tt.want {
				t.Errorf("removeComments() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveStringLiterals(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "simple string literal",
			sql:  "SELECT * FROM users WHERE name = 'Alice'",
			want: "SELECT * FROM users WHERE name = ''",
		},
		{
			name: "multiple string literals",
			sql:  "SELECT * FROM users WHERE name = 'Alice' AND city = 'NYC'",
			want: "SELECT * FROM users WHERE name = '' AND city = ''",
		},
		{
			name: "WHERE inside string literal",
			sql:  "UPDATE users SET description = 'WHERE is my data'",
			want: "UPDATE users SET description = ''",
		},
		{
			name: "escaped single quote",
			sql:  "INSERT INTO users (name) VALUES ('can''t')",
			want: "INSERT INTO users (name) VALUES ('')",
		},
		{
			name: "no string literals",
			sql:  "SELECT * FROM users WHERE id = 123",
			want: "SELECT * FROM users WHERE id = 123",
		},
		{
			name: "empty string",
			sql:  "SELECT * FROM users WHERE name = ''",
			want: "SELECT * FROM users WHERE name = ''",
		},
		{
			name: "string with SQL keywords",
			sql:  "UPDATE users SET note = 'DELETE FROM WHERE UPDATE'",
			want: "UPDATE users SET note = ''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeStringLiterals(tt.sql)
			if got != tt.want {
				t.Errorf("removeStringLiterals() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkCheckSafety_UPDATE_WithWhere(b *testing.B) {
	sql := "UPDATE users SET name = 'Alice' WHERE id = 1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckSafety(sql)
	}
}

func BenchmarkCheckSafety_DELETE_WithWhere(b *testing.B) {
	sql := "DELETE FROM users WHERE id = 1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckSafety(sql)
	}
}

func BenchmarkHasWhereClause(b *testing.B) {
	sql := "SELECT * FROM users WHERE id = 1 AND name = 'test'"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hasWhereClause(sql)
	}
}

func BenchmarkRemoveComments(b *testing.B) {
	sql := "SELECT * FROM users -- comment 1\nWHERE id = 1 /* comment 2 */"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = removeComments(sql)
	}
}
