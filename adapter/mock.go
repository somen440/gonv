package adapter

import "github.com/somen440/gonv/structure"

func createSampleTableSt() *structure.TableStructure {
	return &structure.TableStructure{
		Table:          "sample",
		Type:           structure.TableType,
		Comment:        "sample table",
		Engine:         "InnoDB",
		DefaultCharset: "utf8mb4",
		Collate:        "utf8mb4_unicode_ci",
		Properties:     []string{},
		ColumnStructureList: []*structure.MySQL57ColumnStructure{
			{
				Field:   structure.ColumnField("id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "Sample ID",
				Attributes: []structure.Attribute{
					structure.AutoIncrement,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
			},
			{
				Field:                structure.ColumnField("name"),
				Type:                 "varchar(255)",
				Default:              "",
				Comment:              "Sample Name",
				Attributes:           []structure.Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
			},
			{
				Field:                structure.ColumnField("created"),
				Type:                 "datetime",
				Default:              "CURRENT_TIMESTAMP",
				Comment:              "Sample Name",
				Attributes:           []structure.Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
			},
			{
				Field:                structure.ColumnField("modified"),
				Type:                 "datetime",
				Default:              "CURRENT_TIMESTAMP",
				Comment:              "Sample Name",
				Attributes:           []structure.Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
			},
		},
		IndexStructureList: []*structure.IndexStructure{},
		Partition:          nil,
	}
}
