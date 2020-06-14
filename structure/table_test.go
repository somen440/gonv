/*
Copyright 2020 somen440

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTableSt() *TableStructure {
	return &TableStructure{
		Table:          "sample",
		Type:           TableType,
		Comment:        "sample table",
		Engine:         "InnoDB",
		DefaultCharset: "utf8mb4",
		Collate:        "utf8mb4_unicode_ci",
		Properties:     []string{},
		ColumnStructureList: map[ColumnField]*MySQL57ColumnStructure{
			ColumnField("id"): {
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
			ColumnField("name"): {
				Field:                ColumnField("name"),
				Type:                 "varchar(255)",
				Default:              "",
				Comment:              "Sample Name",
				Attributes:           []Attribute{},
				CollationName:        "utf8mb4_unicode_ci",
				Properties:           []string{},
				GenerationExpression: "",
			},
		},
		IndexStructureList: map[IndexKey]*IndexStructure{},
		Partition:          nil,
	}
}

func TestGenerateModifiedColumnStructureSetMap(t *testing.T) {
	t1 := createTableSt()
	t2 := createTableSt()

	t2.ColumnStructureList[ColumnField("id")].Comment = "Sample ID dayo"
	t2.ColumnStructureList[ColumnField("id")].Attributes = []Attribute{Unsigned, Nullable}
	t2.ColumnStructureList[ColumnField("id")].Default = "10"
	t2.ColumnStructureList[ColumnField("name")].Field = ColumnField("fullname")
	t2.ColumnStructureList[ColumnField("name")].Default = "sample"
	t2.ColumnStructureList[ColumnField("name")].Type = "text"

	setMap, err := t1.GenerateModifiedColumnStructureSetMap(t2, RenamedField{
		"name": "fullname",
	})
	assert.Nil(t, err)

	tests := []struct {
		field        ColumnField
		expectedUp   string
		expectedDown string
	}{
		{
			ColumnField("id"),
			"CHANGE `id` `id` bigint(20) unsigned DEFAULT 10 COMMENT 'Sample ID dayo'",
			"CHANGE `id` `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Sample ID'",
		},
		{
			ColumnField("name"),
			"CHANGE `fullname` `fullname` text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name'",
			"CHANGE `name` `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Sample Name'",
		},
	}

	for _, tt := range tests {
		modified, ok := setMap[tt.field]
		if !ok {
			assert.Fail(t, "missiong "+string(tt.field)+" modified.")
		}

		actualUp := modified.Up.GenerateChangeQuery()
		assert.Equal(t, tt.expectedUp, actualUp)

		actualDown := modified.Down.GenerateChangeQuery()
		assert.Equal(t, tt.expectedDown, actualDown)
	}
}
