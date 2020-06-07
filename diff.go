package main

import (
	"fmt"

	"github.com/gookit/color"
)

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
func (d *Diff) Exec(beforeDbName string, schema string) error {
	defer d.gdo.Close()

	color.Info.Tips("create before database structures from db")
	fmt.Println()
	before, err := d.factory.CreateDatabaseStructure(beforeDbName)
	if err != nil {
		return err
	}
	fmt.Println(before.String())

	color.Info.Tips("create after database structures from schema")
	fmt.Println()
	after, err := d.factory.CreateDatabaseStructureFromSchema("tmp_"+beforeDbName, schema)
	if err != nil {
		return err
	}
	fmt.Println(after.String())

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
