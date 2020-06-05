package migration

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

// List migration list
type List struct {
	list []tableMigration
}

// Add table migration
func (l *List) Add(migration tableMigration) {
	l.list = append(l.list, migration)
}
