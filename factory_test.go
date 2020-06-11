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
	factory, df := setUpTestFactory()
	defer df()

	expected := converter.CreateMockDatabaseStructure()
	actual, err := factory.CreateDatabaseStructure("test")
	assert.Nil(t, err)

	assertEqualTable(t, expected.Map[structure.TableName("sample")], actual.Map[structure.TableName("sample")])
	assertEqualTable(t, expected.Map[structure.TableName("sample_log")], actual.Map[structure.TableName("sample_log")])
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

func assertEqualColumn(t *testing.T, expected, actual []*structure.MySQL57ColumnStructure) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, actualColumn := range actual {
		expectedColumn := expected[i]
		assert.Equal(t, expectedColumn, actualColumn)
	}
}

func assertEqualIndex(t *testing.T, expected, actual []*structure.IndexStructure) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, actualColumn := range actual {
		expectedColumn := expected[i]
		assert.Equal(t, expectedColumn, actualColumn)
	}
}