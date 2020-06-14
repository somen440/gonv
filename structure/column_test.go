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
