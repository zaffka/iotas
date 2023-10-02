package parse

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/rs/zerolog"
	"golang.org/x/tools/go/packages"
)

var (
	ErrMultiplePackages    = errors.New("multiple packages within a folder")
	ErrDuplicatedTypeParam = errors.New("duplicated type param received")
	ErrEmptyTypeName       = errors.New("empty type name param received")
)

// Deps struct is a set of required parser's dependencies.
type Deps struct {
	Dir       string
	TypeNames []string
	Logger    zerolog.Logger
}

// NewParser creates and initializes a Parser for Go package.
// It loads a package with all types and syntax info, gets package name, etc.
func NewParser(deps Deps) (*Parser, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
	}, deps.Dir)
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

	tnMap := make(constantsByType, len(deps.TypeNames))

	for _, tn := range deps.TypeNames {
		if tn == "" {
			return nil, ErrEmptyTypeName
		}

		_, exist := tnMap[tn]
		if exist {
			return nil, fmt.Errorf("%w: %s", ErrDuplicatedTypeParam, tn)
		}

		tnMap[tn] = nil
	}

	return &Parser{
		pkg: struct {
			name     string
			astFiles []*ast.File
		}{name: pkg.Name, astFiles: pkg.Syntax},
		log:             deps.Logger,
		constantsByType: tnMap,
	}, nil
}

// Parser holds a raw ast-tree data with pre-initialized map for constant names grouped by type name.
type Parser struct {
	pkg struct {
		name     string
		astFiles []*ast.File
	}

	log zerolog.Logger

	constantsByType constantsByType
}

func (p Parser) GetPackageName() string {
	return p.pkg.name
}

func (p Parser) GetConstantsByType() map[string][]string {
	return p.constantsByType
}

type constantsByType map[string][]string
