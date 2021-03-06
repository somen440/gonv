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

// ToViewDropMigration DatabaseStructure -> ViewDropMigration
func (c *Converter) ToViewDropMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	// todo: drop view #2

	return results
}

// ToViewAlterMigration DatabaseStructure -> ViewAlterMigration
func (c *Converter) ToViewAlterMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	// todo: alter view #3

	return results
}

// ToViewRenameMigration DatabaseStructure -> ViewRenameMigration
func (c *Converter) ToViewRenameMigration(
	before, after *structure.DatabaseStructure,
	a *ViewAnswer,
) *migration.List {
	results := &migration.List{}

	// todo: rename view #4

	return results
}

// ToViewCreateMigration DatabaseStructure -> ViewCreateMigration
func (c *Converter) ToViewCreateMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	// todo: create view #5

	return results
}
