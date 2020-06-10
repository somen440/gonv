package converter

import (
	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
)

// Converter structure -> migration converter
type Converter struct{}

// ConvertAll cnvert struct -> migration all
//   1. DROP
//   2. MODIFY
//     2.1. DROP Index
//     2.2. DROP
//     2.3. MODIFY
//     2.4. ADD
//     2.5. ADD
//   3. ADD
func (c *Converter) ConvertAll(
	before, after *structure.DatabaseStructure,
	ask *ModifiedAsk,
) *migration.List {
	results := &migration.List{}

	tMigration := c.convertTableAll(before, after, ask.Table)

	results.Merge(tMigration)

	return results
}
