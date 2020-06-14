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
