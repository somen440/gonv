package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/somen440/gonv/structure"
)

// Factory structure factory
type Factory struct {
	gdo *GDO
}

// CreateDatabaseStructureFromSchema create database structure
func (f *Factory) CreateDatabaseStructureFromSchema(dbName, schema string) (*structure.DatabaseStructure, error) {
	if f, err := os.Stat(schema); !(!os.IsNotExist(err) && f.IsDir()) {
		return nil, fmt.Errorf("%s id not dir: %w", schema, err)
	}

	if err := f.gdo.CreateDatabase(dbName); err != nil {
		return nil, err
	}
	if err := f.gdo.SwitchDb(dbName); err != nil {
		return nil, fmt.Errorf("SwitchDb error: %w", err)
	}
	defer func() {
		err := f.gdo.DropDatabaseIfExists(dbName)
		if err != nil {
			panic(err)
		}
	}()

	files, err := ioutil.ReadDir(schema)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		path := filepath.Join(schema, file.Name())
		contains, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		f.gdo.ExecSchema(string(contains))
	}

	return f.CreateDatabaseStructure(dbName)
}

// CreateDatabaseStructure create database structure
func (f *Factory) CreateDatabaseStructure(dbName string) (*structure.DatabaseStructure, error) {
	if err := f.gdo.SwitchDb(dbName); err != nil {
		return nil, fmt.Errorf("SwitchDb error: %w", err)
	}

	result := &structure.DatabaseStructure{
		Map: map[structure.TableName]*structure.TableStructure{},
	}
	for _, table := range f.gdo.ShowTables() {
		ts, err := f.createTableStructure(dbName, table)
		if err != nil {
			return nil, err
		}
		result.Add(ts)
	}
	return result, nil
}

func (f *Factory) createTableStructure(dbName, tableName string) (*structure.TableStructure, error) {
	tableStatus, err := f.gdo.ShowTableStatusLike(tableName)
	if err != nil {
		return nil, fmt.Errorf("ShowTableStatusLike error: %w", err)
	}
	createTable := f.gdo.ShowCreateTable(tableName)
	defaultCharaset := SearchDefaultCharaset(createTable.schema)

	partition, err := f.createPartitionStructure(dbName, tableName)
	if err != nil {
		return nil, fmt.Errorf("createPartitionStructure error: %w", err)
	}

	columns, err := f.createColumnStructureList(dbName, tableName)
	if err != nil {
		return nil, fmt.Errorf("createColumnStructureList error: %w", err)
	}

	indexes, err := f.createIndexStructureList(tableName)
	if err != nil {
		return nil, fmt.Errorf("createIndexStructureList error: %w", err)
	}

	return &structure.TableStructure{
		Table:               tableName,
		Comment:             tableStatus.Comment,
		Engine:              tableStatus.Engine,
		Collate:             tableStatus.Collation,
		DefaultCharset:      defaultCharaset,
		Partition:           partition,
		ColumnStructureList: columns,
		IndexStructureList:  indexes,
	}, nil
}

