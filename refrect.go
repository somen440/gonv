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

	if _, err := os.Stat(r.output); os.IsNotExist(err) {
		if err := os.MkdirAll(r.output, os.ModePerm); err != nil {
			return err
		}
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
