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

package migration

import "github.com/somen440/gonv/structure"

// ColumnStructure column
type ColumnStructure interface {
	GenerateCreateQuery() string
	GenerateDropQuery() string
	GenerateBaseQuery() string
}

// ModifiedColumnStructure modified
type ModifiedColumnStructure interface {
	GenerateAddQuery() string
	GetColumn() *structure.MySQL57ColumnStructure
}

// ModifiedColumnStructureSet up down set
type ModifiedColumnStructureSet interface {
	UpStructure() *structure.ModifiedColumnStructure
	DownStructure() *structure.ModifiedColumnStructure
}

// IndexStructure index
type IndexStructure interface {
	GenerateCreateQuery() string
	GenerateAddQuery() string
	GenerateDropQuery() string
}

// TableStructure table
type TableStructure interface {
	GetTable() string
	GetColumnStructureList() map[structure.ColumnField]*structure.MySQL57ColumnStructure
	GetIndexStructureList() map[structure.IndexKey]*structure.IndexStructure
	GetEngine() string
	GetDefaultCharset() string
	GetCollate() string
	GetComment() string
	GetPartition() structure.PartitionStructure
}

// ViewStructure view
type ViewStructure interface {
	GetName() string
	CompareQuery() string
	CreateQueryToFormat() string
}

// PartitionStructure partition
type PartitionStructure interface {
	Query() string
}
