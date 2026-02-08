package query

import (
	"strings"

	"github.com/vibesql/vibe/internal/postgres"
)

const (
	// MaxQuerySize is the maximum allowed SQL query length (10KB)
	MaxQuerySize = 10 * 1024 // 10KB in bytes
)

// ValidateQuery validates a SQL query for basic requirements
func ValidateQuery(sql string) error {
	// Check for empty SQL
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return postgres.NewVibeError(
			postgres.ErrorCodeMissingRequiredField,
			"Missing required field",
			"The 'sql' field is required and cannot be empty",
		)
	}

	// Check query length (10KB limit)
	if len(sql) > MaxQuerySize {
		return postgres.NewVibeError(
			postgres.ErrorCodeQueryTooLarge,
			"Query too large",
			"SQL query exceeds the maximum allowed size of 10KB",
		)
	}

	// Basic SQL syntax validation - check for at least one SQL keyword
	// Detailed syntax validation is deferred to PostgreSQL engine,
	// which returns SQLSTATE codes that we map to INVALID_SQL errors
	upperSQL := strings.ToUpper(trimmed)
	validKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE",
		"CREATE", "DROP", "ALTER", "TRUNCATE",
	}

	hasValidKeyword := false
	for _, keyword := range validKeywords {
		if strings.HasPrefix(upperSQL, keyword) {
			hasValidKeyword = true
			break
		}
	}

	if !hasValidKeyword {
		return postgres.NewVibeError(
			postgres.ErrorCodeInvalidSQL,
			"Invalid SQL syntax",
			"Query must start with a valid SQL keyword (SELECT, INSERT, UPDATE, DELETE, CREATE, DROP)",
		)
	}

	return nil
}
