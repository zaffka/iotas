package parse_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zaffka/iotas/pkg/parse"
)

func TestParser_Parse(t *testing.T) {
	testPkgDir, err := filepath.Abs("./testpkg")
	if err != nil {
		t.Fatal("failed to get testpkg path")
	}

	t.Parallel()

	type testCase struct {
		name      string
		typeNames []string

		err          error
		resConstants map[string][]string
	}

	tests := []testCase{
		// {
		// 	name:         "block started not from zero-iota",
		// 	typeNames:    []string{"TestType"},
		// 	errS:         "first spec has no zeroed-iota value",
		// 	resConstants: map[string][]string{"TestType": nil},
		// },
		// {
		// 	name:         "two const block with same type",
		// 	typeNames:    []string{"TestType2"},
		// 	errS:         "duplicated iota sequence",
		// 	resConstants: map[string][]string{},
		// },
		{
			name:         "handling stopped at const alias",
			typeNames:    []string{"TestType3"},
			err:          nil,
			resConstants: map[string][]string{"TestType3": nil},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p, err := parse.NewParser(testPkgDir, tt.typeNames)
			require.NoError(t, err)

			require.Equal(t, p.Parse(), tt.err)
			require.Equal(t, tt.resConstants, p.GetConstantsByType())
		})
	}
}
