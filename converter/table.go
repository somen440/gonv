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
	"fmt"

	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// ToTableDropMigration DatabaseStructure -> TableDropMigration
func (c *Converter) ToTableDropMigration(
	before, after *structure.DatabaseStructure,
	a *TableAnswer,
) *migration.List {
	if c.HasError() {
		return nil
	}

	results := &migration.List{}

	if len(a.DroppedTableList) == 0 {
		return nil
	}

	beforeAll := before.ListToFilterTableType()
	afterAll := after.ListToFilterTableType()

	for _, table := range a.DroppedTableList {
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
	a *TableAnswer,
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
		migration := c.toTableAlterMigration(beforeSt, afterSt, a)
		if migration.IsAltered {
			results.Add(migration)
		}
	}

	for beforeTable, afterTable := range a.RenamedTableList {
		beforeSt, ok := beforeList[beforeTable]
		if !ok {
			c.Err = fmt.Errorf("ToTableAlterMigrationAll not found table %s from before %v", beforeTable, beforeList)
			return nil
		}
		afterSt, ok := afterList[afterTable]
		if !ok {
			c.Err = fmt.Errorf("ToTableAlterMigrationAll not found table %s from before %v", afterTable, afterList)
			return nil
		}
		migration := c.toTableAlterMigration(beforeSt, afterSt, a)
		if migration.IsAltered {
			results.Add(migration)
		}
	}

	return results
}

func (c *Converter) toTableAlterMigration(
	before, after *structure.TableStructure,
	a *TableAnswer,
) *migration.TableAlterMigration {
	return migration.NewTableAlterMigration(
		before.Table,
		after.Table,
		c.toTableMigrationLineList(before, after, a),
		a.RenamedColumnListAsStrings(),
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
