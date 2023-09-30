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

	tnMap := make(constantsByType, len(typeNames))
	for _, tn := range typeNames {
		_, exist := tnMap[tn]
		if exist {
			return nil, fmt.Errorf("%w: %s", ErrDuplicatedTypeParam, tn)
		}

		tnMap[tn] = nil
	}

	return &Parser{
		pkgName:         pkg.Name,
		constantsByType: tnMap,

		astFiles: pkg.Syntax,
	}, nil
}

// Parser holds a raw ast-tree data with pre-initialized map for constant names grouped by type name.
type Parser struct {
	pkgName         string
	astFiles        []*ast.File
	constantsByType constantsByType

	errs []error
}

func (p Parser) GetPackageName() string {
	return p.pkgName
}

func (p Parser) GetConstantsByType() map[string][]string {
	return p.constantsByType
}

type constantsByType map[string][]string
