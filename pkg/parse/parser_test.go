package parse_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zaffka/iotas/pkg/parse"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	type args struct {
		dir       string
		typeNames []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantRes bool
	}{
		{
			name: "no package",
			args: args{
				dir:       "..",
				typeNames: []string{"MatrixType"},
			},
			wantErr: true,
			wantRes: false,
		},
		{
			name: "ok",
			args: args{
				dir:       "../../examples",
				typeNames: []string{"MatrixType"},
			},
			wantErr: false,
			wantRes: true,
		},
		{
			name: "duplicated type name",
			args: args{
				dir:       "../../examples",
				typeNames: []string{"MatrixType", "MatrixType"},
			},
			wantErr: true,
			wantRes: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parse.NewParser(tt.args.dir, tt.args.typeNames)
			if tt.wantErr {
				require.NotNil(t, err)
			}

			if tt.wantRes {
				require.NotNil(t, got)
			}
		})
	}
}
