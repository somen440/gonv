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
func (d *Diff) Exec(beforeDbName string, schema string, ignores []string) error {
	defer d.gdo.Close()

	color.Info.Tips("create before database structures from db")
	fmt.Println()
	before, err := d.factory.CreateDatabaseStructure(beforeDbName, ignores)
	if err != nil {
		return err
	}
	fmt.Println(before)

	color.Info.Tips("create after database structures from schema")
	fmt.Println()
	after, err := d.factory.CreateDatabaseStructureFromSchema("tmp_"+beforeDbName, schema, ignores)
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
