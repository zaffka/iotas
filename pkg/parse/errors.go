package parse

import (
	"errors"
	"strings"
)

// Errors holds a set of strings with errors found while parsing the ast-tree.
type Errors []string

// Error method to make Error satisfy the Error interface.
func (ee Errors) Error() string {
	return strings.Join(ee, ":")
}

const (
	// Error messages for parse.Errors type.
	iotaDuplicatedSequence = "duplicated iota sequence"
	iotaIdentExpected      = "first spec has no zeroed-iota value"

	emptySpecs     = "specs list is empty"
	emptyIdentList = "ident list is empty"

	isNotAValueSpec    = "value is not a ast.ValueSpec"
	isUntypedValueSpec = "value spec is untyped"
	isNotAnIdentNode   = "value is not an ast.Ident node"

	noValuesAtValueSpec = "no values in spec"
)

var errSkipToken = errors.New("skip needless token")
