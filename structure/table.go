package structure

import (
	"bytes"
	"strings"
)

// TableName table name
type TableName string

// TableStructureType type
type TableStructureType string

// TableStructureTypes
const (
	TableType   TableStructureType = "table"
	ViewType    TableStructureType = "view"
	ViewRawType TableStructureType = "view_raw"
)

// Is match return true
func (t TableStructureType) Is(target TableStructureType) bool {
	return t == target
}

// TableStructure table
type TableStructure struct {
	Table               string
	Type                TableStructureType
	Comment             string
	Engine              string
	DefaultCharset      string
	Collate             string
	ColumnStructureList []MySQL57ColumnStructure
	IndexStructureList  []IndexStructure
	Partition           PartitionStructure
	Properties          []string
}

// String structure tu string
func (ts *TableStructure) String() string {
	var out bytes.Buffer

	out.WriteString("name: " + ts.Table + "\n")
	out.WriteString("type: " + string(ts.Type) + "\n")
	out.WriteString("comment: " + ts.Comment + "\n")
	out.WriteString("engine: " + ts.Engine + "\n")
	out.WriteString("default_charset: " + ts.DefaultCharset + "\n")
	out.WriteString("collate: " + ts.Collate + "\n")
	out.WriteString("properties: " + strings.Join(ts.Properties, ", ") + "\n")

	out.WriteString("columns:\n")
	for _, column := range ts.ColumnStructureList {
		out.WriteString("\t" + column.GenerateCreateQuery() + "\n")
	}

	return out.String()
}

// ColumnStructureMap column map
type ColumnStructureMap map[ColumnField]MySQL57ColumnStructure

// ColumnFieldMap column map
type ColumnFieldMap map[ColumnField]bool

// GetColumnStructureMap return columnmap
func (ts *TableStructure) GetColumnStructureMap() (result ColumnStructureMap) {
	for _, structure := range ts.ColumnStructureList {
		result[structure.Field] = structure
	}
	return
}

// GetOrderColumnStructureMap return order
func (ts *TableStructure) GetOrderColumnStructureMap(diff, renamed ColumnStructureMap) (result ColumnStructureMap) {
	for _, structure := range ts.ColumnStructureList {
		_, ok := diff[structure.Field]
		if ok {
			continue
		}
		t, ok := renamed[structure.Field]
		if ok {
			result[structure.Field] = t
		}
	}
	return
}

// IndexMap index key structure map
type IndexMap map[IndexKey]IndexStructure

// GetIndexMap return index map
func (ts *TableStructure) GetIndexMap() (result IndexMap) {
	for _, structure := range ts.IndexStructureList {
		result[structure.Key] = structure
	}
	return
}

// GetDiffColumnList return diff column
func (ts *TableStructure) GetDiffColumnList(target *TableStructure) (result ColumnStructureMap) {
	targetMap := target.GetColumnStructureMap()
	for field, structure := range ts.GetColumnStructureMap() {
		_, ok := targetMap[field]
		if ok {
			continue
		}
		result[field] = structure
	}
	return
}

// ModifiedColumnStructureSetMap modified column
type ModifiedColumnStructureSetMap map[ColumnField]ModifiedColumnStructureSet

// GenerateModifiedColumnStructureSetMap return a
func (ts *TableStructure) GenerateModifiedColumnStructureSetMap(
	target *TableStructure,
	renamed RenamedColumnFielMap,
) (result ModifiedColumnStructureSetMap) {
	targetMap := target.GetColumnStructureMap()
	for beforeField, before := range ts.GetColumnStructureMap() {
		afterField, ok := renamed[beforeField]
		if ok {
			after := targetMap[afterField]
			result[beforeField] = ModifiedColumnStructureSet{
				Up: ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      after,
				},
				Down: ModifiedColumnStructure{
					BeforeField: afterField,
					Column:      before,
				},
			}
		}
		after := targetMap[beforeField]
		if before.IsChanged(after) {
			result[beforeField] = ModifiedColumnStructureSet{
				Up: ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      after,
				},
				Down: ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      before,
				},
			}
		}
	}
	return
}
