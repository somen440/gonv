package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// ConvertTableAll convert table
//   1. DROP
//   2. MODIFY
//   3. ADD
func (c *Converter) convertTableAll(
	before, after *structure.DatabaseStructure,
	ask *TableAsk,
) *migration.List {
	results := &migration.List{}

	results.Merge(
		c.ConvertTableDropMigration(before, after),
		c.convertModifyTableAll(before, after, ask),
		c.ConvertTableCreateMigration(before, after),
	)

	return results
}

// ConvertTableDropMigration converter
func (c *Converter) ConvertTableDropMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	return results
}

// convertModifyTableAll
func (c *Converter) convertModifyTableAll(
	before, after *structure.DatabaseStructure,
	ask *TableAsk,
) *migration.List {
	results := &migration.List{}

	beforeList := before.ListToFilterTableType()
	afterList := after.ListToFilterTableType()

	for beforeTable, beforeSt := range beforeList {
		afterSt, ok := afterList[beforeTable]
		if !ok {
			continue
		}
		migration := c.convertModifyTable(beforeSt, afterSt, ask)
		results.Add(migration)
	}

	for beforeTable, afterTable := range ask.RenamedTableList {
		beforeSt := beforeList[beforeTable]
		afterSt := afterList[afterTable]
		migration := c.convertModifyTable(beforeSt, afterSt, ask)
		results.Add(migration)
	}

	return results
}

func (c *Converter) convertModifyTable(
	before, after *structure.TableStructure,
	ask *TableAsk,
) *migration.TableAlterMigration {
	return migration.NewTableAlterMigration(
		before.Table,
		after.Table,
		c.convertTableMigrationLineList(before, after, ask),
		ask.RenamedColumnListAsStrings(),
		c.ConvertTablePartitionMigration(before, after),
	)
}

// ConvertTableCreateMigration convert DatabaseStructure -> TableCreateMigration
func (c *Converter) ConvertTableCreateMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	afterAll := after.ListToFilterTableType()
	unknowns := after.DiffListToFilterTableType(before)

	if len(unknowns) == 0 {
		return nil
	}

	for table := range unknowns {
		migration := migration.NewTableCreateMigration(afterAll[table])
		results.Add(migration)
	}

	return results
}
