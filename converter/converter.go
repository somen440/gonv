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
