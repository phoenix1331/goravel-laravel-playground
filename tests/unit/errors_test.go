// Laravel equivalent: expecting a specific exception class - see docs/chapters/24-error-handling.md
package unit

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"goravel/app/results"
)

// #region test
func TestWrappedErrorsKeepTheirIdentity(t *testing.T) {
	// What results.Record returns, wrapped once more by a caller
	err := fmt.Errorf("sync: %w", fmt.Errorf("record fixture 7: %w", results.ErrAlreadyRecorded))

	// errors.Is unwraps the chain looking for the sentinel, like
	// catch (AlreadyRecorded $e) matching a subclass
	assert.True(t, errors.Is(err, results.ErrAlreadyRecorded))

	// == only compares the outer value, so it misses the wrapped sentinel
	assert.False(t, err == results.ErrAlreadyRecorded)

	// Each layer added context to the message
	assert.Equal(t, "sync: record fixture 7: fixture already has a result", err.Error())
}

// #endregion test
