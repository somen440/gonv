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

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// Converter structure -> migration converter
type Converter struct {
	Err error
}

// NewConverter create Converter
func NewConverter() *Converter {
	return &Converter{
		Err: nil,
	}
}

// HasError error not eq nil return true
func (c *Converter) HasError() bool {
	return c.Err != nil
}

// ConvertAll cnvert struct -> migration all
//   1. DROP
//   2. MODIFY
//     2.1. DROP Index
//     2.2. DROP
//     2.3. MODIFY
//     2.4. ADD
//     2.5. ADD
//   3. ADD
func (c *Converter) ConvertAll(
	before, after *structure.DatabaseStructure,
	a *ModifiedAnswer,
) *migration.List {
	results := &migration.List{}

	tableAnswer := &TableAnswer{}
	if a != nil {
		tableAnswer = a.Table
	}
	viewAnser := &ViewAnswer{}
	if a != nil {
		viewAnser = a.View
	}

	// table
	results.Merge(
		c.ToTableDropMigration(before, after, tableAnswer),
		c.ToTableAlterMigrationAll(before, after, tableAnswer),
		c.ToTableCreateMigration(before, after),
	)

	// view
	results.Merge(
		c.ToViewDropMigration(before, after),
		c.ToViewAlterMigration(before, after),
		c.ToViewRenameMigration(before, after, viewAnser),
		c.ToViewCreateMigration(before, after),
	)

	return results
}
