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

package migration

type migrationLine struct {
	Line

	upLineList   []string
	downLineList []string
}

func (line *migrationLine) UpList() []string {
	return line.upLineList
}

func (line *migrationLine) DownList() []string {
	return line.downLineList
}

// ColumnAddMigrationLine ALTER TABLE ~ ADD ~
type ColumnAddMigrationLine struct {
	migrationLine
}

// NewColumnAddMigrationLine create ColumnAddMigrationLine
func NewColumnAddMigrationLine(list []ModifiedColumnStructure) *ColumnAddMigrationLine {
	if len(list) == 0 {
		return nil
	}

	line := &ColumnAddMigrationLine{}

	for _, column := range list {
		line.upLineList = append(line.upLineList, column.GenerateAddQuery())
		line.downLineList = append(line.downLineList, column.GetColumn().GenerateDropQuery())
	}

	return line
}

// ColumnDropMigrationLine ALTER TABLE ~ DROP ~
type ColumnDropMigrationLine struct {
	migrationLine
}

// NewColumnDropMigrationLine create ColumnDropMigrationLine
func NewColumnDropMigrationLine(columns []ModifiedColumnStructure) *ColumnDropMigrationLine {
	if len(columns) == 0 {
		return nil
	}

	line := &ColumnDropMigrationLine{}

	cam := NewColumnAddMigrationLine(columns)
	line.upLineList = cam.downLineList
	line.downLineList = cam.upLineList

	return line
}

// ColumnModifyMigrationLine ALTER TABLE ~ MODIFY ~
type ColumnModifyMigrationLine struct {
	migrationLine
}

// NewColumnModifyMigrationLine create ColumnModifyMigrationLine
func NewColumnModifyMigrationLine(list []ModifiedColumnStructureSet) *ColumnModifyMigrationLine {
	if len(list) == 0 {
		return nil
	}
	line := &ColumnModifyMigrationLine{}
	for _, set := range list {
		line.upLineList = append(line.upLineList, set.UpStructure().GenerateChangeQuery())
		line.downLineList = append(line.downLineList, set.DownStructure().GenerateChangeQuery())
	}
	return line
}

// IndexAddMigrationLine ALTER TABLE ~ KEY ~
type IndexAddMigrationLine struct {
	migrationLine
}

// NewIndexAddMigrationLine create IndexAddMigrationLine
func NewIndexAddMigrationLine(list []IndexStructure) *IndexAddMigrationLine {
	if len(list) == 0 {
		return nil
	}
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
	if len(list) == 0 {
		return nil
	}
	line := &IndexDropMigrationLine{}

	iam := NewIndexAddMigrationLine(list)
	line.upLineList = iam.downLineList
	line.downLineList = iam.upLineList

	return line
}

// IndexAllMigrationLine ALTER TABLE ~ KEY ~
type IndexAllMigrationLine struct {
	First *IndexDropMigrationLine
	Last  *IndexAddMigrationLine
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
	PartitionMigration

	Up   string
	Down string
}

func (pm *partitionMigration) UpQuery() string {
	return pm.Up
}

func (pm *partitionMigration) DownQuery() string {
	return pm.Down
}

// PartitionRemoveMigration REMOVE PARTITIONING
type PartitionRemoveMigration struct {
	partitionMigration
}

// NewPartitionRemoveMigration create PartitionRemoveMigration
func NewPartitionRemoveMigration(before PartitionStructure) *PartitionRemoveMigration {
	if before == nil {
		return nil
	}
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
func NewPartitionResetMigration(before, after PartitionStructure) *PartitionResetMigration {
	if before == nil || after == nil {
		return nil
	}
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
	line.upLineList = append(line.upLineList, "COMMENT '"+after+"'")
	line.downLineList = append(line.downLineList, "COMMENT '"+before+"'")
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

// TableRenameMigrationLine ALTER TABLE ~ RENAME TO ~
type TableRenameMigrationLine struct {
	migrationLine
}

// NewTableRenameMigrationLine create TableRenameMigrationLine
func NewTableRenameMigrationLine(before, after string) *TableRenameMigrationLine {
	migration := &TableRenameMigrationLine{}
	migration.upLineList = append(migration.upLineList, "RENAME TO `"+after+"`")
	migration.downLineList = append(migration.downLineList, "RENAME TO `"+before+"`")
	return migration
}
