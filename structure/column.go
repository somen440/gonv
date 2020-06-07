package structure

import (
	"strconv"
	"strings"
)

// ColumnField column field
type ColumnField string

// MySQL57ColumnStructure mysql 5.7
type MySQL57ColumnStructure struct {
	Field                ColumnField
	Type                 string
	Default              string
	Comment              string
	Attributes           []Attribute
	CollationName        string
	Properties           []string
	GenerationExpression string
}

// GenerateCreateQuery return creqte query
func (mc *MySQL57ColumnStructure) GenerateCreateQuery() string {
	return string(mc.Field) + " " + mc.GenerateBaseQuery()
}

// GenerateBaseQuery return query base
func (mc *MySQL57ColumnStructure) GenerateBaseQuery() string {
	query := []string{mc.Type}
	if mc.IsUnsigned() {
		query = append(query, mc.Type)
	}
	if mc.CollationName != "" {
		query = append(query, "COLLATE", mc.CollationName)
	}
	if mc.GenerationExpression != "" {
		query = append(query, "AS", "("+mc.GenerationExpression+")")
	}
	if mc.IsStored() {
		query = append(query, "STORED")
	}
	if !mc.IsNullable() {
		query = append(query, "NOT NULL")
	} else if mc.IsForceNull() {
		query = append(query, "NULL")
	}
	if mc.GenerationExpression == "" {
		if mc.IsAutoIncrement() {
			query = append(query, "AUTO_INCREMENT")
		} else if mc.Default != "" {
			query = append(query, "DEFAULT", mc.DefaultNecessaryQuot())
		} else if mc.IsNullable() {
			query = append(query, "DEFAULT NULL")
		}
	}
	query = append(query, "COMMENT", mc.Comment)
	return strings.Join(query, " ")
}

// GenerateDropQuery return drop query
func (mc *MySQL57ColumnStructure) GenerateDropQuery() string {
	return "DROP COLUMN " + string(mc.Field)
}

// DefaultNecessaryQuot return creqte query
func (mc *MySQL57ColumnStructure) DefaultNecessaryQuot() string {
	_, err := strconv.Atoi(mc.Default)
	if err == nil {
		return mc.Default
	}

	t := mc.Default
	if mc.IsForceNull() && t == "CURRENT_TIMESTAMP" {
		return t
	}
	return "'" + t + "'"
}

// IsChanged is not match return true
func (mc *MySQL57ColumnStructure) IsChanged(target MySQL57ColumnStructure) bool {
	return !(mc.Type == target.Type &&
		mc.Comment == target.Comment &&
		mc.IsNullable() == target.IsNullable() &&
		mc.IsUnsigned() == target.IsUnsigned() &&
		mc.Default == target.Default &&
		mc.IsAutoIncrement() == target.IsAutoIncrement() &&
		mc.CollationName == target.CollationName &&
		mc.IsStored() == target.IsStored())
}

// IsUnsigned has unsigned true
func (mc *MySQL57ColumnStructure) IsUnsigned() bool {
	for _, attribete := range mc.Attributes {
		if attribete == Unsigned {
			return true
		}
	}
	return false
}

// IsForceNull return null true
func (mc *MySQL57ColumnStructure) IsForceNull() bool {
	return strings.ToUpper(mc.Type) == "TIMESTAMP"
}

// IsNullable has nullable true
func (mc *MySQL57ColumnStructure) IsNullable() bool {
	for _, attribete := range mc.Attributes {
		if attribete == Nullable {
			return true
		}
	}
	return false
}

// IsAutoIncrement has nullable true
func (mc *MySQL57ColumnStructure) IsAutoIncrement() bool {
	for _, attribete := range mc.Attributes {
		if attribete == AutoIncrement {
			return true
		}
	}
	return false
}

// IsStored has nullable true
func (mc *MySQL57ColumnStructure) IsStored() bool {
	for _, attribete := range mc.Attributes {
		if attribete == Stored {
			return true
		}
	}
	return false
}

// ModifiedColumnStructure modified column
type ModifiedColumnStructure struct {
	BeforeField   ColumnField
	Column        MySQL57ColumnStructure
	ModifiedAfter string
}

// GenerateAddQuery return add query
func (ms *ModifiedColumnStructure) GenerateAddQuery() (query string) {
	query = "ADD COLUMN " + string(ms.Column.Field)
	query += ms.Column.GenerateBaseQuery()
	if ms.IsOrderChanged() {
		query += " " + ms.ModifiedAfter
	}
	return
}

// GenerateChangeQuery return change query
func (ms *ModifiedColumnStructure) GenerateChangeQuery() (query string) {
	query = "CHANGE " + string(ms.BeforeField) + " " + string(ms.Column.Field) + " "
	query += ms.Column.GenerateBaseQuery()
	if ms.IsOrderChanged() {
		query += " " + ms.ModifiedAfter
	}
	return
}

// IsRenamed return is not match field return true
func (ms *ModifiedColumnStructure) IsRenamed() bool {
	return ms.BeforeField != ms.Column.Field
}

// IsOrderChanged return is set modifier return true.
func (ms *ModifiedColumnStructure) IsOrderChanged() bool {
	return ms.ModifiedAfter != ""
}

// SetModifiedAfter set modifier after
func (ms *ModifiedColumnStructure) SetModifiedAfter(modifier string) {
	if modifier == "" {
		ms.ModifiedAfter = "FIRST"
	} else {
		ms.ModifiedAfter = "AFTER " + modifier
	}
}

// ModifiedColumnStructureSet up down set
type ModifiedColumnStructureSet struct {
	Up   ModifiedColumnStructure
	Down ModifiedColumnStructure
}

// DroppedColumnFielList drop column list
type DroppedColumnFielList []ColumnField

// RenamedColumnFielMap renamed column [before->after]
type RenamedColumnFielMap map[ColumnField]ColumnField

// AddedFieldList added column list
type AddedFieldList []ColumnField
