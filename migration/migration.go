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
	if migration == nil {
		return
	}
	l.list = append(l.list, migration)
}

// Merge migrations
func (l *List) Merge(targetsList ...*List) {
	for _, targets := range targetsList {
		if targets == nil || len(targets.list) == 0 {
			continue
		}
		l.list = append(l.list, targets.list...)
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
