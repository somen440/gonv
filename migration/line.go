package migration

import "github.com/somen440/gonv/structure"

type migrationLine struct {
	upLineList   []string
	downLineList []string
}

func (line *migrationLine) Up() []string {
	return line.upLineList
}

func (line *migrationLine) Down() []string {
	return line.upLineList
}

// ColumnAddMigrationLine ALTER TABLE ~ ADD ~
type ColumnAddMigrationLine struct {
	migrationLine
}

// NewColumnAddMigrationLine create ColumnAddMigrationLine
func NewColumnAddMigrationLine(list []ModifiedColumnStructure) *ColumnAddMigrationLine {
	line := &ColumnAddMigrationLine{}

	for _, column := range list {
		line.upLineList = append(line.upLineList, column.GenerateAddQuery())
		line.downLineList = append(line.downLineList, column.GetColumn().GenerateBaseQuery())
	}

	return line
}

// ColumnDropMigrationLine ALTER TABLE ~ DROP ~
type ColumnDropMigrationLine struct {
	migrationLine
}

// NewColumnDropMigrationLine create ColumnDropMigrationLine
func NewColumnDropMigrationLine(columns []ModifiedColumnStructure) *ColumnDropMigrationLine {
	line := &ColumnDropMigrationLine{}

	cam := NewColumnAddMigrationLine(columns)
	line.upLineList = cam.downLineList
	line.downLineList = cam.upLineList

	return line
}

// ModifiedColumnStructureSet up down set
type ModifiedColumnStructureSet interface {
	UpColumn() ModifiedColumnStructure
	DownColumn() ModifiedColumnStructure
}

// ColumnModifyMigrationLine ALTER TABLE ~ MODIFY ~
type ColumnModifyMigrationLine struct {
	migrationLine
}

// NewColumnModifyMigrationLine create ColumnModifyMigrationLine
func NewColumnModifyMigrationLine(list []ModifiedColumnStructureSet) *ColumnModifyMigrationLine {
	line := &ColumnModifyMigrationLine{}

	for _, set := range list {
		line.upLineList = append(line.upLineList, set.UpColumn().GenerateAddQuery())
		line.downLineList = append(line.downLineList, set.DownColumn().GenerateAddQuery())
	}

	return line
}

// IndexAddMigrationLine ALTER TABLE ~ KEY ~
type IndexAddMigrationLine struct {
	migrationLine
}

// NewIndexAddMigrationLine create IndexAddMigrationLine
func NewIndexAddMigrationLine(list []IndexStructure) *IndexAddMigrationLine {
	line := &IndexAddMigrationLine{}

	for _, i := range list {
		line.upLineList = append(line.upLineList, i.GenerateAddQuery())
		line.downLineList = append(line.downLineList, i.GenerateDropQuery())
	}

	return line
}

// IndexDropMigrationLine ALTER TABLE ~ KEY ~
type IndexDropMigrationLine struct {
	migrationLine
}

// NewIndexIndexDropMigrationLine create IndexDropMigrationLine
func NewIndexIndexDropMigrationLine(list []IndexStructure) *IndexDropMigrationLine {
	line := &IndexDropMigrationLine{}

	iam := NewIndexAddMigrationLine(list)
	line.upLineList = iam.downLineList
	line.downLineList = iam.upLineList

	return line
}

// IndexAllMigrationLine ALTER TABLE ~ KEY ~
type IndexAllMigrationLine struct {
	First *IndexAddMigrationLine
	Last  *IndexDropMigrationLine
}

// IsFirstExist exist first return true
func (line *IndexAllMigrationLine) IsFirstExist() bool {
	return line.First != nil
}

// IsLastExist exist last return true
func (line *IndexAllMigrationLine) IsLastExist() bool {
	return line.Last != nil
}

type partitionMigration struct {
	Up   string
	Down string
}

// PartitionRemoveMigration REMOVE PARTITIONING
type PartitionRemoveMigration struct {
	partitionMigration
}

// NewPartitionRemoveMigration create PartitionRemoveMigration
func NewPartitionRemoveMigration(before structure.PartitionStructure) *PartitionRemoveMigration {
	p := &PartitionRemoveMigration{}
	p.Up = "REMOVE PARTITIONING"
	p.Down = before.Query()
	return p
}

// PartitionResetMigration PARTITION BY ~
type PartitionResetMigration struct {
	partitionMigration
}

// NewPartitionResetMigration create PartitionResetMigration
func NewPartitionResetMigration(before, after structure.PartitionStructure) *PartitionResetMigration {
	p := &PartitionResetMigration{}
	p.Up = after.Query()
	p.Down = before.Query()
	return p
}

// TableCollateMigrationLine ALTER TABLE ~ COLLATE ~
type TableCollateMigrationLine struct {
	migrationLine
}

// NewTableCollateMigrationLine create TableCollateMigrationLine
func NewTableCollateMigrationLine(before, after string) *TableCollateMigrationLine {
	line := &TableCollateMigrationLine{}
	line.upLineList = append(line.upLineList, "COLLATE "+after)
	line.downLineList = append(line.downLineList, "COLLATE "+before)
	return line
}

// TableCommentMigrationLine ALTER TABLE ~ COMMENT ~
type TableCommentMigrationLine struct {
	migrationLine
}

// NewTableCommentMigrationLine create TableCommentMigrationLine
func NewTableCommentMigrationLine(before, after string) *TableCommentMigrationLine {
	line := &TableCommentMigrationLine{}
	line.upLineList = append(line.upLineList, "COMMENT "+after)
	line.downLineList = append(line.downLineList, "COMMENT "+before)
	return line
}

// TableDefaultCharsetMigrationLine ALTER TABLE ~ DEFAULT CHARSET ~
type TableDefaultCharsetMigrationLine struct {
	migrationLine
}

// NewTableDefaultCharsetMigrationLine create TableDefaultCharsetMigrationLine
func NewTableDefaultCharsetMigrationLine(before, after string) *TableDefaultCharsetMigrationLine {
	line := &TableDefaultCharsetMigrationLine{}
	line.upLineList = append(line.upLineList, "DEFAULT CHARSET "+after)
	line.downLineList = append(line.downLineList, "DEFAULT CHARASET "+before)
	return line
}

// TableEngineMigrationLine ALTER TABLE ~ ENGINE ~
type TableEngineMigrationLine struct {
	migrationLine
}

// NewTableEngineMigrationLine create TableEngineMigrationLine
func NewTableEngineMigrationLine(before, after string) *TableEngineMigrationLine {
	line := &TableEngineMigrationLine{}
	line.upLineList = append(line.upLineList, "ENGINE "+after)
	line.downLineList = append(line.downLineList, "ENGINE "+before)
	return line
}
