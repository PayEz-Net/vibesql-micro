package query

import (
	"strings"
	"testing"

	"github.com/vibesql/vibe/internal/postgres"
)

func TestValidateQuery_EmptySQL(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "empty string",
			sql:     "",
			wantErr: true,
			errCode: postgres.ErrorCodeMissingRequiredField,
		},
		{
			name:    "whitespace only",
			sql:     "   \t\n  ",
			wantErr: true,
			errCode: postgres.ErrorCodeMissingRequiredField,
		},
		{
			name:    "single space",
			sql:     " ",
			wantErr: true,
			errCode: postgres.ErrorCodeMissingRequiredField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
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

func TestValidateQuery_QueryTooLarge(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "exactly 10KB",
			sql:     "SELECT " + strings.Repeat("a", MaxQuerySize-7), // "SELECT " is 7 bytes
			wantErr: false,
		},
		{
			name:    "one byte over 10KB",
			sql:     "SELECT " + strings.Repeat("a", MaxQuerySize-6), // One byte over
			wantErr: true,
			errCode: postgres.ErrorCodeQueryTooLarge,
		},
		{
			name:    "significantly over 10KB",
			sql:     "SELECT " + strings.Repeat("a", MaxQuerySize*2),
			wantErr: true,
			errCode: postgres.ErrorCodeQueryTooLarge,
		},
		{
			name:    "small query",
			sql:     "SELECT 1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
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

func TestValidateQuery_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "random text",
			sql:     "this is not sql",
			wantErr: true,
			errCode: postgres.ErrorCodeInvalidSQL,
		},
		{
			name:    "gibberish",
			sql:     "asdfghjkl",
			wantErr: true,
			errCode: postgres.ErrorCodeInvalidSQL,
		},
		{
			name:    "number only",
			sql:     "12345",
			wantErr: true,
			errCode: postgres.ErrorCodeInvalidSQL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
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

func TestValidateQuery_ValidQueries(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "simple SELECT",
			sql:  "SELECT 1",
		},
		{
			name: "SELECT with WHERE",
			sql:  "SELECT * FROM users WHERE id = 1",
		},
		{
			name: "INSERT",
			sql:  "INSERT INTO users (name) VALUES ('Alice')",
		},
		{
			name: "UPDATE with WHERE",
			sql:  "UPDATE users SET name = 'Bob' WHERE id = 1",
		},
		{
			name: "DELETE with WHERE",
			sql:  "DELETE FROM users WHERE id = 1",
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
			name: "lowercase SELECT",
			sql:  "select * from users",
		},
		{
			name: "mixed case",
			sql:  "SeLeCt * FrOm users",
		},
		{
			name: "leading whitespace",
			sql:  "   SELECT 1",
		},
		{
			name: "trailing whitespace",
			sql:  "SELECT 1   ",
		},
		{
			name: "newlines",
			sql:  "SELECT *\nFROM users\nWHERE id = 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.sql)
			if err != nil {
				t.Errorf("ValidateQuery() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateQuery_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
		errCode string
	}{
		{
			name:    "ALTER (valid keyword)",
			sql:     "ALTER TABLE users ADD COLUMN email TEXT",
			wantErr: false,
		},
		{
			name:    "TRUNCATE (valid keyword)",
			sql:     "TRUNCATE TABLE users",
			wantErr: false,
		},
		{
			name:    "SELECT in middle of junk",
			sql:     "junk SELECT 1",
			wantErr: true,
			errCode: postgres.ErrorCodeInvalidSQL,
		},
		{
			name:    "tabs and spaces",
			sql:     "\t\t  SELECT 1  \t\t",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
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

func BenchmarkValidateQuery_Simple(b *testing.B) {
	sql := "SELECT * FROM users WHERE id = 1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateQuery(sql)
	}
}

func BenchmarkValidateQuery_Large(b *testing.B) {
	sql := "SELECT " + strings.Repeat("column_name, ", 100) + "1 FROM users"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateQuery(sql)
	}
}
