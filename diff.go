package main

import "fmt"

// Diff スキーマと db の diff を取る
type Diff struct {
	gdo     *GDO
	factory *Factory
}

// NewDiff Diff object の生成
func NewDiff(conf *DBConfig) *Diff {
	gdo := NewGDO(conf)
	factory := &Factory{
		gdo: gdo,
	}
	return &Diff{
		gdo:     gdo,
		factory: factory,
	}
}

// Exec 実行
func (d *Diff) Exec(beforeDbName string) error {
	defer d.gdo.Close()

	fmt.Println("create table structures")
	for _, table := range d.gdo.ShowTables() {
		fmt.Println("table: " + table)
		tableSt, err := d.factory.CreateTableStructure(beforeDbName, table)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", tableSt)
	}

	return nil
}

// Migration interface
type Migration interface {
	Up() string
	Down() string
}

func (d *Diff) generate(before, after interface{}) (Migration, error) {
	// DROP → MODIFY → ADD
	// table drop

	// view drop

	// table alter

	// view alter

	// view rename

	// table create

	// view create

	return nil, nil
}
