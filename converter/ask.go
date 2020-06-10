package converter

import "github.com/somen440/gonv/structure"

// TableAsk ask
type TableAsk struct {
	DroppedTableList []structure.TableName
	RenamedTableList map[structure.TableName]structure.TableName

	DroppedColumnList []structure.ColumnField
	RenamedColumnList map[structure.ColumnField]structure.ColumnField
}

// ViewAsk ask
type ViewAsk struct {
	// todo: []ViewName ... etc
}

// ModifiedAsk ask
type ModifiedAsk struct {
	Table *TableAsk
	View  *ViewAsk
}

// NewModifiedAsk return ModifiedAsk
func NewModifiedAsk() *ModifiedAsk {
	return &ModifiedAsk{}
}
