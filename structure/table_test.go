package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t1 := &TableStructure{
		Table:          "sample",
		Type:           TableType,
		Comment:        "sample table",
		Engine:         "InnoDB",
		DefaultCharset: "utf8mb4",
		Collate:        "utf8mb4_unicode_ci",
		Properties:     []string{},
		ColumnStructureList: []*MySQL57ColumnStructure{
			{
				Field:   ColumnField("id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "Sample ID",
				Attributes: []Attribute{
					AutoIncrement,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
			},
		},
		IndexStructureList: []*IndexStructure{},
		Partition:          nil,
	}
	actual := t1.String()

	expected := "name: sample\n"
	expected += "type: table\n"
	expected += "comment: sample table\n"
	expected += "engine: InnoDB\n"
	expected += "default_charset: utf8mb4\n"
	expected += "collate: utf8mb4_unicode_ci\n"
	expected += "properties: \n"
	expected += "columns:\n"
	expected += "\t`id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Sample ID'\n"

	assert.Equal(t, expected, actual)
}
