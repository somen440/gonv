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
)

// TableStructureMap table name structure map
type TableStructureMap map[TableName]*TableStructure

// TableStructureTypeMap table structure type map
type TableStructureTypeMap map[TableStructureType]bool

// DatabaseStructure db
type DatabaseStructure struct {
	Map map[TableName]*TableStructure
}

func (ds *DatabaseStructure) String() string {
	var out bytes.Buffer

	for table, ts := range ds.Map {
		out.WriteString("*************************** table: " + string(table) + " ***************************" + "\n")
		out.WriteString(ts.String() + "\n")
	}

	return out.String()
}

// Add add table structure
func (ds *DatabaseStructure) Add(ts *TableStructure) {
	ds.Map[TableName(ts.Table)] = ts
}

// ListToFilter filter table type
func (ds *DatabaseStructure) ListToFilter(filters TableStructureTypeMap) TableStructureMap {
	m := TableStructureMap{}

	for table, structure := range ds.Map {
		_, ok := filters[structure.Type]
		if ok {
			m[table] = structure
		}
	}

	return m
}

// ListToFilterTableType return table type structure
func (ds *DatabaseStructure) ListToFilterTableType() TableStructureMap {
	filter := TableStructureTypeMap{
		TableType: true,
	}
	return ds.ListToFilter(filter)
}

// ListToFilterViewType return view type structure
func (ds *DatabaseStructure) ListToFilterViewType() TableStructureMap {
	filter := TableStructureTypeMap{
		ViewType: true,
	}
	return ds.ListToFilter(filter)
}

// DiffListToFilterTableType return diff list filter table
func (ds *DatabaseStructure) DiffListToFilterTableType(target *DatabaseStructure) TableStructureMap {
	m := TableStructureMap{}

	before := ds.ListToFilterTableType()
	after := target.ListToFilterTableType()

	for table, structure := range before {
		_, ok := after[table]
		if !ok {
			m[table] = structure
		}
	}

	return m
}

// DiffListToFilterView return diff list filter view
func (ds *DatabaseStructure) DiffListToFilterView(target *DatabaseStructure) TableStructureMap {
	m := TableStructureMap{}

	filter := TableStructureTypeMap{
		ViewType:    true,
		ViewRawType: true,
	}
	before := ds.ListToFilter(filter)
	after := target.ListToFilter(filter)

	for table, structure := range before {
		_, ok := after[table]
		if !ok {
			m[table] = structure
		}
	}

	return m
}
