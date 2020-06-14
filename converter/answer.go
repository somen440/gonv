package converter

import "github.com/somen440/gonv/structure"

// TableAnswer answer table
type TableAnswer struct {
	DroppedTableList []structure.TableName
	RenamedTableList map[structure.TableName]structure.TableName

	DroppedColumnList []structure.ColumnField
	RenamedColumnList map[structure.ColumnField]structure.ColumnField
}

// RenamedColumnListAsStrings to strings
func (a *TableAnswer) RenamedColumnListAsStrings() []string {
	results := []string{}

	for _, table := range a.RenamedColumnList {
		results = append(results, string(table))
	}

	return results
}

// ViewAnswer answer view
type ViewAnswer struct {
	// todo: []ViewName ... etc
}

// ModifiedAnswer answer modofied
type ModifiedAnswer struct {
	Table *TableAnswer
	View  *ViewAnswer
}

// NewModifiedAnswer return ModifiedAnswer
func NewModifiedAnswer(table *TableAnswer, view *ViewAnswer) *ModifiedAnswer {
	return &ModifiedAnswer{
		Table: table,
		View:  view,
	}
}
