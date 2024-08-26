package util

import (
	"gotest.tools/v3/assert"
	"testing"
)

func Test_merkleHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"expected result",
			"31323334353637",
			"8bb0cf6eb9b17d0f7d22b456f121257dc1254e1f01665370476383ea776df414",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MerkleHash(tt.input)
			assert.NilError(t, err)
			assert.Equal(t, tt.expected, got)

		})
	}
}
