package gonv

import (
	"fmt"

	"github.com/somen440/gonv/structure"
)

// Factory structure factory
type Factory struct {
	gdo *GDO
}

// CreateTableStructure create table structure
func (f *Factory) CreateTableStructure(dbName, tableName string) (*structure.TableStructure, error) {
	if err := f.gdo.SwitchDb(dbName); err != nil {
		return nil, fmt.Errorf("SwitchDb error: %w", err)
	}
	tableStatus, err := f.gdo.ShowTableStatusLike(tableName)
	if err != nil {
		return nil, err
	}
	// createTable := f.gdo.ShowCreateTable(tableName)

	return &structure.TableStructure{
		Table:   tableName,
		Comment: tableStatus.Comment,
		Engine:  tableStatus.Engine,
		// NatableName,
		// tableStatus.Comment,
		// tableStatus.Engine,
	}, nil
}
