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
	"github.com/somen440/gonv/util"
)

func (c *Converter) toTableMigrationLineList(
	before, after *structure.TableStructure,
	a *TableAnswer,
) *migration.LineList {
	if c.HasError() {
		return nil
	}

	results := migration.NewMigrationLineList()

	var indexFirst *migration.IndexDropMigrationLine
	var indexLast *migration.IndexAddMigrationLine
	indexAll := c.toIndexAllMigrationLine(before, after)
	if indexAll != nil {
		indexFirst = indexAll.First
		indexLast = indexAll.Last
	}

	results.Merge(
		c.toTableRenameMigrationLine(before, after),
		c.toTableCommentMigrationLine(before, after),
		c.toTableEngineMigrationLine(before, after),
		c.toTableDefaultCharsetMigrationLine(before, after),
		c.toTableCollateMigrationLine(before, after),
		c.toTableCollateMigrationLine(before, after),
		indexFirst,
		c.toColumnDropMigrationLine(before, after, a),
		c.toColumnModifyMigrationLine(before, after, a),
		c.toColumnAddMigrationLine(before, after),
		indexLast,
	)

	return results
}

func (c *Converter) toTableRenameMigrationLine(before, after *structure.TableStructure) *migration.TableRenameMigrationLine {
	if c.HasError() {
		return nil
	}

	if before.Table == after.Table {
		return nil
	}

	return migration.NewTableRenameMigrationLine(before.Table, after.Table)
}

func (c *Converter) toTableCommentMigrationLine(before, after *structure.TableStructure) *migration.TableCommentMigrationLine {
	if c.HasError() {
		return nil
	}

	if before.Comment == after.Comment {
		return nil
	}

	return migration.NewTableCommentMigrationLine(before.Comment, after.Comment)
}

func (c *Converter) toTableEngineMigrationLine(before, after *structure.TableStructure) *migration.TableEngineMigrationLine {
	if c.HasError() {
		return nil
	}

	if before.Engine == after.Engine {
		return nil
	}

	return migration.NewTableEngineMigrationLine(before.Engine, after.Engine)
}

func (c *Converter) toTableDefaultCharsetMigrationLine(before, after *structure.TableStructure) *migration.TableDefaultCharsetMigrationLine {
	if c.HasError() {
		return nil
	}

	if before.DefaultCharset == after.DefaultCharset {
		return nil
	}

	return migration.NewTableDefaultCharsetMigrationLine(before.DefaultCharset, after.DefaultCharset)
}

func (c *Converter) toTableCollateMigrationLine(before, after *structure.TableStructure) *migration.TableCollateMigrationLine {
	if c.HasError() {
		return nil
	}

	if before.Collate == after.Collate {
		return nil
	}

	return migration.NewTableCollateMigrationLine(before.Collate, after.Collate)
}

func (c *Converter) toIndexAllMigrationLine(before, after *structure.TableStructure) *migration.IndexAllMigrationLine {
	if c.HasError() {
		return nil
	}
	droppedList := []migration.IndexStructure{}
	addedList := []migration.IndexStructure{}

	biList := before.IndexStructureList
	aiList := after.IndexStructureList

	for _, key := range util.MapDiffKeys(biList, aiList) {
		bVal := biList[structure.IndexKey(key)]
		droppedList = append(droppedList, bVal)
	}

	for aKey, aVal := range aiList {
		bVal, ok := biList[aKey]
		if !ok {
			addedList = append(addedList, aVal)
			continue
		}
		if aVal.IsChanged(bVal) {
			droppedList = append(droppedList, bVal)
			addedList = append(addedList, aVal)
		}
	}

	if len(droppedList) == 0 && len(addedList) == 0 {
		return nil
	}

	return &migration.IndexAllMigrationLine{
		First: migration.NewIndexIndexDropMigrationLine(droppedList),
		Last:  migration.NewIndexAddMigrationLine(addedList),
	}
}

func (c *Converter) toColumnDropMigrationLine(
	before, after *structure.TableStructure,
	a *TableAnswer,
) *migration.ColumnDropMigrationLine {
	if c.HasError() {
		return nil
	}
	dropped := a.DroppedColumnList
	for _, v := range util.MapDiffKeys(before.ColumnStructureList, after.ColumnStructureList) {
		dropped = append(dropped, structure.ColumnField(v))
	}
	list := before.GetModifiedColumnList(dropped)
	if len(list) == 0 {
		return nil
	}
	results := []migration.ModifiedColumnStructure{}
	for _, v := range list {
		results = append(results, v)
	}
	return migration.NewColumnDropMigrationLine(results)
}

func (c *Converter) toColumnModifyMigrationLine(
	before, after *structure.TableStructure,
	a *TableAnswer,
) *migration.ColumnModifyMigrationLine {
	if c.HasError() {
		return nil
	}

	results := []migration.ModifiedColumnStructureSet{}

	mSetList, err := before.GenerateModifiedColumnStructureSetMap(after, a.RenamedColumnList)
	if err != nil {
		c.Err = err
	}
	for _, mSet := range mSetList {
		results = append(results, mSet)
	}

	// todo: moved #6

	if len(results) == 0 {
		return nil
	}

	return migration.NewColumnModifyMigrationLine(results)
}

func (c *Converter) toColumnAddMigrationLine(
	before, after *structure.TableStructure,
) *migration.ColumnAddMigrationLine {
	if c.HasError() {
		return nil
	}
	added := []structure.ColumnField{}
	for _, v := range util.MapDiffKeys(after.ColumnStructureList, before.ColumnStructureList) {
		added = append(added, structure.ColumnField(v))
	}
	list := after.GetModifiedColumnList(added)
	if len(list) == 0 {
		return nil
	}
	results := []migration.ModifiedColumnStructure{}
	for _, v := range list {
		results = append(results, v)
	}
	return migration.NewColumnAddMigrationLine(results)
}

func (c *Converter) toTablePartitionMigration(before, after *structure.TableStructure) migration.PartitionMigration {
	// todo: partition #1
	return nil
}
