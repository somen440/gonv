package structure

import (
	"bytes"
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
	Order                int
}

// String return string
func (mc *MySQL57ColumnStructure) String() string {
	var out bytes.Buffer

	out.WriteString("\tfield: " + string(mc.Field) + "\n")
	out.WriteString("\t\ttype: " + string(mc.Type) + "\n")
	out.WriteString("\t\tdefault: " + string(mc.Default) + "\n")
	out.WriteString("\t\tcomment: " + string(mc.Comment) + "\n")

	out.WriteString("\t\tattributes:" + "\n")
	for _, attribute := range mc.Attributes {
		out.WriteString("\t\t\t- " + string(attribute) + "\n")
	}

	out.WriteString("\t\tcollation_name: " + string(mc.CollationName) + "\n")

	out.WriteString("\t\tproperties:" + "\n")
	for _, property := range mc.Properties {
		out.WriteString("\t\t\t- " + property + "\n")
	}

	out.WriteString("\t\tgeneration_expression: " + string(mc.GenerationExpression) + "\n")

	out.WriteString("\t\tgenerate: " + mc.GenerateCreateQuery())

	return out.String()
}

// GenerateCreateQuery return creqte query
func (mc *MySQL57ColumnStructure) GenerateCreateQuery() string {
	return "`" + string(mc.Field) + "`" + " " + mc.GenerateBaseQuery()
}

// GenerateBaseQuery return query base
func (mc *MySQL57ColumnStructure) GenerateBaseQuery() string {
	query := []string{mc.Type}
	if mc.IsUnsigned() {
		query = append(query, "unsigned")
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
		if mc.IsOnUpdateCurrentTimestamp() {
			query = append(query, "ON UPDATE CURRENT_TIMESTAMP")
		}
	}
	if mc.Comment != "" {
		query = append(query, "COMMENT", "'"+mc.Comment+"'")
	}
	return strings.Join(query, " ")
}

// GenerateDropQuery return drop query
func (mc *MySQL57ColumnStructure) GenerateDropQuery() string {
	return "DROP COLUMN `" + string(mc.Field) + "`"
}

// DefaultNecessaryQuot return creqte query
func (mc *MySQL57ColumnStructure) DefaultNecessaryQuot() string {
	_, err := strconv.Atoi(mc.Default)
	if err == nil {
		return mc.Default
	}

	t := mc.Default
	if strings.ToUpper(mc.Type) == "DATETIME" && t == "CURRENT_TIMESTAMP" {
		return t
	}
	return "'" + t + "'"
}

// Diff is not match return true
func (mc *MySQL57ColumnStructure) Diff(target *MySQL57ColumnStructure) ([]string, bool) {
	results := []string{}

	if mc.Type != target.Type {
		results = append(results, "Type")
	}
	if mc.Comment != target.Comment {
		results = append(results, "Comment")
	}
	if mc.IsNullable() != target.IsNullable() {
		results = append(results, "IsNullable")
	}
	if mc.IsUnsigned() != target.IsUnsigned() {
		results = append(results, "IsUnsigned")
	}
	if mc.Default != target.Default {
		results = append(results, "Default")
	}
	if mc.IsAutoIncrement() != target.IsAutoIncrement() {
		results = append(results, "IsAutoIncrement")
	}
	if mc.CollationName != target.CollationName {
		results = append(results, "CollationName")
	}
	if mc.IsStored() != target.IsStored() {
		results = append(results, "IsStored")
	}

	return results, len(results) > 0
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

// IsOnUpdateCurrentTimestamp has OnUpdateCurrentTimestamp true
func (mc *MySQL57ColumnStructure) IsOnUpdateCurrentTimestamp() bool {
	for _, attribete := range mc.Attributes {
		if attribete == OnUpdateCurrentTimestamp {
			return true
		}
	}
	return false
}

// IsForceNull return null true
func (mc *MySQL57ColumnStructure) IsForceNull() bool {
	return strings.ToUpper(mc.Type) == "DATETIME"
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
	Column        *MySQL57ColumnStructure
	ModifiedAfter string
}

// GenerateAddQuery return add query
func (ms *ModifiedColumnStructure) GenerateAddQuery() (query string) {
	query = "ADD COLUMN `" + string(ms.Column.Field) + "` "
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

// GetColumn implements migration ModifiedColumnStructure
func (ms *ModifiedColumnStructure) GetColumn() *MySQL57ColumnStructure {
	return ms.Column
}

// ModifiedColumnStructureSet up down set
type ModifiedColumnStructureSet struct {
	Up   *ModifiedColumnStructure
	Down *ModifiedColumnStructure
}

// UpStructure return up
func (ms *ModifiedColumnStructureSet) UpStructure() *ModifiedColumnStructure {
	return ms.Up
}

// DownStructure return down
func (ms *ModifiedColumnStructureSet) DownStructure() *ModifiedColumnStructure {
	return ms.Down
}

// DroppedColumnFielList drop column list
type DroppedColumnFielList []ColumnField

// RenamedColumnFielMap renamed column [before->after]
type RenamedColumnFielMap map[ColumnField]ColumnField

// AddedFieldList added column list
type AddedFieldList []ColumnField
