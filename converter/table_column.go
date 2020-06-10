package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

func (c *Converter) convertTableMigrationLineList(
	before, after *structure.TableStructure,
	ask *TableAsk,
) *migration.LineList {
	results := migration.NewMigrationLineList()

	indexAll := c.ConvertIndexAllMigrationLine(before, after)

	results.Merge(
		c.ConvertTableRenameMigrationLine(before, after),
		c.ConvertTableCommentMigrationLine(before, after),
		c.ConvertTableEngineMigrationLine(before, after),
		c.ConvertTableDefaultCharsetMigrationLine(before, after),
		c.ConvertTableCollateMigrationLine(before, after),
		indexAll.First,
		c.ConvertColumnDropMigrationLine(before, after),
		c.ConvertColumnModifyMigrationLine(before, after),
		c.ConverterColumnAddMigrationLine(before, after),
		indexAll.Last,
	)

	return results
}

// ConvertTableRenameMigrationLine convert
func (c *Converter) ConvertTableRenameMigrationLine(before, after *structure.TableStructure) *migration.TableRenameMigrationLine {
	return nil
}

// ConvertTableCommentMigrationLine convert
func (c *Converter) ConvertTableCommentMigrationLine(before, after *structure.TableStructure) *migration.TableCommentMigrationLine {
	return nil
}

// ConvertTableEngineMigrationLine convert
func (c *Converter) ConvertTableEngineMigrationLine(before, after *structure.TableStructure) *migration.TableEngineMigrationLine {
	return nil
}

// ConvertTableDefaultCharsetMigrationLine convert
func (c *Converter) ConvertTableDefaultCharsetMigrationLine(before, after *structure.TableStructure) *migration.TableDefaultCharsetMigrationLine {
	return nil
}

// ConvertTableCollateMigrationLine convert
func (c *Converter) ConvertTableCollateMigrationLine(before, after *structure.TableStructure) *migration.TableCollateMigrationLine {
	return nil
}

// ConvertIndexAllMigrationLine convert
func (c *Converter) ConvertIndexAllMigrationLine(before, after *structure.TableStructure) *migration.IndexAllMigrationLine {
	return nil
}

// ConvertColumnDropMigrationLine convert
func (c *Converter) ConvertColumnDropMigrationLine(before, after *structure.TableStructure) *migration.ColumnDropMigrationLine {
	return nil
}

// ConvertColumnModifyMigrationLine convert
func (c *Converter) ConvertColumnModifyMigrationLine(before, after *structure.TableStructure) *migration.ColumnModifyMigrationLine {
	return nil
}

// ConverterColumnAddMigrationLine convert
func (c *Converter) ConverterColumnAddMigrationLine(before, after *structure.TableStructure) *migration.ColumnAddMigrationLine {
	return nil
}

// ConvertTablePartitionMigration convert
func (c *Converter) ConvertTablePartitionMigration(before, after *structure.TableStructure) migration.PartitionMigration {
	return nil
}
