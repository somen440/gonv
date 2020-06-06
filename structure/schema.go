package structure

// // SchemaKey schema key
// type SchemaKey string

// // SchemaKeys
// const (
// 	TableType            SchemaKey = "table"
// 	TableComment         SchemaKey = "comment"
// 	TableColumn          SchemaKey = "column"
// 	TablePrimaryKey      SchemaKey = "primary_key"
// 	TableIndex           SchemaKey = "index"
// 	TableEngine          SchemaKey = "engine"
// 	TableDefaultCharaset SchemaKey = "default_charaset"
// 	TableCollate         SchemaKey = "collate"
// 	TablePartition       SchemaKey = "partition"

// 	ColumnType      SchemaKey = "type"
// 	ColumnDefault   SchemaKey = "default"
// 	ColumnComment   SchemaKey = "comment"
// 	ColumnAttribute SchemaKey = "attribute"

// 	IndexType   SchemaKey = "is_unique"
// 	IndexColumn SchemaKey = "column"

// 	PartitionBy          SchemaKey = "by"
// 	PartitionValue       SchemaKey = "value"
// 	PartitionList        SchemaKey = "list"
// 	PartitionLessThan    SchemaKey = "less_than"
// 	PartitionIn          SchemaKey = "in"
// 	PartitionEngine      SchemaKey = "engine"
// 	PartitionPartComment SchemaKey = "comment"
// 	PartitionNum         SchemaKey = "num"

// 	ViewAlgorithm SchemaKey = "algorithm"
// 	ViewAlias     SchemaKey = "alias"
// 	ViewColumn    SchemaKey = "column"
// 	ViewFrom      SchemaKey = "from"

// 	JoinReference SchemaKey = "reference"
// 	JoinJoins     SchemaKey = "joins"
// 	JoinFactor    SchemaKey = "factor"
// 	JoinOn        SchemaKey = "on"

// 	ViewRawQuery SchemaKey = "query"
// )

// // SchemaKeySet schema key set
// type SchemaKeySet map[SchemaKey]interface{}

// // SchemaKeys keys
// var (
// 	TableKeys = [9]SchemaKey{
// 		TableType,
// 		TableComment,
// 		TableColumn,
// 		TablePrimaryKey,
// 		TableIndex,
// 		TableEngine,
// 		TableDefaultCharaset,
// 		TableCollate,
// 		TablePartition,
// 	}

// 	TableRequireKeys = [2]SchemaKey{
// 		TableType,
// 		TableColumn,
// 	}

// 	TableOptionalKeys = [7]SchemaKey{
// 		TableComment,
// 		TablePrimaryKey,
// 		TableIndex,
// 		TableEngine,
// 		TableDefaultCharaset,
// 		TableCollate,
// 		TablePartition,
// 	}

// 	ColumnKeys = [4]SchemaKey{
// 		ColumnType,
// 		ColumnDefault,
// 		ColumnComment,
// 		ColumnAttribute,
// 	}

// 	ColumnRequireKeys = [1]SchemaKey{
// 		ColumnType,
// 	}

// 	ColumnOptionalKeys = [3]SchemaKey{
// 		ColumnDefault,
// 		ColumnComment,
// 		ColumnAttribute,
// 	}

// 	IndexRequireKeys = [2]SchemaKey{
// 		IndexType,
// 		IndexColumn,
// 	}

// 	ViewKeys = [3]SchemaKey{
// 		ViewColumn,
// 		ViewFrom,
// 		ViewAlias,
// 	}

// 	ViewRequireKeys = [2]SchemaKey{
// 		ViewColumn,
// 		ViewFrom,
// 	}

// 	ViewOptionalKeys = [1]SchemaKey{
// 		ViewAlias,
// 	}
// )
