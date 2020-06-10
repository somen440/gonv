package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// ToTableDropMigration DatabaseStructure -> TableDropMigration
func (c *Converter) ToTableDropMigration(before, after *structure.DatabaseStructure) *migration.List {
	results := &migration.List{}

	return results
}

// ToTableAlterMigrationAll DatabaseStructure -> TableAlterMigration
func (c *Converter) ToTableAlterMigrationAll(
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
		migration := c.toTableAlterMigration(beforeSt, afterSt, ask)
		results.Add(migration)
	}

	for beforeTable, afterTable := range ask.RenamedTableList {
		beforeSt := beforeList[beforeTable]
		afterSt := afterList[afterTable]
		migration := c.toTableAlterMigration(beforeSt, afterSt, ask)
		results.Add(migration)
	}

	return results
}

func (c *Converter) toTableAlterMigration(
	before, after *structure.TableStructure,
	ask *TableAsk,
) *migration.TableAlterMigration {
	return migration.NewTableAlterMigration(
		before.Table,
		after.Table,
		c.toTableMigrationLineList(before, after, ask),
		ask.RenamedColumnListAsStrings(),
		c.toTablePartitionMigration(before, after),
	)
}

// ToTableCreateMigration DatabaseStructure -> TableCreateMigration
func (c *Converter) ToTableCreateMigration(before, after *structure.DatabaseStructure) *migration.List {
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
