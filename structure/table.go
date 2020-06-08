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

// PartitionStructure interface
type PartitionStructure interface {
	Query() string
}

// TableStructure table
type TableStructure struct {
	Table               string
	Type                TableStructureType
	Comment             string
	Engine              string
	DefaultCharset      string
	Collate             string
	ColumnStructureList []*MySQL57ColumnStructure
	IndexStructureList  []*IndexStructure
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

	if len(ts.IndexStructureList) > 0 {
		out.WriteString("index:\n")
		for _, index := range ts.IndexStructureList {
			out.WriteString("\t" + index.GenerateCreateQuery() + "\n")
		}
	}

	if ts.Partition != nil {
		out.WriteString("partitions:\n")
		out.WriteString(ts.Partition.Query())
	}

	return out.String()
}

// ColumnStructureMap column map
type ColumnStructureMap map[ColumnField]*MySQL57ColumnStructure

// ColumnFieldMap column map
type ColumnFieldMap map[ColumnField]bool

// GetColumnStructureMap return columnmap
func (ts *TableStructure) GetColumnStructureMap() (result ColumnStructureMap) {
	for _, structure := range ts.ColumnStructureList {
		result[structure.Field] = structure
	}
	return
}

// RenamedField before -> after
type RenamedField map[ColumnField]ColumnField

// GetOrderColumnStructureMap return order
func (ts *TableStructure) GetOrderColumnStructureMap(diff ColumnFieldMap, renamed RenamedField) (result RenamedField) {
	for _, structure := range ts.ColumnStructureList {
		before := structure.Field

		_, ok := diff[before]
		if ok {
			continue
		}

		after, ok := renamed[before]
		if ok {
			result[structure.Field] = after
		}
	}
	return
}

// IndexMap index key structure map
type IndexMap map[IndexKey]*IndexStructure

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
type ModifiedColumnStructureSetMap map[ColumnField]*ModifiedColumnStructureSet

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
			result[beforeField] = &ModifiedColumnStructureSet{
				Up: &ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      after,
				},
				Down: &ModifiedColumnStructure{
					BeforeField: afterField,
					Column:      before,
				},
			}
		}
		after := targetMap[beforeField]
		if before.IsChanged(after) {
			result[beforeField] = &ModifiedColumnStructureSet{
				Up: &ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      after,
				},
				Down: &ModifiedColumnStructure{
					BeforeField: beforeField,
					Column:      before,
				},
			}
		}
	}
	return
}

// GetModifiedColumnList modified column
func (ts *TableStructure) GetModifiedColumnList(fieldList []ColumnField) []*ModifiedColumnStructure {
	columns := ts.GetColumnStructureMap()
	orders := ts.GetOrderColumnStructureMap(ColumnFieldMap{}, RenamedField{})
	for _, field := range fieldList {
		modified := ModifiedColumnStructure{
			BeforeField: field,
			Column:      columns[field],
		}
		order, ok := orders[field]
		if ok {
			modified.SetModifiedAfter(string(order))
		}
	}
	return nil
}

// GetTable implements migrations TableStructure
func (ts *TableStructure) GetTable() string {
	return ts.Table
}

// GetColumnStructureList implements migrations TableStructure
func (ts *TableStructure) GetColumnStructureList() []*MySQL57ColumnStructure {
	return ts.ColumnStructureList
}

// GetIndexStructureList implements migrations TableStructure
func (ts *TableStructure) GetIndexStructureList() []*IndexStructure {
	return ts.IndexStructureList
}

// GetEngine implements migrations TableStructure
func (ts *TableStructure) GetEngine() string {
	return ts.Engine
}

// GetDefaultCharset implements migrations TableStructure
func (ts *TableStructure) GetDefaultCharset() string {
	return ts.DefaultCharset
}

// GetCollate implements migrations TableStructure
func (ts *TableStructure) GetCollate() string {
	return ts.Collate
}

// GetComment implements migrations TableStructure
func (ts *TableStructure) GetComment() string {
	return ts.Comment
}

// GetPartition implements migrations TableStructure
func (ts *TableStructure) GetPartition() PartitionStructure {
	return ts.Partition
}
