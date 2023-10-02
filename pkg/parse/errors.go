package parse

import (
	"errors"
)

const (
	// Error messages for parse.Errors type.
	iotaDuplicatedSequenceSkipped = "duplicated iota sequence found and skipped"
	iotaIdentExpected             = "first spec has no zero-iota value"

	emptySpecs     = "specs list is empty"
	emptyIdentList = "ident list is empty"
)

var errSkipToken = errors.New("skipping needless token")
