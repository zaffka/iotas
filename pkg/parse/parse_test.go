package parse_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
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
		name         string
		typeNames    []string
		logErrMsg    string
		resConstants map[string][]string
	}

	tests := []testCase{
		{
			name:         "block started from non-zero value",
			typeNames:    []string{"TestType"},
			logErrMsg:    "{\"level\":\"error\",\"error\":\"first spec has no zero-iota value\",\"type_name\":\"TestType\",\"message\":\"parsing interrupted\"}\n",
			resConstants: map[string][]string{"TestType": nil},
		},
		{
			name:         "two const block with same type",
			typeNames:    []string{"TestType2"},
			logErrMsg:    "{\"level\":\"warn\",\"type_name\":\"TestType2\",\"message\":\"duplicated iota sequence found and skipped\"}\n",
			resConstants: map[string][]string{"TestType2": {"TestType2X"}},
		},
		{
			name:         "handling stopped at const alias",
			typeNames:    []string{"TestType3"},
			logErrMsg:    "",
			resConstants: map[string][]string{"TestType3": {"TestType3X"}},
		},
		{
			name:         "handling stopped at second iota-declaration",
			typeNames:    []string{"TestType4"},
			logErrMsg:    "",
			resConstants: map[string][]string{"TestType4": {"TestType4X"}},
		},
		{
			name:         "block started from untyped const",
			typeNames:    []string{"TestType5"},
			logErrMsg:    "",
			resConstants: map[string][]string{"TestType5": nil},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}

			p, err := parse.NewParser(parse.Deps{
				Dir:       testPkgDir,
				TypeNames: tt.typeNames,
				Logger:    zerolog.New(buf),
			})
			require.NoError(t, err)

			p.Parse()

			require.Equal(t, tt.logErrMsg, buf.String())
			require.Equal(t, tt.resConstants, p.GetConstantsByType())
		})
	}
}
