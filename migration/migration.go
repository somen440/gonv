package migration

import (
	"bytes"
)

// Type t
type Type int

// Migration Types
const (
	CreateType Type = iota
	AlterType
	DropType
	CreateOrReplaceType
	ViewCreateType
	ViewDropType
	ViewRenameType
)

// Migration interface
type Migration interface {
	UpQuery() string
	DownQuery() string
}

type tableMigration struct {
	Migration

	Table string
	Type  Type
	Up    string
	Down  string
}

func (m *tableMigration) UpQuery() string {
	return m.Up
}

func (m *tableMigration) DownQuery() string {
	return m.Down
}

// Line migration line
type Line interface {
	UpList() []string
	DownList() []string
}

// List migration list
type List struct {
	list []Migration
}

// Add table migration
func (l *List) Add(migration Migration) {
	l.list = append(l.list, migration)
}

// Merge migrations
func (l *List) Merge(targetsList ...*List) {
	for _, targets := range targetsList {
		for _, migration := range targets.list {
			l.list = append(l.list, migration)
		}
	}
}

// List list
func (l *List) List() []Migration {
	return l.list
}

// String to string
func (l *List) String() string {
	var out bytes.Buffer

	if len(l.list) == 0 {
		return "no migrations."
	}

	out.WriteString("*************************** migration up ***************************" + "\n")
	for _, migration := range l.list {
		out.WriteString(migration.UpQuery() + "\n")
	}

	out.WriteString("*************************** migration down ***************************" + "\n")
	for _, migration := range l.list {
		out.WriteString(migration.DownQuery() + "\n")
	}

	return out.String()
}
