package gonv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInSlice(t *testing.T) {
	tests := []struct {
		title    string
		expect   bool
		needle   string
		haystack []string
	}{
		{"success", true, "a", []string{"a", "b", "c"}},
		{"failed", false, "d", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			actual := InSlice(tt.needle, tt.haystack)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
