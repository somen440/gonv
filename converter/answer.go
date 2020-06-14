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

package converter

import "github.com/somen440/gonv/structure"

// TableAnswer answer table
type TableAnswer struct {
	DroppedTableList []structure.TableName
	RenamedTableList map[structure.TableName]structure.TableName

	DroppedColumnList []structure.ColumnField
	RenamedColumnList map[structure.ColumnField]structure.ColumnField
}

// RenamedColumnListAsStrings to strings
func (a *TableAnswer) RenamedColumnListAsStrings() []string {
	results := []string{}

	for _, table := range a.RenamedColumnList {
		results = append(results, string(table))
	}

	return results
}

// ViewAnswer answer view
type ViewAnswer struct {
	// todo: []ViewName ... etc
}

// ModifiedAnswer answer modofied
type ModifiedAnswer struct {
	Table *TableAnswer
	View  *ViewAnswer
}

// NewModifiedAnswer return ModifiedAnswer
func NewModifiedAnswer(table *TableAnswer, view *ViewAnswer) *ModifiedAnswer {
	return &ModifiedAnswer{
		Table: table,
		View:  view,
	}
}
