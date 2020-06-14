/*
Copyright 2020 somen440

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package structure

import (
	"bytes"
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
		"PARTITION %s VALUES %s (%s) ENGINE = InnoDB",
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
	PartitionStructure

	Type  string
	Value string
	Parts []*PartitionPartStructure
}

func (ps *PartitionLongStructure) String() string {
	var out bytes.Buffer

	out.WriteString("type: " + ps.Type + "\n")
	out.WriteString("value: " + ps.Value + "\n")
	out.WriteString("parts:\n")
	for _, part := range ps.Parts {
		out.WriteString(part.Query())
	}

	return out.String()
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

// PartitionShortStructure partition shart
type PartitionShortStructure struct {
	PartitionStructure

	Type  string
	Value string
	Num   int
}

// Query return qury
func (ps *PartitionShortStructure) Query() (query string) {
	query = "PARTITION BY " + ps.Type + "(" + ps.Value + ")\n"
	query += "PARTITIONS " + strconv.Itoa(ps.Num)
	return
}
