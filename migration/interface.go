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
	Column() ColumnStructure
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
	Table() string
	ColumnStructureList() []ColumnStructure
	IndexStructureList() []IndexStructure
	Engine() string
	DefaultCharset() string
	Collate() string
	Comment() string
	Partition() PartitionStructure
}

// ViewStructure view
type ViewStructure interface {
	Name() string
	CompareQuery() string
	CreateQuery() string
}
