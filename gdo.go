package main

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// GDO golang data object
type GDO struct {
	conf *DBConfig
	db   *sql.DB
	errs []error
}

// NewGDO DGO を生成
func NewGDO(conf *DBConfig) *GDO {
	errs := []error{}
	db, err := sql.Open(conf.Driver.AsString(), conf.DataSourceName())
	if err != nil {
		errs = append(errs, err)
	}
	return &GDO{
		conf: conf,
		db:   db,
		errs: []error{},
	}
}

// Close 終了処理
func (gdo *GDO) Close() {
	gdo.db.Close()
}

// HasError error を持っているか
func (gdo *GDO) HasError() bool {
	return len(gdo.errs) > 0
}

// Error error 情報の取得
func (gdo *GDO) Error() string {
	var results []string
	for _, err := range gdo.errs {
		results = append(results, err.Error())
	}
	return strings.Join(results, "\n")
}

// ShowTables table 一覧の取得
func (gdo *GDO) ShowTables() []string {
	if gdo.HasError() {
		return []string{}
	}

	results := []string{}

	rows, err := gdo.db.Query("show tables;")
	if err != nil {
		gdo.errs = append(gdo.errs, err)
		return []string{}
	}
	defer rows.Close()

	for rows.Next() {
		var table string

		if err := rows.Scan(&table); err != nil {
			gdo.errs = append(gdo.errs, err)
			return []string{}
		}

		results = append(results, table)
	}

	return results
}

// ShowCreateTableResult show create table した時の結果
type ShowCreateTableResult struct {
	table  string
	schema string
}

// ShowCreateTable table create 文の取得
func (gdo *GDO) ShowCreateTable(table string) *ShowCreateTableResult {
	if gdo.HasError() {
		return nil
	}

	rows, err := gdo.db.Query("show create table " + table + ";")
	if err != nil {
		gdo.errs = append(gdo.errs, err)
		return nil
	}
	defer rows.Close()

	var result ShowCreateTableResult

	for rows.Next() {
		if err := rows.Scan(&result.table, &result.schema); err != nil {
			gdo.errs = append(gdo.errs, err)
			return nil
		}
	}

	return &result
}
