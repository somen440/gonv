package adapter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// TableCreateMigration adapter DatabaseStructure -> TableCreateMigration
func (a *Adapter) TableCreateMigration(before, after *structure.DatabaseStructure) *migration.TableCreateMigration {
	return nil
}
