package query

import (
	"regexp"
	"strings"

	"github.com/vibesql/vibe/internal/postgres"
)

var (
	whereClausePattern = regexp.MustCompile(`\bWHERE\b`)
	singleLineComment  = regexp.MustCompile(`--[^\n]*`)
	multiLineComment   = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	stringLiteral      = regexp.MustCompile(`'(?:[^']|'')*'`)
)

// CheckSafety enforces safety rules on SQL queries
func CheckSafety(sql string) error {
	trimmed := strings.TrimSpace(sql)
	upperSQL := strings.ToUpper(trimmed)

	// Check UPDATE without WHERE
	if strings.HasPrefix(upperSQL, "UPDATE") {
		if !hasWhereClause(trimmed) {
			return postgres.NewVibeError(
				postgres.ErrorCodeUnsafeQuery,
				"Unsafe query: UPDATE without WHERE clause",
				"UPDATE queries must include a WHERE clause. Use 'WHERE 1=1' to update all rows explicitly",
			)
		}
	}

	// Check DELETE without WHERE
	if strings.HasPrefix(upperSQL, "DELETE") {
		if !hasWhereClause(trimmed) {
			return postgres.NewVibeError(
				postgres.ErrorCodeUnsafeQuery,
				"Unsafe query: DELETE without WHERE clause",
				"DELETE queries must include a WHERE clause. Use 'WHERE 1=1' to delete all rows explicitly",
			)
		}
	}

	return nil
}

// hasWhereClause checks if a SQL query contains a WHERE clause
// It removes comments and string literals to avoid false positives
func hasWhereClause(sql string) bool {
	// Remove SQL comments before checking
	sql = removeComments(sql)
	
	// Remove string literals to avoid false positives
	// e.g., UPDATE users SET desc = 'WHERE is my data' should not match
	sql = removeStringLiterals(sql)
	
	// Convert to uppercase for case-insensitive matching
	upperSQL := strings.ToUpper(sql)
	
	// Check for WHERE keyword
	// Using word boundary matching to avoid false positives like "SOMEWHERE"
	return whereClausePattern.MatchString(upperSQL)
}

// removeComments removes SQL comments from the query
// Note: Nested /* */ comments are not fully supported
// (matches PostgreSQL default behavior)
func removeComments(sql string) string {
	// Remove single-line comments (-- comment)
	sql = singleLineComment.ReplaceAllString(sql, "")
	
	// Remove multi-line comments (/* comment */)
	sql = multiLineComment.ReplaceAllString(sql, "")
	
	return sql
}

// removeStringLiterals removes SQL string literals from the query
// This prevents false positives when WHERE appears inside strings
// Handles PostgreSQL string escaping: 'can''t' (doubled single quotes)
func removeStringLiterals(sql string) string {
	return stringLiteral.ReplaceAllString(sql, "''")
}
