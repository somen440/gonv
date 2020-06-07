package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
		errs = append(errs, err) // nolint: staticcheck
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

// SwitchDb switch db
func (gdo *GDO) SwitchDb(name string) error {
	if gdo.HasError() {
		return errors.New(gdo.Error())
	}
	_, err := gdo.db.Exec("use " + name)
	return err
}

// ShowTableStatusLikeResult SHOW TABLE STATUS LIKE result
type ShowTableStatusLikeResult struct {
	Name          string
	Engine        string
	Version       string
	RowFormat     string
	Rows          int
	AvgRowLength  int
	DataLength    int
	MaxDataLength int
	IndexLength   int
	DataFree      int
	AutoIncrement int
	CreateTime    time.Time
	UpdateTime    time.Time
	CheckTime     time.Time
	Collation     string
	Checksum      string
	CreateOptions string
	Comment       string
}

// ShowTableStatusLike SHOW TABLE STATUS LIKE
func (gdo *GDO) ShowTableStatusLike(table string) (*ShowTableStatusLikeResult, error) {
	if gdo.HasError() {
		return nil, errors.New(gdo.Error())
	}
	rows, err := gdo.db.Query("show table status like '" + table + "'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *ShowTableStatusLikeResult
	for rows.Next() {
		if err := rows.Scan(
			&result.Name,
			&result.Engine,
			&result.Version,
			&result.RowFormat,
			&result.Rows,
			&result.AvgRowLength,
			&result.DataLength,
			&result.MaxDataLength,
			&result.IndexLength,
			&result.DataFree,
			&result.AutoIncrement,
			&result.CreateTime,
			&result.UpdateTime,
			&result.CheckTime,
			&result.Collation,
			&result.Checksum,
			&result.CreateOptions,
			&result.Comment,
		); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// SelectPartitionsResult select
type SelectPartitionsResult struct {
	TableCatalog                string
	TableSchema                 string
	TableName                   string
	PartitionName               string
	SubPartitionName            string
	PartitionOrdinalPosition    int
	SubpartitionOrdinalPosition string
	PartitionMethod             string
	SubPartitionMethod          string
	PartitionExpression         string
	SubPartitionExpression      string
	PartitionDescription        string
	TableRows                   int
	AvgRowLength                int
	DataLength                  int
	MaxDataLength               string
	IndexLength                 int
	DataFree                    int
	CreateTime                  time.Time
	UpdateTime                  time.Time
	CheckTime                   time.Time
	Checksum                    string
	PartitionComment            string
	Nodegroup                   string
	TablespaceName              string
}

func (gdo *GDO) SelectPartitions(dbName, tableName string) ([]SelectPartitionsResult, error) {
	if gdo.HasError() {
		return nil, errors.New(gdo.Error())
	}
	rows, err := gdo.db.Query(
		fmt.Sprintf(`select * from information_schema.PARTITIONS
where table_schema = '%s'
	and table_name = '%s' 
order by PARTITION_ORDINAL_POSITION asc;
`, dbName, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SelectPartitionsResult
	for rows.Next() {
		var col SelectPartitionsResult
		if err := rows.Scan(
			&col.TableCatalog,
			&col.TableSchema,
			&col.TableName,
			&col.PartitionName,
			&col.SubPartitionName,
			&col.PartitionOrdinalPosition,
			&col.SubpartitionOrdinalPosition,
			&col.PartitionMethod,
			&col.SubPartitionMethod,
			&col.PartitionExpression,
			&col.SubPartitionExpression,
			&col.PartitionDescription,
			&col.TableRows,
			&col.AvgRowLength,
			&col.DataLength,
			&col.MaxDataLength,
			&col.IndexLength,
			&col.DataFree,
			&col.CreateTime,
			&col.UpdateTime,
			&col.CheckTime,
			&col.Checksum,
			&col.PartitionComment,
			&col.Nodegroup,
			&col.TablespaceName,
		); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// SelectColumnsResult select information_schema.COLUMNS result
type SelectColumnsResult struct {
	TableCatalog           string
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int
	ColumnDefault          string
	IsNullable             string
	DataType               string
	CharacterMaximumLength string
	CharacterOctetLength   string
	NumericPrecision       int
	NumericScale           int
	DateTimePrecision      int
	CharacterSetName       string
	CollationName          string
	ColumnType             string
	ColumnKey              string
	Extra                  string
	Privileges             string
	ColumnComment          string
	GenerationExpression   string
}

func (gdo *GDO) SelectColumns(dbName, tableName string) ([]SelectColumnsResult, int, error) {
	if gdo.HasError() {
		return nil, 0, errors.New(gdo.Error())
	}
	rows, err := gdo.db.Query(
		fmt.Sprintf(`select * from information_schema.COLUMNS
where table_schema = '%s'
	and table_name = '%s' 
order by ORDINAL_POSITION asc;
`, dbName, tableName))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []SelectColumnsResult

	count := 0
	for rows.Next() {
		var column SelectColumnsResult
		if err := rows.Scan(
			&column.TableCatalog,
			&column.TableSchema,
			&column.TableName,
			&column.ColumnName,
			&column.OrdinalPosition,
			&column.ColumnDefault,
			&column.IsNullable,
			&column.DataType,
			&column.CharacterMaximumLength,
			&column.CharacterOctetLength,
			&column.NumericPrecision,
			&column.NumericScale,
			&column.DateTimePrecision,
			&column.CharacterSetName,
			&column.CollationName,
			&column.ColumnType,
			&column.ColumnKey,
			&column.Extra,
			&column.Privileges,
			&column.ColumnComment,
			&column.GenerationExpression,
		); err != nil {
			return nil, 0, err
		}
		count++
	}

	return result, count, nil
}
