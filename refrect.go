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
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Reflect スキーマを反映する
type Reflect struct {
	gdo    *GDO
	output string
}

// NewReflect Reflect object の生成
func NewReflect(conf *DBConfig, output string) *Reflect {
	return &Reflect{
		gdo:    NewGDO(conf),
		output: output,
	}
}

// Exec reflect 実行
func (r *Reflect) Exec() error {
	results := r.getSchemas()

	if r.gdo.HasError() {
		return errors.New(r.gdo.Error())
	}

	if err := os.RemoveAll(r.output); err != nil {
		return err
	}
	if err := os.MkdirAll(r.output, os.ModePerm); err != nil {
		return err
	}

	for _, result := range results {
		filename := filepath.Join(r.output, result.table+".sql")
		fmt.Println(filename)
		if err := r.outputFile(filename, result.schema); err != nil {
			return err
		}
	}

	return nil
}

func (r *Reflect) getSchemas() []*ShowCreateTableResult {
	defer r.gdo.Close()

	var results []*ShowCreateTableResult
	for _, table := range r.gdo.ShowTables() {
		result := r.gdo.ShowCreateTable(table)
		results = append(results, result)
	}

	return results
}

func (r *Reflect) outputFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()
	_, err = fmt.Fprintln(file, content)
	return err
}