func (f *Factory) createPartitionStructure(dbName, tableName string) (structure.PartitionStructure, error) {
	partitions, err := f.gdo.SelectPartitions(dbName, tableName)
	if err != nil {
		return nil, err
	}

	if len(partitions) == 1 && !partitions[0].PartitionMethod.Valid {
		// not exist partition
		return nil, nil
	}

	type PartitionSummary struct {
		Name        string
		Description string
		Comment     string
	}
	type PartitionMethod string
	type PartitionExpression string
	type PartitionOrdinalPosition int
	type partitionRows map[PartitionOrdinalPosition]PartitionSummary
	type PartitionRowsMap map[PartitionExpression]partitionRows
	type PartitionRowsMapGroup map[PartitionMethod]PartitionRowsMap

	partitionRowsMapGroup := PartitionRowsMapGroup{}
	methodMap := map[PartitionMethod][]SelectPartitionsResult{}
	for _, partition := range partitions {
		method := PartitionMethod(partition.PartitionMethod.String)
		methodMap[method] = append(methodMap[method], partition)
	}
	for method, list := range methodMap {
		rowsMap := PartitionRowsMap{}
		for _, partition := range list {
			ordinal := PartitionOrdinalPosition(partition.PartitionOrdinalPosition.Int32)
			expression := PartitionExpression(partition.PartitionExpression.String)

			_, ok := rowsMap[expression]
			if !ok {
				rowsMap[expression] = partitionRows{}
			}

			rowsMap[expression][ordinal] = PartitionSummary{
				Name:        partition.PartitionName.String,
				Description: partition.PartitionDescription.String,
				Comment:     partition.PartitionComment,
			}
		}
		partitionRowsMapGroup[method] = rowsMap
	}

	var partition structure.PartitionStructure
	for method, group := range partitionRowsMapGroup {
		m := string(method)
		switch t := structure.PartitionMethodTypeMap[m]; t {
		case structure.PartitionTypeShort:
			for value, raws := range group {
				partition = &structure.PartitionShartStructure{
					Type:  m,
					Value: string(value),
					Num:   len(raws),
				}
			}
			break
		case structure.PartitionTypeLong:
			for value, rows := range group {
				parts := []structure.PartitionPartStructure{}
				var orders []int
				for order := range rows {
					orders = append(orders, int(order))
				}
				sort.Ints(orders)

				for _, order := range orders {
					summary := rows[PartitionOrdinalPosition(order)]
					parts = append(parts, structure.PartitionPartStructure{
						Name:     summary.Name,
						Operator: structure.PartitionMethodOperatorMap[m],
						Value:    summary.Description,
						Comment:  summary.Comment,
					})
				}

				partition = &structure.PartitionLongStructure{
					Type:  m,
					Value: string(value),
					Parts: parts,
				}
			}
			break
		}
	}

	return partition, nil
}

func (f *Factory) createColumnStructureList(dbName, tableName string) ([]structure.MySQL57ColumnStructure, error) {
	columns, count, err := f.gdo.SelectColumns(dbName, tableName)
	if err != nil {
		return nil, err
	}

	resultsIndex := 0
	results := make([]structure.MySQL57ColumnStructure, count)
	for _, column := range columns {
		result, err := f.createColumnStructure(column)
		if err != nil {
			return nil, err
		}
		results[resultsIndex] = result
		resultsIndex++
	}

	return results, nil
}

func (f *Factory) createColumnStructure(column SelectColumnsResult) (structure.MySQL57ColumnStructure, error) {
	var attributes []structure.Attribute

	if strings.Contains(column.Extra, "auto_increment") {
		attributes = append(attributes, structure.AutoIncrement)
	}
	if column.IsNullable == "YES" {
		attributes = append(attributes, structure.Nullable)
	}
	if strings.Contains(column.ColumnType, "unsigned") {
		attributes = append(attributes, structure.Unsigned)
	}
	if strings.Contains(column.Extra, "STORED") {
		attributes = append(attributes, structure.Stored)
	}

	return structure.MySQL57ColumnStructure{
		Field:         structure.ColumnField(column.ColumnName),
		Type:          TrimUnsigned(column.ColumnType),
		Default:       column.ColumnDefault.String,
		Comment:       column.ColumnComment,
		CollationName: column.CollationName.String,
		Attributes:    attributes,
	}, nil
}

func (f *Factory) createIndexStructureList(tableName string) ([]structure.IndexStructure, error) {
	indexes, err := f.gdo.ShowIndex(tableName)
	if err != nil {
		return nil, err
	}

	type IndexSummary struct {
		IsUnique   bool
		ColumnName string
		IndexType  string
	}
	type IndexSummaryGroup map[string][]IndexSummary

	group := IndexSummaryGroup{}
	for _, index := range indexes {
		columnName := index.ColumnName
		if index.SubPart.Valid {
			columnName = fmt.Sprintf("%s(%s)", columnName, index.SubPart.String)
		}
		group[index.KeyName] = append(group[index.KeyName], IndexSummary{
			IsUnique:   index.NonUnique == 0,
			ColumnName: columnName,
			IndexType:  index.IndexType,
		})
	}

	var results []structure.IndexStructure
	for keyName, list := range group {
		columnNameList := make([]string, len(list))
		for i, index := range list {
			columnNameList[i] = index.ColumnName
		}
		results = append(results, structure.NewIndexStructure(
			keyName,
			list[0].IndexType,
			list[0].IsUnique,
			columnNameList,
		))
	}
	return results, nil
}
