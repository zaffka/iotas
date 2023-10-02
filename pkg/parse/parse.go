package parse

import (
	"errors"
	"go/ast"
	"go/token"
)

const (
	iotaName = "iota"
)

type typeCons struct {
	typeName string
	conNames []string
}

type err string

func (e err) Error() string {
	return string(e)
}

// Parse inspects nodes of the ast-tree, gets constant blocks,
// finds necessary types and return errors if gets any in process.
func (p *Parser) Parse() {
	for _, file := range p.pkg.astFiles {
		ast.Inspect(file, p.findConstantDecl)
	}
}

func (p *Parser) findConstantDecl(astNode ast.Node) bool {
	declaration, ok := astNode.(*ast.GenDecl)
	if !ok || declaration.Tok != token.CONST {
		return true
	}

	// finding a constant sequence and its type
	tc, err := getTypeAndConstantNames(declaration.Specs)
	if errors.Is(err, errSkipToken) {
		return true
	}

	// skip unnecessary
	savedFirst, ok := p.constantsByType[tc.typeName]
	if !ok {
		return true
	}

	// handling an error if type is required
	if err != nil {
		p.log.Err(err).Str("type_name", tc.typeName).Msg("parsing interrupted")

		return true
	}

	// if second block of constants with the same type name found
	// then we delete the entire type keeping an error
	if savedFirst != nil {
		p.log.Warn().Str("type_name", tc.typeName).Msg(iotaDuplicatedSequenceSkipped)

		return true
	}

	// if we need such a type and there were no const block previously found
	// keeping the constants
	p.constantsByType[tc.typeName] = tc.conNames

	return false
}

// getTypeAndConstantNames returns constant slice with its type name.
func getTypeAndConstantNames(specs []ast.Spec) (*typeCons, error) {
	if len(specs) == 0 {
		return nil, err(emptySpecs)
	}

	tc, err := getFirstConst(specs[0])
	if err != nil {
		return tc, err
	}

	for _, spec := range specs[1:] {
		next, ok := getNext(spec)
		if !ok {
			break
		}

		tc.conNames = append(tc.conNames, next)
	}

	return tc, nil
}

// getFirstConst gets first constant in a block with name and its type name.
func getFirstConst(spec ast.Spec) (*typeCons, error) {
	val, ok := spec.(*ast.ValueSpec)
	if !ok {
		return nil, errSkipToken
	}

	if val.Type == nil {
		return nil, errSkipToken
	}

	if len(val.Values) == 0 {
		return nil, errSkipToken
	}

	valType, ok := val.Type.(*ast.Ident)
	if !ok {
		return nil, errSkipToken
	}

	fc := &typeCons{
		typeName: valType.Name,
	}

	valValue, ok := val.Values[0].(*ast.Ident)
	if !ok || valValue.Name != iotaName {
		return fc, err(iotaIdentExpected)
	}

	constName, err := getFirstName(val.Names)
	if err != nil {
		return fc, err
	}

	fc.conNames = []string{constName}

	return fc, nil
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
		return "", err(emptyIdentList)
	}

	return idents[0].Name, nil
}
