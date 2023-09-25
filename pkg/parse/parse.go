package parse

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

const (
	iotaName = "iota"
)

// Parse inspects nodes of the ast-tree, gets constant blocks,
// finds necessary types and return errors if gets any in process.
func (p *Parser) Parse() error {
	for _, file := range p.astFiles {
		ast.Inspect(file, p.findConstantDecl)
	}

	return errors.Join(p.errs...)
}

func (p *Parser) findConstantDecl(astNode ast.Node) bool {
	declaration, ok := astNode.(*ast.GenDecl)
	if !ok || declaration.Tok != token.CONST {
		return true
	}

	// finding a constant sequence and its type
	sequenceTypeName, constants, err := getTypeAndConstantNames(declaration.Specs)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("constant parsing finished with an error: %w", err))

		return true
	}

	savedFirst, ok := p.ConstantsByType[sequenceTypeName]
	// if second block of constants with the same type name found
	// then we delete the entire type keeping an error
	if ok && savedFirst != nil {
		delete(p.ConstantsByType, sequenceTypeName)

		p.errs = append(p.errs, Errors{iotaDuplicatedSequence, sequenceTypeName})

		return true
	}

	if ok {
		p.ConstantsByType[sequenceTypeName] = constants
	}

	return false
}

// getTypeAndConstantNames returns constant slice with its type name.
func getTypeAndConstantNames(specs []ast.Spec) (string, []string, error) {
	if len(specs) == 0 {
		return "", nil, Errors{emptySpecs}
	}

	typeName, firstConstName, err := getFirst(specs[0])
	if err != nil {
		return "", nil, fmt.Errorf("failed to get a first iota const: %w", err)
	}

	res := []string{firstConstName}

	for _, spec := range specs[1:] {
		next, ok := getNext(spec)
		if !ok {
			break
		}

		res = append(res, next)
	}

	return typeName, res, nil
}

// getFirst gets first constant in a block with name and its type name.
func getFirst(spec ast.Spec) (string, string, error) {
	val, ok := spec.(*ast.ValueSpec)
	if !ok {
		return "", "", Errors{isNotAValueSpec, fmt.Sprintf("%t", spec)}
	}

	if val.Type == nil {
		return "", "", Errors{isUntypedValueSpec}
	}

	if len(val.Values) == 0 {
		return "", "", Errors{noValuesAtValueSpec}
	}

	valType, ok := val.Type.(*ast.Ident)
	if !ok {
		return "", "", Errors{isNotAnIdentNode, fmt.Sprintf("%t", val.Type)}
	}

	valValue, ok := val.Values[0].(*ast.Ident)
	if !ok || valValue.Name != iotaName {
		return "", "", Errors{iotaIdentExpected}
	}

	constName, err := getFirstName(val.Names)
	if err != nil {
		return "", "", err
	}

	return valType.Name, constName, nil
}

// getNext gets a constant name from the spec.
func getNext(spec ast.Spec) (string, bool) {
	val, ok := spec.(*ast.ValueSpec)
	if !ok {
		return "", false
	}

	// all specs after the first one must have nil Type and Value fields
	if val.Type != nil || val.Values != nil {
		return "", false
	}

	constName, err := getFirstName(val.Names)
	if err != nil {
		return "", false
	}

	return constName, true
}

// getFirstName gets a name from a slice of idents.
func getFirstName(idents []*ast.Ident) (string, error) {
	if len(idents) == 0 {
		return "", Errors{emptyIdentList}
	}

	return idents[0].Name, nil
}
