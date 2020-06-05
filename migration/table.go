package migration

import (
	"fmt"
	"strings"
)

type tableMigration struct {
	Table string
	Type  Type
	Up    string
	Down  string
}

type Line interface {
	Up() []string
	Down() []string
}

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
		upLineList = append(line.Up())
	}
	return " " + strings.Join(upLineList, ",\n ")
}

// Down return down query
func (ll *LineList) Down() string {
	downLineList := []string{}
	for _, line := range ll.list {
		downLineList = append(line.Up())
	}
	return " " + strings.Join(downLineList, ",\n ")
}

// TableAlterMigration ALTER TABLE
type TableAlterMigration struct {
	tableMigration

	IsAltered       bool
	RenamedNameList []string
}

type PartitionMigration interface {
	Up() string
	Down() string
}

// NewTableAlterMigration create TableAlterMigration
func NewTableAlterMigration(
	beforeTableName, afterTableName string,
	lineList LineList, renamedNameList []string,
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

	migration.Up = "ALTER TABLE " + beforeTableName
	migration.Down = "ALTER TABLE " + afterTableName

	if lineList.IsMigratable() {
		migration.Up += "\n" + lineList.Up()
		migration.Down += "\n" + lineList.Down()
	}

	if partitionMigration != nil {
		migration.Up += "\n" + partitionMigration.Up()
		migration.Down += "\n" + partitionMigration.Down()
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

	migration.Table = ts.Table()
	migration.Type = CreateType

	migration.Up = "CREATE TABLE " + ts.Table() + " (\n"

	bodies := []string{}
	for _, column := range ts.ColumnStructureList() {
		bodies = append(bodies, column.GenerateCreateQuery())
	}
	for _, index := range ts.IndexStructureList() {
		bodies = append(bodies, index.GenerateCreateQuery())
	}
	migration.Up += " " + strings.Join(bodies, ",\n ") + "\n"

	migration.Up += ")"
	migration.Up += " ENGINE=" + ts.Engine()
	migration.Up += " DEFAULT CHARASET=" + ts.DefaultCharset()
	migration.Up += " COLLATE=" + ts.Collate()
	migration.Up += " COMMENT=" + ts.Comment()

	if ts.Partition() != nil {
		migration.Up += fmt.Sprintf("\n/*!50100 %s */", ts.Partition().Query())
	}

	migration.Down = "DROP TABLE " + ts.Table()

	return migration
}

// TableDropMigration DROP TABLE
type TableDropMigration struct {
	tableMigration
}

// NewTableDropMigration create TableDropMigration
func NewTableDropMigration(ts TableStructure) *TableDropMigration {
	migration := &TableDropMigration{}
	migration.Table = ts.Table()
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
	migration.Table = after.Name()
	migration.Type = CreateOrReplaceType
	migration.IsAltered = before.CompareQuery() != after.CompareQuery()
	if !migration.IsAltered {
		return migration
	}

	migration.Up = strings.Replace(after.CompareQuery(), "CREATE", "CREATE OR REPLACE", 0)
	migration.Down = strings.Replace(before.CompareQuery(), "CREATE", "CREATE OR REPLACE", 0)

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

	migration.Down = strings.Replace(migration.Down, before.Name(), after.Name(), 0)

	return migration
}

// ViewCreateMigration CREATE VIEW
type ViewCreateMigration struct {
	tableMigration
}

// NewViewCreateMigration create ViewCreateMigration
func NewViewCreateMigration(view ViewStructure) *ViewCreateMigration {
	migration := &ViewCreateMigration{}
	migration.Table = view.Name()
	migration.Type = ViewCreateType
	migration.Up = view.CreateQuery()
	migration.Down = "DROP VIEW " + view.Name() + ";"
	return migration
}

// ViewDropMigration DROP VIEW
type ViewDropMigration struct {
	tableMigration
}

// NewViewDropMigration return ViewDropMigration
func NewViewDropMigration(view ViewStructure) *ViewDropMigration {
	migration := &ViewDropMigration{}
	migration.Table = view.Name()
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
	migration.Table = after.Name()
	migration.Type = ViewRenameType
	migration.Up = "RENAME TABLE " + before.Name() + " TO " + after.Name() + ";"
	migration.Down = "RENAME TABLE " + after.Name() + " TO " + before.Name() + ";"
	return migration
}
