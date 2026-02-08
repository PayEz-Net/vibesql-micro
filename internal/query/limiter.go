package query

import (
	"github.com/vibesql/vibe/internal/postgres"
)

const (
	MaxResultRows = 1000
)

func CheckRowLimit(currentRowCount int) error {
	if currentRowCount >= MaxResultRows {
		return postgres.NewVibeError(
			postgres.ErrorCodeResultTooLarge,
			"Result set too large",
			"Query returned more than the maximum allowed 1000 rows",
		)
	}
	return nil
}
