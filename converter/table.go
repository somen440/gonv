package converter

import (
	"fmt"

	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// ToTableDropMigration DatabaseStructure -> TableDropMigration
func (c *Converter) ToTableDropMigration(
	before, after *structure.DatabaseStructure,
	ask *TableAsk,
) *migration.List {
	if c.HasError() {
		return nil
	}

	results := &migration.List{}

	if len(ask.DroppedTableList) == 0 {
		return nil
	}

	beforeAll := before.ListToFilterTableType()
	afterAll := after.ListToFilterTableType()

	for _, table := range ask.DroppedTableList {
		before, ok := beforeAll[table]
		if !ok {
			c.Err = fmt.Errorf("ToTableDropMigration not found table %s from %v", table, beforeAll)
			return nil
		}
		_, ok = afterAll[table]
		if ok {
			c.Err = fmt.Errorf("ToTableDropMigration found table %s from %v", table, afterAll)
			return nil
		}
		migration := migration.NewTableDropMigration(before)
		results.Add(migration)
	}

	return results
}

// ToTableAlterMigrationAll DatabaseStructure -> TableAlterMigration
func (c *Converter) ToTableAlterMigrationAll(
	before, after *structure.DatabaseStructure,
	ask *TableAsk,
) *migration.List {
	if c.HasError() {
		return nil
	}

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
	if c.HasError() {
		return nil
	}

	results := &migration.List{}

	afterAll := after.ListToFilterTableType()
	unknowns := after.DiffListToFilterTableType(before)

	if len(unknowns) == 0 {
		return nil
	}

	for table := range unknowns {
		after, ok := afterAll[table]
		if !ok {
			c.Err = fmt.Errorf("ToTableCreateMigration not found table %s from after %v", table, afterAll)
			return nil
		}
		migration := migration.NewTableCreateMigration(after)
		results.Add(migration)
	}

	return results
}
