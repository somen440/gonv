package gonv

// Diff スキーマと db の diff を取る
type Diff struct {
	before *DatabaseStructure
	after  *DatabaseStructure
}

// NewDiff Diff object の生成
func NewDiff(conf *DBConfig, schema string) *Diff {
	gdo := NewGDO(conf)

	before := gdo.CreateDatabaseStructure()
	after := CreateDatabaseStructureFromSchema(gdo, schema)

	return &Diff{
		before: before,
		after:  after,
	}
}

// Exec 実行
func (d *Diff) Exec() error {
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
