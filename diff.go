package main

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/somen440/gonv/converter"
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// implements check
var (
	_ migration.ColumnStructure            = &structure.MySQL57ColumnStructure{}
	_ migration.ModifiedColumnStructure    = &structure.ModifiedColumnStructure{}
	_ migration.ModifiedColumnStructureSet = &structure.ModifiedColumnStructureSet{}
	_ migration.IndexStructure             = &structure.IndexStructure{}
	_ migration.TableStructure             = &structure.TableStructure{}
	_ migration.ViewStructure              = &structure.ViewStructure{}
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
	fmt.Println(before)

	color.Info.Tips("create after database structures from schema")
	fmt.Println()
	after, err := d.factory.CreateDatabaseStructureFromSchema("tmp_"+beforeDbName, schema)
	if err != nil {
		return err
	}
	fmt.Println(after)

	// operate ask
	color.Info.Tips("question")
	o := NewOperator(before, after)
	answer := o.Ask()

	color.Info.Tips("migrations")
	c := converter.NewConverter()
	mList := c.ConvertAll(before, after, answer)
	fmt.Println(mList)

	return nil
}
