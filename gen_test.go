package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	type args struct {
		yamlPath string
	}
	tests := []struct {
		name           string
		args           args
		wantOutputPath string
		wantErr        bool
	}{
		{
			args: args{
				yamlPath: "test/desc.yaml",
			},
			wantOutputPath: "test/output",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yamlBytes, err := os.ReadFile(tt.args.yamlPath)
			require.NoError(t, err)

			wantOutputBytes, err := os.ReadFile(tt.wantOutputPath)
			require.NoError(t, err)

			gotOutputBytes, err := generate(yamlBytes)
			require.NoError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, string(wantOutputBytes), string(gotOutputBytes))
		})
	}
}
