package parse

import (
	"errors"
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

var (
	ErrMultiplePackages    = errors.New("multiple packages within a folder")
	ErrDuplicatedTypeParam = errors.New("duplicated type param received")
)

// NewParser creates and initializes a Parser for Go package.
// It loads a package with all types and syntax info, gets package name, etc.
func NewParser(dir string, typeNames []string) (*Parser, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
	}, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load a package: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, ErrMultiplePackages
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("package loading finished with an error: %w", pkg.Errors[0])
	}

	tnMap := make(map[string][]string, len(typeNames))
	for _, tn := range typeNames {
		_, exist := tnMap[tn]
		if exist {
			return nil, fmt.Errorf("%w: %s", ErrDuplicatedTypeParam, tn)
		}

		tnMap[tn] = nil
	}

	return &Parser{
		PkgName:         pkg.Name,
		ConstantsByType: tnMap,

		astFiles: pkg.Syntax,
	}, nil
}

// Parser holds a raw ast-tree data with pre-initialized map for constant names grouped by type name.
type Parser struct {
	PkgName         string
	astFiles        []*ast.File
	ConstantsByType map[string][]string

	errs []error
}
