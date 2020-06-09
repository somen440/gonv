package adapter

import (
	"testing"

	"github.com/somen440/gonv/structure"
)

func TestTableCreateMigration(t *testing.T) {
	adapter := &Adapter{}

	t1 := createSampleTableSt()
	t2 := createTableSt()

	db1 := &structure.DatabaseStructure{
		structure.TableName("sample"): t1,
	}
}
