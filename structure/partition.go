package structure

import (
	"fmt"
	"strconv"
	"strings"
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
	query = "PARTITION BY " + ps.Type + "(" + ps.Value + ") (\n"
	body := []string{}
	for _, part := range ps.Parts {
		body = append(body, part.Query())
	}
	query += " " + strings.Join(body, ",\n )")
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
