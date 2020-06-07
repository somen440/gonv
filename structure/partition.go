package structure

import (
	"fmt"
	"strconv"
	"strings"
)

// PartitionType short or long
type PartitionType int

// PartitionTypes
const (
	PartitionTypeShort PartitionType = iota
	PartitionTypeLong
)

// PartitionMethod method
type PartitionMethod string

// PartitionMethods
const (
	PartitionMethodKey PartitionMethod = "key"
	LinearKey          PartitionMethod = "linear_key"
	Hash               PartitionMethod = "hash"
	LinearHash         PartitionMethod = "linear_hash"
	List               PartitionMethod = "list"
	ListColumns        PartitionMethod = "list_columns"
	Range              PartitionMethod = "range"
	RangeColumns       PartitionMethod = "range_columns"
)

// Partition maps
var (
	PartitionMethodMap = map[string]PartitionMethod{
		"KEY":           PartitionMethodKey,
		"LINEAR KEY":    LinearKey,
		"HASH":          Hash,
		"LINEAR HASH":   LinearHash,
		"LIST":          List,
		"LIST COLUMNS":  ListColumns,
		"RANGE":         Range,
		"RANGE COLUMNS": RangeColumns,
	}

	PartitionMethodTypeMap = map[string]PartitionType{
		"KEY":           PartitionTypeShort,
		"LINEAR KEY":    PartitionTypeShort,
		"HASH":          PartitionTypeShort,
		"LINEAR HASH":   PartitionTypeShort,
		"LIST":          PartitionTypeLong,
		"LIST COLUMNS":  PartitionTypeLong,
		"RANGE":         PartitionTypeLong,
		"RANGE COLUMNS": PartitionTypeLong,
	}

	PartitionMethodOperatorMap = map[string]string{
		"LIST":          "IN",
		"LIST COLUMNS":  "IN",
		"RANGE":         "LESS THAN",
		"RANGE COLUMNS": "LESS THAN",
	}
)

// PartitionStructure interface
type PartitionStructure interface {
	Query() string
}

// PartitionPartStructure partition part
type PartitionPartStructure struct {
	Name     string
	Operator string
	Value    string
	Comment  string
}

// Query return query
func (ps *PartitionPartStructure) Query() (query string) {
	query = fmt.Sprintf(
		"PARTITION %s VALUES %s (%s)",
		ps.Name,
		strings.ToUpper(ps.Operator),
		ps.Value,
	)
	if ps.Comment != "" {
		query += " COMMENT = " + ps.Comment
	}
	return
}

// PartitionLongStructure partition long
type PartitionLongStructure struct {
	Type  string
	Value string
	Parts []PartitionPartStructure
}

// Query return query
func (ps *PartitionLongStructure) Query() (query string) {
	query = "PARTITION BY " + ps.Type + "(" + ps.Value + ")\n("
	body := []string{}
	for _, part := range ps.Parts {
		body = append(body, part.Query())
	}
	query += strings.Join(body, ",\n ") + ")"
	return
}

// PartitionShartStructure partition shart
type PartitionShartStructure struct {
	Type  string
	Value string
	Num   int
}

// Query return qury
func (ps *PartitionShartStructure) Query() (query string) {
	query = "PARTITION BY " + ps.Type + "(" + ps.Value + ")\n"
	query += "PARTITIONS " + strconv.Itoa(ps.Num)
	return
}
