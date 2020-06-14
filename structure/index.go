package structure

import (
	"fmt"
	"reflect"
	"strings"
)

// IndexKey index key
type IndexKey string

// IndexStructure index
type IndexStructure struct {
	Key            IndexKey
	IndexType      string
	ColumnNameList []string
	IsPrimary      bool
	IsUnique       bool
	IsBtree        bool
	Order          int
}

// NewIndexStructure return IndexStructure
func NewIndexStructure(keyName IndexKey, indexType string, isUnique bool, columnNameList []string, order int) *IndexStructure {
	return &IndexStructure{
		Key:            IndexKey(keyName),
		IndexType:      indexType,
		ColumnNameList: columnNameList,
		IsUnique:       isUnique,
		IsPrimary:      keyName == "PRIMARY",
		IsBtree:        strings.ToUpper(indexType) == "BTREE",
		Order:          order,
	}
}

// GenerateCreateQuery return create query
func (is *IndexStructure) GenerateCreateQuery() (query string) {
	if is.IsPrimary {
		query += "PRIMARY KEY "
	} else {
		if is.IsUnique {
			query += "UNIQUE "
		} else if !is.IsBtree {
			query += strings.ToUpper(is.IndexType) + " "
		}
		query += "KEY "
		query += "`" + string(is.Key) + "` "
	}
	query += is.GenerateIndexText()
	return
}

// GenerateAddQuery return add query
func (is *IndexStructure) GenerateAddQuery() (query string) {
	query += "ADD "
	if is.IsPrimary {
		query += "PRIMARY KEY "
	} else {
		if is.IsUnique {
			query += "UNIQUE "
		} else if !is.IsBtree {
			query += strings.ToUpper(is.IndexType) + " "
		} else {
			query += "INDEX "
		}
		query += "`" + string(is.Key) + "` "
	}
	query += is.GenerateIndexText()
	return
}

// GenerateDropQuery return drop query
func (is *IndexStructure) GenerateDropQuery() string {
	if is.IsPrimary {
		return "DROP PRIMARY KEY"
	}
	return "DROP INDEX `" + string(is.Key) + "`"
}

// GenerateIndexText return index text
func (is *IndexStructure) GenerateIndexText() (text string) {
	columns := func() (result []string) {
		for _, name := range is.ColumnNameList {
			result = append(result, fmt.Sprintf("`%s`", name))
		}
		return
	}()
	return fmt.Sprintf("(%s)", strings.Join(columns, ", "))
}

// IsChanged is not match target return true
func (is *IndexStructure) IsChanged(target *IndexStructure) bool {
	return !(reflect.DeepEqual(
		is.ColumnNameList,
		target.ColumnNameList,
	) && is.IsUnique == target.IsUnique &&
		is.IsBtree == target.IsBtree)
}
