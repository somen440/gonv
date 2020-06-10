package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// ToViewDropMigration DatabaseStructure -> ViewDropMigration
func (c *Converter) ToViewDropMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	return results
}

// ToViewAlterMigration DatabaseStructure -> ViewAlterMigration
func (c *Converter) ToViewAlterMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	return results
}

// ToViewRenameMigration DatabaseStructure -> ViewRenameMigration
func (c *Converter) ToViewRenameMigration(
	before, after *structure.DatabaseStructure,
	ask *ViewAsk,
) *migration.List {
	results := &migration.List{}

	return results
}

// ToViewCreateMigration DatabaseStructure -> ViewCreateMigration
func (c *Converter) ToViewCreateMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	return results
}
