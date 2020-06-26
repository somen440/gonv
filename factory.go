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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/somen440/gonv/structure"
	"github.com/somen440/gonv/util"
)

// Factory structure factory
type Factory struct {
	gdo *GDO
}

// CreateDatabaseStructureFromSchema create database structure
func (f *Factory) CreateDatabaseStructureFromSchema(dbName, schema string, ignores []string) (*structure.DatabaseStructure, error) {
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
		if err := f.gdo.ExecSchema(string(contains)); err != nil {
			return nil, err
		}
	}

	return f.CreateDatabaseStructure(dbName, ignores)
}

// CreateDatabaseStructure create database structure
func (f *Factory) CreateDatabaseStructure(dbName string, ignores []string) (*structure.DatabaseStructure, error) {
	if err := f.gdo.SwitchDb(dbName); err != nil {
		return nil, fmt.Errorf("SwitchDb error: %w", err)
	}

	result := &structure.DatabaseStructure{
		Map: map[structure.TableName]*structure.TableStructure{},
	}
	for _, table := range f.gdo.ShowTables() {
		isFound := false
		for _, v := range ignores {
			if table == v {
				isFound = true
			}
		}
		if isFound {
			continue
		}
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
	defaultCharaset := util.SearchDefaultCharaset(createTable.schema)

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
		Type:                structure.TableType,
		Comment:             tableStatus.Comment,
		Engine:              tableStatus.Engine,
		Collate:             tableStatus.Collation,
		DefaultCharset:      defaultCharaset,
		Partition:           partition,
		ColumnStructureList: columns,
		IndexStructureList:  indexes,
		Properties:          []string{},
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
				partition = &structure.PartitionShortStructure{
					Type:  m,
					Value: string(value),
					Num:   len(raws),
				}
			}
		case structure.PartitionTypeLong:
			for value, rows := range group {
				parts := []*structure.PartitionPartStructure{}
				var orders []int
				for order := range rows {
					orders = append(orders, int(order))
				}
				sort.Ints(orders)

				for _, order := range orders {
					summary := rows[PartitionOrdinalPosition(order)]
					parts = append(parts, &structure.PartitionPartStructure{
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
		}
	}

	return partition, nil
}

func (f *Factory) createColumnStructureList(dbName, tableName string) (map[structure.ColumnField]*structure.MySQL57ColumnStructure, error) {
	columns, count, err := f.gdo.SelectColumns(dbName, tableName)
	if err != nil {
		return nil, err
	}

	resultsIndex := 0
	results := make(map[structure.ColumnField]*structure.MySQL57ColumnStructure, count)
	for _, column := range columns {
		result, err := f.createColumnStructure(column, resultsIndex)
		if err != nil {
			return nil, err
		}
		results[structure.ColumnField(column.ColumnName)] = result
		resultsIndex++
	}

	return results, nil
}

func (f *Factory) createColumnStructure(column SelectColumnsResult, index int) (*structure.MySQL57ColumnStructure, error) {
	attributes := []structure.Attribute{}

	if strings.Contains(column.Extra, string(structure.AutoIncrement)) {
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
	if strings.Contains(column.Extra, string(structure.OnUpdateCurrentTimestamp)) {
		attributes = append(attributes, structure.OnUpdateCurrentTimestamp)
	}

	return &structure.MySQL57ColumnStructure{
		Field:         structure.ColumnField(column.ColumnName),
		Type:          util.TrimUnsigned(column.ColumnType),
		Default:       column.ColumnDefault.String,
		Comment:       column.ColumnComment,
		CollationName: column.CollationName.String,
		Attributes:    attributes,
		Properties:    []string{},
		Order:         index,
	}, nil
}

func (f *Factory) createIndexStructureList(tableName string) (map[structure.IndexKey]*structure.IndexStructure, error) {
	indexes, err := f.gdo.ShowIndex(tableName)
	if err != nil {
		return nil, err
	}

	type IndexSummary struct {
		IsUnique   bool
		ColumnName string
		IndexType  string
	}
	type IndexSummaryGroup map[structure.IndexKey][]IndexSummary

	indexOrders := []structure.IndexKey{}
	group := IndexSummaryGroup{}
	for _, index := range indexes {
		columnName := index.ColumnName
		if index.SubPart.Valid {
			columnName = fmt.Sprintf("%s(%s)", columnName, index.SubPart.String)
		}
		k := structure.IndexKey(index.KeyName)
		group[k] = append(group[k], IndexSummary{
			IsUnique:   index.NonUnique == 0,
			ColumnName: columnName,
			IndexType:  index.IndexType,
		})
		found := false
		for _, indexOrder := range indexOrders {
			if k == indexOrder {
				found = true
			}
		}
		if found {
			continue
		}
		indexOrders = append(indexOrders, k)
	}

	results := map[structure.IndexKey]*structure.IndexStructure{}
	for keyName, list := range group {
		columnNameList := make([]string, len(list))
		for i, index := range list {
			columnNameList[i] = index.ColumnName
		}
		indexOrder := -1
		for i, order := range indexOrders {
			if keyName == order {
				indexOrder = i
				break
			}
		}
		if indexOrder == -1 {
			return nil, errors.New("not found index")
		}
		results[keyName] = structure.NewIndexStructure(
			keyName,
			list[0].IndexType,
			list[0].IsUnique,
			columnNameList,
			indexOrder,
		)
	}
	return results, nil
}
