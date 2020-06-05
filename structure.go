package gonv

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
)

// DatabaseStructure db
type DatabaseStructure struct {
	Tables []*TableStructure
}

// TableStructure table
type TableStructure struct {
	Rows []*DescribeRow
}

// DescribeRow describe <table> row
type DescribeRow struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default []byte
	Extra   string
}

// CreateTableStructure describe tables
func (gdo *GDO) CreateTableStructure(table string) *TableStructure {
	if gdo.HasError() {
		return nil
	}
	rows, err := gdo.db.Query("describe " + table)
	if err != nil {
		gdo.errs = append(gdo.errs, err)
		return nil
	}

	results := &TableStructure{
		Rows: []*DescribeRow{},
	}
	for rows.Next() {
		result := &DescribeRow{}
		if err := rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra); err != nil {
			gdo.errs = append(gdo.errs, err)
			return nil
		}
		results.Rows = append(results.Rows, result)
	}

	return results
}

// CreateDatabaseStructure db
func (gdo *GDO) CreateDatabaseStructure() *DatabaseStructure {
	result := &DatabaseStructure{
		Tables: []*TableStructure{},
	}

	for _, table := range gdo.ShowTables() {
		tableSt := gdo.CreateTableStructure(table)
		result.Tables = append(result.Tables, tableSt)
	}

	return result
}

// CreateDatabaseStructureFromSchema スキーマからデータベース構造を作成
func CreateDatabaseStructureFromSchema(gdo *GDO, dir string) *DatabaseStructure {
	db, err := sql.Open(gdo.conf.Driver.AsString(), gdo.conf.DataSourceNameNoDatabase())
	if err != nil {
		panic(err)
	}

	dbName := "tmp_" + gdo.conf.Database

	_, err = db.Exec("create database " + dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("use " + dbName)
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		sql, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}

	conf := *gdo.conf
	conf.Database = dbName
	return NewGDO(&conf).CreateDatabaseStructure()
}
