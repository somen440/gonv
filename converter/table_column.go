package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

func (c *Converter) toTableMigrationLineList(
	before, after *structure.TableStructure,
	ask *TableAsk,
) *migration.LineList {
	results := migration.NewMigrationLineList()

	indexAll := c.toIndexAllMigrationLine(before, after)

	results.Merge(
		c.toTableRenameMigrationLine(before, after),
		c.toTableCommentMigrationLine(before, after),
		c.toTableEngineMigrationLine(before, after),
		c.toTableDefaultCharsetMigrationLine(before, after),
		c.toTableCollateMigrationLine(before, after),
		indexAll.First,
		c.toColumnDropMigrationLine(before, after),
		c.toColumnModifyMigrationLine(before, after),
		c.toColumnAddMigrationLine(before, after),
		indexAll.Last,
	)

	return results
}

func (c *Converter) toTableRenameMigrationLine(before, after *structure.TableStructure) *migration.TableRenameMigrationLine {
	return nil
}

func (c *Converter) toTableCommentMigrationLine(before, after *structure.TableStructure) *migration.TableCommentMigrationLine {
	return nil
}

func (c *Converter) toTableEngineMigrationLine(before, after *structure.TableStructure) *migration.TableEngineMigrationLine {
	return nil
}

func (c *Converter) toTableDefaultCharsetMigrationLine(before, after *structure.TableStructure) *migration.TableDefaultCharsetMigrationLine {
	return nil
}

func (c *Converter) toTableCollateMigrationLine(before, after *structure.TableStructure) *migration.TableCollateMigrationLine {
	return nil
}

func (c *Converter) toIndexAllMigrationLine(before, after *structure.TableStructure) *migration.IndexAllMigrationLine {
	return nil
}

func (c *Converter) toColumnDropMigrationLine(before, after *structure.TableStructure) *migration.ColumnDropMigrationLine {
	return nil
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
