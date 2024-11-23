package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testfiles = [][]string{
	{"test/desc.yaml", "test/output"},
}

func TestGenerate(t *testing.T) {
	for i, tt := range testfiles {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			inputPath := tt[0]
			wantOutputPath := tt[1]

			yamlBytes, err := os.ReadFile(inputPath)
			require.NoError(t, err)

			wantOutputBytes, err := os.ReadFile(wantOutputPath)
			require.NoError(t, err)

			gotOutputBytes, err := Generate(yamlBytes)
			require.NoError(t, err)

			fmt.Println(string(gotOutputBytes))

			require.Equal(t, string(wantOutputBytes), string(gotOutputBytes))
		})
	}
}

func TestUpdateTestfile(t *testing.T) {
	for i, tt := range testfiles {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			inputPath := tt[0]

			yamlBytes, err := os.ReadFile(inputPath)
			require.NoError(t, err)

			gotOutputBytes, err := Generate(yamlBytes)
			require.NoError(t, err)

			err = os.WriteFile(tt[1], gotOutputBytes, 0644)
			require.NoError(t, err)
		})
	}
}
