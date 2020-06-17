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
	"testing"

	"github.com/somen440/gonv/converter"
	"github.com/somen440/gonv/structure"
	"github.com/stretchr/testify/assert"
)

func setUpTestFactory() (*Factory, func()) {
	config := &DBConfig{
		Driver:   MySQL,
		User:     "root",
		Password: "test",
		Database: "test",
		Host:     "localhost",
		Port:     "33066",
	}
	gdo := NewGDO(config)
	factory := &Factory{
		gdo: gdo,
	}
	return factory, gdo.Close
}

func TestCreateDatabaseStructure(t *testing.T) {
	i := 1000
	for i > 0 {
		func() {
			factory, df := setUpTestFactory()
			defer df()

			expected := converter.CreateMockDatabaseStructure()
			actual, err := factory.CreateDatabaseStructure("test")
			assert.Nil(t, err)

			assertEqualTable(t, expected.Map[structure.TableName("sample")], actual.Map[structure.TableName("sample")])
			assertEqualTable(t, expected.Map[structure.TableName("sample_log")], actual.Map[structure.TableName("sample_log")])
		}()
		i--
	}
}

func assertEqualTable(t *testing.T, expected, actual *structure.TableStructure) {
	assert.Equal(t, expected.Table, actual.Table)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Comment, actual.Comment)
	assert.Equal(t, expected.Engine, actual.Engine)
	assert.Equal(t, expected.Collate, actual.Collate)
	assert.Equal(t, expected.DefaultCharset, actual.DefaultCharset)
	assert.Equal(t, expected.Partition, actual.Partition)
	assertEqualColumn(t, expected.ColumnStructureList, actual.ColumnStructureList)
	assertEqualIndex(t, expected.IndexStructureList, actual.IndexStructureList)
	assert.Equal(t, expected.Properties, actual.Properties)
}

func assertEqualColumn(t *testing.T, expected, actual map[structure.ColumnField]*structure.MySQL57ColumnStructure) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, actualColumn := range actual {
		expectedColumn := expected[i]
		assert.Equal(t, expectedColumn, actualColumn)
	}
}

func assertEqualIndex(t *testing.T, expected, actual map[structure.IndexKey]*structure.IndexStructure) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for key, actualColumn := range actual {
		expectedColumn, ok := expected[key]
		assert.True(t, ok)
		assert.Equal(t, expectedColumn, actualColumn)
	}
}
