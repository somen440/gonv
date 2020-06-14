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

package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFieldOrder(t *testing.T) {
	tests := []struct {
		title    string
		before   []string
		after    []string
		expected FieldOrderMap
	}{
		{
			"完全一致",
			[]string{"A", "B", "C", "D", "E"},
			[]string{"A", "B", "C", "D", "E"},
			FieldOrderMap{},
		},
		{
			"Aの前のBがDの後ろに移動",
			[]string{"A", "B", "C", "D", "E"},
			[]string{"A", "C", "D", "B", "E"},
			FieldOrderMap{
				"B": &FieldOrder{
					Field:              "B",
					PreviousAfterField: "A",
					NextAfterField:     "D",
				},
			},
		},
		{
			"AがBの後ろに移動",
			[]string{"A", "B", "C", "D", "E"},
			[]string{"B", "C", "D", "E", "A"},
			FieldOrderMap{
				"A": &FieldOrder{
					Field:              "A",
					PreviousAfterField: "",
					NextAfterField:     "E",
				},
			},
		},
		{
			"Eが先頭に移動",
			[]string{"A", "B", "C", "D", "E"},
			[]string{"E", "A", "B", "C", "D"},
			FieldOrderMap{
				"E": &FieldOrder{
					Field:              "E",
					PreviousAfterField: "D",
					NextAfterField:     "",
				},
			},
		},
		{
			"Eが先頭に移動",
			[]string{"A", "B", "C", "D", "E"},
			[]string{"E", "B", "C", "D", "A"},
			FieldOrderMap{
				"A": &FieldOrder{
					Field:              "A",
					PreviousAfterField: "",
					NextAfterField:     "D",
				},
				"E": &FieldOrder{
					Field:              "E",
					PreviousAfterField: "D",
					NextAfterField:     "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			actual := GenerateFieldOrderList(tt.before, tt.after)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
