package migration

import (
	"fmt"
	"strings"
)

// LineList list
type LineList struct {
	list []Line
}

// NewMigrationLineList create MigrationLineList
func NewMigrationLineList() *LineList {
	m := &LineList{
		list: []Line{},
	}
	return m
}

// Add add line
func (ll *LineList) Add(line Line) {
	ll.list = append(ll.list, line)
}

// IsMigratable is migratable return true
func (ll *LineList) IsMigratable() bool {
	return len(ll.list) > 0
}

// Up return up query
func (ll *LineList) Up() string {
	upLineList := []string{}
	for _, line := range ll.list {
		upLineList = append(upLineList, line.UpList()...)
	}
	return " " + strings.Join(upLineList, ",\n ")
}

// Down return down query
func (ll *LineList) Down() string {
	downLineList := []string{}
	for _, line := range ll.list {
		downLineList = append(downLineList, line.DownList()...)
	}
	return " " + strings.Join(downLineList, ",\n ")
}

// Merge args
func (ll *LineList) Merge(args ...Line) {
	for _, line := range args {
		if line == nil {
			continue
		}
		ll.list = append(ll.list, line)
	}
}

// TableAlterMigration ALTER TABLE
type TableAlterMigration struct {
	tableMigration

	IsAltered       bool
	RenamedNameList []string
}

type PartitionMigration interface {
	UpQuery() string
	DownQuery() string
}

// NewTableAlterMigration create TableAlterMigration
func NewTableAlterMigration(
	beforeTableName, afterTableName string,
	lineList *LineList, renamedNameList []string,
	partitionMigration PartitionMigration,
) *TableAlterMigration {
	migration := &TableAlterMigration{}

	migration.Table = beforeTableName
	migration.Type = AlterType
	migration.RenamedNameList = renamedNameList
	migration.IsAltered = lineList.IsMigratable() || partitionMigration != nil

	if !migration.IsAltered {
		return migration
	}

	migration.Up = "ALTER TABLE `" + beforeTableName + "`"
	migration.Down = "ALTER TABLE `" + afterTableName + "`"

	if lineList.IsMigratable() {
		migration.Up += "\n" + lineList.Up()
		migration.Down += "\n" + lineList.Down()
	}

	if partitionMigration != nil {
		migration.Up += "\n" + partitionMigration.UpQuery()
		migration.Down += "\n" + partitionMigration.DownQuery()
	}

	migration.Up += ";"
	migration.Down += ";"

	return migration
}

// TableCreateMigration CREATE TABLE
type TableCreateMigration struct {
	tableMigration
}

// NewTableCreateMigration create TableCreateMigration
func NewTableCreateMigration(ts TableStructure) *TableCreateMigration {
	migration := &TableCreateMigration{}

	migration.Table = ts.GetTable()
	migration.Type = CreateType

	migration.Up = "CREATE TABLE `" + ts.GetTable() + "` (\n"

	bodies := []string{}
	for _, column := range ts.GetColumnStructureList() {
		bodies = append(bodies, column.GenerateCreateQuery())
	}
	for _, index := range ts.GetIndexStructureList() {
		bodies = append(bodies, index.GenerateCreateQuery())
	}
	migration.Up += " " + strings.Join(bodies, ",\n ") + "\n"

	migration.Up += ")"
	migration.Up += " ENGINE=" + ts.GetEngine()
	migration.Up += " DEFAULT CHARASET=" + ts.GetDefaultCharset()
	migration.Up += " COLLATE=" + ts.GetCollate()
	migration.Up += " COMMENT='" + ts.GetComment() + "'"

	if ts.GetPartition() != nil {
		migration.Up += fmt.Sprintf("\n/*!50100 %s */", ts.GetPartition().Query())
	}

	migration.Down = "DROP TABLE " + ts.GetTable()

	return migration
}

// TableDropMigration DROP TABLE
type TableDropMigration struct {
	tableMigration
}

// NewTableDropMigration create TableDropMigration
func NewTableDropMigration(ts TableStructure) *TableDropMigration {
	migration := &TableDropMigration{}
	migration.Table = ts.GetTable()
	migration.Type = DropType
	creation := NewTableCreateMigration(ts)
	migration.Up = creation.Down
	migration.Down = creation.Up
	return migration
}

// ViewAlterMigration CREATE OR REPLACE
type ViewAlterMigration struct {
	tableMigration

	IsAltered bool
	IsSplit   bool
}

// NewViewAlterMigration return ViewAlterMigration
func NewViewAlterMigration(before, after ViewStructure, allRenamedNameList [][]string) *ViewAlterMigration {
	migration := &ViewAlterMigration{}
	migration.Table = after.GetName()
	migration.Type = CreateOrReplaceType
	migration.IsAltered = before.CompareQuery() != after.CompareQuery()
	if !migration.IsAltered {
		return migration
	}

	migration.Up = strings.Replace(after.CompareQuery(), "CREATE", "CREATE OR REPLACE", -1)
	migration.Down = strings.Replace(before.CompareQuery(), "CREATE", "CREATE OR REPLACE", -1)

	for _, nameList := range allRenamedNameList {
		count := 0
		for _, name := range nameList {
			if strings.Contains(migration.Up, name) {
				count++
			}
		}
		if count == len(nameList) {
			migration.IsSplit = true
		}
	}

	migration.Down = strings.Replace(migration.Down, before.GetName(), after.GetName(), -1)

	return migration
}

// ViewCreateMigration CREATE VIEW
type ViewCreateMigration struct {
	tableMigration
}

// NewViewCreateMigration create ViewCreateMigration
func NewViewCreateMigration(view ViewStructure) *ViewCreateMigration {
	migration := &ViewCreateMigration{}
	migration.Table = view.GetName()
	migration.Type = ViewCreateType
	migration.Up = view.CreateQueryToFormat()
	migration.Down = "DROP VIEW " + view.GetName() + ";"
	return migration
}

// ViewDropMigration DROP VIEW
type ViewDropMigration struct {
	tableMigration
}

// NewViewDropMigration return ViewDropMigration
func NewViewDropMigration(view ViewStructure) *ViewDropMigration {
	migration := &ViewDropMigration{}
	migration.Table = view.GetName()
	migration.Type = ViewDropType
	create := NewViewCreateMigration(view)
	migration.Up = create.Down
	migration.Down = create.Up
	return migration
}

// ViewRenameMigration RENAME TABLE ~ TO ~
type ViewRenameMigration struct {
	tableMigration

	IsAltered bool
}

// NewViewRenameMigration create ViewRenameMigration
func NewViewRenameMigration(before, after ViewStructure) *ViewRenameMigration {
	migration := &ViewRenameMigration{}
	migration.Table = after.GetName()
	migration.Type = ViewRenameType
	migration.Up = "RENAME TABLE " + before.GetName() + " TO " + after.GetName() + ";"
	migration.Down = "RENAME TABLE " + after.GetName() + " TO " + before.GetName() + ";"
	return migration
}
