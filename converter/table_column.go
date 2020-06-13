package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
	"github.com/somen440/gonv/util"
)

func (c *Converter) toTableMigrationLineList(
	before, after *structure.TableStructure,
	ask *TableAsk,
) *migration.LineList {
	results := migration.NewMigrationLineList()

	indexAll := c.toIndexAllMigrationLine(before, after)

	if line := c.toTableRenameMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toTableCommentMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toTableEngineMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toTableDefaultCharsetMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toTableCollateMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toTableCollateMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if indexAll != nil {
		if line := indexAll.First; line != nil {
			results.Add(line)
		}
	}
	if line := c.toColumnDropMigrationLine(before, after, ask); line != nil {
		results.Add(line)
	}
	if line := c.toColumnModifyMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if line := c.toColumnAddMigrationLine(before, after); line != nil {
		results.Add(line)
	}
	if indexAll != nil {
		if line := indexAll.Last; line != nil {
			results.Add(line)
		}
	}

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
	ask *TableAsk,
) *migration.ColumnDropMigrationLine {
	dropped := ask.DroppedColumnList
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

func (c *Converter) toColumnModifyMigrationLine(before, after *structure.TableStructure) *migration.ColumnModifyMigrationLine {
	return nil
}

func (c *Converter) toColumnAddMigrationLine(before, after *structure.TableStructure) *migration.ColumnAddMigrationLine {
	return nil
}

func (c *Converter) toTablePartitionMigration(before, after *structure.TableStructure) migration.PartitionMigration {
	return nil
}
