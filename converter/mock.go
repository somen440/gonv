package converter

import "github.com/somen440/gonv/structure"

// CreateMockDatabaseStructure db mock from build/schema/*.sql
func CreateMockDatabaseStructure() *structure.DatabaseStructure {
	return &structure.DatabaseStructure{
		Map: map[structure.TableName]*structure.TableStructure{
			structure.TableName("sample"):     CreateMockSampleTableStructure(),
			structure.TableName("sample_log"): CreateMockSampleLogTableStructure(),
		},
	}
}

// CreateMockSampleTableStructure mock sample table
func CreateMockSampleTableStructure() *structure.TableStructure {
	return &structure.TableStructure{
		Table:          "sample",
		Type:           structure.TableType,
		Comment:        "sample table",
		Engine:         "InnoDB",
		DefaultCharset: "utf8mb4",
		Collate:        "utf8mb4_unicode_ci",
		ColumnStructureList: map[structure.ColumnField]*structure.MySQL57ColumnStructure{
			structure.ColumnField("id"): {
				Field:   structure.ColumnField("id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "Sample ID",
				Attributes: []structure.Attribute{
					structure.AutoIncrement,
					structure.Unsigned,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                0,
			},
			structure.ColumnField("name"): {
				Field:                structure.ColumnField("name"),
				Type:                 "varchar(255)",
				Default:              "sample",
				Comment:              "Sample Name",
				Attributes:           []structure.Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                1,
			},
			structure.ColumnField("created"): {
				Field:                structure.ColumnField("created"),
				Type:                 "datetime",
				Default:              "CURRENT_TIMESTAMP",
				Comment:              "Created Time",
				Attributes:           []structure.Attribute{},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                2,
			},
			structure.ColumnField("modified"): {
				Field:   structure.ColumnField("modified"),
				Type:    "datetime",
				Default: "CURRENT_TIMESTAMP",
				Comment: "Modified Time",
				Attributes: []structure.Attribute{
					structure.OnUpdateCurrentTimestamp,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                3,
			},
		},
		IndexStructureList: map[structure.IndexKey]*structure.IndexStructure{
			structure.IndexKey("PRIMARY"): structure.NewIndexStructure("PRIMARY", "BTREE", true, []string{"id"}, 0),
		},
		Partition:  nil,
		Properties: []string{},
	}
}

// CreateMockSampleLogTableStructure mock table sample log
func CreateMockSampleLogTableStructure() *structure.TableStructure {
	return &structure.TableStructure{
		Table:          "sample_log",
		Type:           structure.TableType,
		Comment:        "sample log table",
		Engine:         "InnoDB",
		DefaultCharset: "utf8mb4",
		Collate:        "utf8mb4_unicode_ci",
		ColumnStructureList: map[structure.ColumnField]*structure.MySQL57ColumnStructure{
			structure.ColumnField("id"): {
				Field:   structure.ColumnField("id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "",
				Attributes: []structure.Attribute{
					structure.AutoIncrement,
					structure.Unsigned,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                0,
			},
			structure.ColumnField("month"): {
				Field:   structure.ColumnField("month"),
				Type:    "tinyint(2)",
				Default: "",
				Comment: "",
				Attributes: []structure.Attribute{
					structure.Unsigned,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                1,
			},
			structure.ColumnField("sample_id"): {
				Field:   structure.ColumnField("sample_id"),
				Type:    "bigint(20)",
				Default: "",
				Comment: "",
				Attributes: []structure.Attribute{
					structure.Unsigned,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                2,
			},
			structure.ColumnField("name"): {
				Field:                structure.ColumnField("name"),
				Type:                 "varchar(255)",
				Default:              "",
				Comment:              "",
				Attributes:           []structure.Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                3,
			},
			structure.ColumnField("created"): {
				Field:                structure.ColumnField("created"),
				Type:                 "datetime",
				Default:              "CURRENT_TIMESTAMP",
				Comment:              "",
				Attributes:           []structure.Attribute{},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                4,
			},
			structure.ColumnField("modified"): {
				Field:   structure.ColumnField("modified"),
				Type:    "datetime",
				Default: "CURRENT_TIMESTAMP",
				Comment: "",
				Attributes: []structure.Attribute{
					structure.OnUpdateCurrentTimestamp,
				},
				CollationName:        "",
				Properties:           []string{},
				GenerationExpression: "",
				Order:                5,
			},
		},
		IndexStructureList: map[structure.IndexKey]*structure.IndexStructure{
			structure.IndexKey("PRIMARY"):   structure.NewIndexStructure("PRIMARY", "BTREE", true, []string{"id", "month"}, 0),
			structure.IndexKey("sample_id"): structure.NewIndexStructure("sample_id", "BTREE", false, []string{"sample_id"}, 1),
		},
		Partition: &structure.PartitionLongStructure{
			Type:  "LIST",
			Value: "month",
			Parts: []*structure.PartitionPartStructure{
				{
					Name:     "p1",
					Operator: "IN",
					Value:    "1",
					Comment:  "",
				},
				{
					Name:     "p2",
					Operator: "IN",
					Value:    "2",
					Comment:  "",
				},
				{
					Name:     "p3",
					Operator: "IN",
					Value:    "3",
					Comment:  "",
				},
				{
					Name:     "p4",
					Operator: "IN",
					Value:    "4",
					Comment:  "",
				},
				{
					Name:     "p5",
					Operator: "IN",
					Value:    "5",
					Comment:  "",
				},
				{
					Name:     "p6",
					Operator: "IN",
					Value:    "6",
					Comment:  "",
				},
				{
					Name:     "p7",
					Operator: "IN",
					Value:    "7",
					Comment:  "",
				},
				{
					Name:     "p8",
					Operator: "IN",
					Value:    "8",
					Comment:  "",
				},
				{
					Name:     "p9",
					Operator: "IN",
					Value:    "9",
					Comment:  "",
				},
				{
					Name:     "p10",
					Operator: "IN",
					Value:    "10",
					Comment:  "",
				},
				{
					Name:     "p11",
					Operator: "IN",
					Value:    "11",
					Comment:  "",
				},
				{
					Name:     "p12",
					Operator: "IN",
					Value:    "12",
					Comment:  "",
				},
			},
		},
		Properties: []string{},
	}
}
