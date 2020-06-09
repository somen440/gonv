package main

import (
	"testing"

	"github.com/somen440/gonv/structure"
	"github.com/stretchr/testify/assert"
)

func setUpTestFactory() (*Factory, func()) {
	config := &DBConfig{
		Driver:   MySQL,
		User:     "root",
		Password: "test",
		Database: "test",
		Host:     "localhost",
		Port:     "33066",
	}
	gdo := NewGDO(config)
	factory := &Factory{
		gdo: gdo,
	}
	return factory, gdo.Close
}

func TestCreateDatabaseStructure(t *testing.T) {
	factory, df := setUpTestFactory()
	defer df()

	expected := &structure.DatabaseStructure{
		Map: map[structure.TableName]*structure.TableStructure{
			structure.TableName("sample"): {
				Table:          "sample",
				Type:           structure.TableType,
				Comment:        "",
				Engine:         "InnoDB",
				DefaultCharset: "utf8mb4",
				Collate:        "utf8mb4_unicode_ci",
				ColumnStructureList: []*structure.MySQL57ColumnStructure{
					{
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
					},
					{
						Field:                structure.ColumnField("name"),
						Type:                 "varchar(255)",
						Default:              "sample",
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
						Comment:              "Created Time",
						Attributes:           []structure.Attribute{},
						CollationName:        "",
						Properties:           []string{},
						GenerationExpression: "",
					},
					{
						Field:                structure.ColumnField("modified"),
						Type:                 "datetime",
						Default:              "CURRENT_TIMESTAMP",
						Comment:              "Modified Time",
						Attributes:           []structure.Attribute{},
						CollationName:        "",
						Properties:           []string{},
						GenerationExpression: "",
					},
				},
				IndexStructureList: []*structure.IndexStructure{
					structure.NewIndexStructure("PRIMARY", "BTREE", true, []string{"id"}),
				},
				Partition:  nil,
				Properties: []string{},
			},
			structure.TableName("sample_log"): {
				Table:               "sample_log",
				Type:                structure.TableType,
				Comment:             "",
				Engine:              "InnoDB",
				DefaultCharset:      "utf8mb4",
				Collate:             "utf8mb4_unicode_ci",
				ColumnStructureList: []*structure.MySQL57ColumnStructure{},
				IndexStructureList:  []*structure.IndexStructure{},
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
			},
		},
	}

	actual, err := factory.CreateDatabaseStructure("test")
	assert.Nil(t, err)

	assertEqualTable(t, expected.Map[structure.TableName("sample")], actual.Map[structure.TableName("sample")])
	assertEqualTable(t, expected.Map[structure.TableName("sample_log")], actual.Map[structure.TableName("sample_log")])
}

func assertEqualTable(t *testing.T, expected, actual *structure.TableStructure) {
	assert.Equal(t, expected.Table, actual.Table)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Comment, actual.Comment)
	assert.Equal(t, expected.Engine, actual.Engine)
	assert.Equal(t, expected.Collate, actual.Collate)
	assert.Equal(t, expected.DefaultCharset, actual.DefaultCharset)
	assert.Equal(t, expected.Partition, actual.Partition)
	for i, expectedColumn := range expected.ColumnStructureList {
		actualColumn := actual.ColumnStructureList[i]
		assert.Equal(t, expectedColumn, actualColumn)
	}
	for i, expectedIndex := range expected.IndexStructureList {
		actualIndex := actual.IndexStructureList[i]
		assert.Equal(t, expectedIndex, actualIndex)
	}
	assert.Equal(t, expected.Properties, actual.Properties)
}
