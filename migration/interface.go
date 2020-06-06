package migration

// ColumnStructure column
type ColumnStructure interface {
	GenerateCreateQuery() string
	GenerateDropQuery() string
	GenerateBaseQuery() string
}

// ModifiedColumnStructure modified
type ModifiedColumnStructure interface {
	GenerateAddQuery() string
	GetColumn() ColumnStructure
}

// IndexStructure index
type IndexStructure interface {
	GenerateCreateQuery() string
	GenerateAddQuery() string
	GenerateDropQuery() string
}

// PartitionStructure partition
type PartitionStructure interface {
	Query() string
}

// TableStructure table
type TableStructure interface {
	GetTable() string
	GetColumnStructureList() []ColumnStructure
	GetIndexStructureList() []IndexStructure
	GetEngine() string
	GetDefaultCharset() string
	GetCollate() string
	GetComment() string
	GetPartition() PartitionStructure
}

// ViewStructure view
type ViewStructure interface {
	GetName() string
	GetCompareQuery() string
	GetCreateQuery() string
}
