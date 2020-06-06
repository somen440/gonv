package structure

// TableStructureMap table name structure map
type TableStructureMap map[TableName]TableStructure

// TableStructureTypeMap table structure type map
type TableStructureTypeMap map[TableStructureType]bool

// DatabaseStructure db
type DatabaseStructure struct {
	Map map[TableName]TableStructure
}

// ListToFilter filter table type
func (ds *DatabaseStructure) ListToFilter(filters TableStructureTypeMap) (m TableStructureMap) {
	for table, structure := range ds.Map {
		_, ok := filters[structure.Type]
		if ok {
			m[table] = structure
		}
	}
	return
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
func (ds *DatabaseStructure) DiffListToFilterTableType(target *DatabaseStructure) (m TableStructureMap) {
	before := ds.ListToFilterTableType()
	after := target.ListToFilterTableType()

	for table, structure := range before {
		_, ok := after[table]
		if !ok {
			m[table] = structure
		}
	}
	return
}

// DiffListToFilterView return diff list filter view
func (ds *DatabaseStructure) DiffListToFilterView(target *DatabaseStructure) (m TableStructureMap) {
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
	return
}
