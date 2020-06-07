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
		return nil, err
	}
	createTable := f.gdo.ShowCreateTable(tableName)
	defaultCharaset := SearchDefaultCharaset(createTable.schema)

	partition, err := f.createPartitionStructure(dbName, tableName)
	if err != nil {
		return nil, err
	}

	return &structure.TableStructure{
		Table:          tableName,
		Comment:        tableStatus.Comment,
		Engine:         tableStatus.Engine,
		Collate:        tableStatus.Collation,
		DefaultCharset: defaultCharaset,
		Partition:      partition,
	}, nil
}

func (f *Factory) createPartitionStructure(dbName, tableName string) (structure.PartitionStructure, error) {
	partitions, err := f.gdo.SelectPartitions(dbName, tableName)
	if err != nil {
		return nil, err
	}
	if len(partitions) == 0 {
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

	if len(partitions) > 1 && partitions[0].PartitionMethod != "NULL" {
		var methodMap map[PartitionMethod][]SelectPartitionsResult
		for _, partition := range partitions {
			method := PartitionMethod(partition.PartitionMethod)
			methodMap[method] = append(methodMap[method], partition)
		}
		for method, list := range methodMap {
			var rowsMap PartitionRowsMap
			for _, partition := range list {
				var rows partitionRows
				ordinal := PartitionOrdinalPosition(partition.PartitionOrdinalPosition)
				rows[ordinal] = PartitionSummary{
					Name:        partition.PartitionName,
					Description: partition.PartitionDescription,
					Comment:     partition.PartitionComment,
				}

				expression := PartitionExpression(partition.PartitionExpression)
				rowsMap[expression] = rows
			}
			partitionRowsMapGroup[method] = rowsMap
		}
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
	return structure.MySQL57ColumnStructure{
		Field: structure.ColumnField(column.ColumnName),
		Type:  strings.Replace(column.ColumnType, " unsigned", "", -1),
	}, nil
}
