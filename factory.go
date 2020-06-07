package main

import (
	"fmt"
	"strings"

	"github.com/somen440/gonv/structure"
)

// Factory structure factory
type Factory struct {
	gdo *GDO
}

// CreateTableStructure create table structure
func (f *Factory) CreateTableStructure(dbName, tableName string) (*structure.TableStructure, error) {
	if err := f.gdo.SwitchDb(dbName); err != nil {
		return nil, fmt.Errorf("SwitchDb error: %w", err)
	}
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
	var partitionRowsMapGroup PartitionRowsMapGroup

	methodMap := map[PartitionMethod][]SelectPartitionsResult{}
	for _, partition := range partitions {
		method := PartitionMethod(partition.PartitionMethod.String)
		methodMap[method] = append(methodMap[method], partition)
	}
	for method, list := range methodMap {
		var rowsMap PartitionRowsMap
		for _, partition := range list {
			var rows partitionRows
			ordinal := PartitionOrdinalPosition(partition.PartitionOrdinalPosition.Int32)
			rows[ordinal] = PartitionSummary{
				Name:        partition.PartitionName.String,
				Description: partition.PartitionDescription.String,
				Comment:     partition.PartitionComment,
			}

			expression := PartitionExpression(partition.PartitionExpression.String)
			rowsMap[expression] = rows
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
			for value, raws := range group {
				var partMap map[int]structure.PartitionPartStructure
				for order, summary := range raws {
					partMap[int(order)] = structure.PartitionPartStructure{
						Name:     summary.Name,
						Operator: structure.PartitionMethodOperatorMap[m],
						Value:    summary.Description,
						Comment:  summary.Comment,
					}
				}
				partition = &structure.PartitionLongStructure{
					Type:    m,
					Value:   string(value),
					PartMap: partMap,
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
