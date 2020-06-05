package gonv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Generate migration ファイルの生成
type Generate struct {
	output string
}

// NewGenerate Generate を生成
func NewGenerate(output string) *Generate {
	return &Generate{
		output: output,
	}
}

// Exec 実行
func (g *Generate) Exec() error {
	if _, err := os.Stat(g.output); os.IsNotExist(err) {
		if err := os.MkdirAll(g.output, os.ModePerm); err != nil {
			return err
		}
	}

	version := strconv.Itoa(int(time.Now().Unix()))

	migrationTypes := []string{
		"create",
		"update",
		"alter",
		"delete",
	}
	for i, v := range migrationTypes {
		fmt.Printf("[%d] %s\n", i, v)
	}
	fmt.Print("> ")
	var titleIndex int
	fmt.Scan(&titleIndex)
	title := migrationTypes[titleIndex]

	var table string
	fmt.Println("table ?")
	fmt.Print("> ")
	fmt.Scan(&table)

	upF := filepath.Join(g.output, fmt.Sprintf("%s_%s_%s.up.sql", version, title, table))
	if err := ioutil.WriteFile(upF, []byte(""), 0644); err != nil {
		return err
	}
	fmt.Println(upF)

	downF := filepath.Join(g.output, fmt.Sprintf("%s_%s_%s.down.sql", version, title, table))
	if err := ioutil.WriteFile(downF, []byte(""), 0644); err != nil {
		return err
	}
	fmt.Println(downF)

	return nil
}
