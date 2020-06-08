package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCreateQuery(t *testing.T) {
	tests := []struct {
		column   *MySQL57ColumnStructure
		expected string
	}{
		{
			&MySQL57ColumnStructure{
				Field:   ColumnField("id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "Sample ID",
				Attributes: []Attribute{
					AutoIncrement,
					Unsigned,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
			},
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID'",
		},
	}
	for _, tt := range tests {
		actual := tt.column.GenerateCreateQuery()
		assert.Equal(t, tt.expected, actual)
	}
}
