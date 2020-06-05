package gonv

import (
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setUp(t *testing.T) (*GDO, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gdo := &GDO{
		db:   db,
		errs: []error{},
	}
	return gdo, mock, func() { db.Close() }
}

func TestShowTables(t *testing.T) {
	t.Helper()
	gdo, mock, close := setUp(t)
	defer close()

	rows := sqlmock.NewRows([]string{"table"}).
		AddRow("t_hoge").
		AddRow("t_foo").
		AddRow("t_bar")
	mock.ExpectQuery("show tables;").
		WillReturnRows(rows)

	expected := []string{
		"t_hoge",
		"t_foo",
		"t_bar",
	}
	actual := gdo.ShowTables()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, expected, actual)
}

func TestShowCreateTable(t *testing.T) {
	t.Helper()
	gdo, mock, close := setUp(t)
	defer close()

	contains, _ := ioutil.ReadFile("schema/sample.sql")
	rows := sqlmock.NewRows([]string{"Table", "Create Table"}).
		AddRow("sample", string(contains))
	mock.ExpectQuery("show create table sample;").
		WillReturnRows(rows)

	expected := &ShowCreateTableResult{
		table:  "sample",
		schema: string(contains),
	}
	actual := gdo.ShowCreateTable("sample")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, expected, actual)
}
