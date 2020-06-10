package converter

import (
	"testing"

	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
	"github.com/stretchr/testify/assert"
)

type expectedMigration struct {
	migration    migration.Migration
	expectedUp   string
	expectedDown string
}

func TestToTableDropMigration(t *testing.T) {
	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	converter := &Converter{}
	migrationList := converter.ToTableDropMigration(db1, db2)
	assert.Nil(t, migrationList)
}

func TestToTableCreateMigration(t *testing.T) {
	converter := &Converter{}

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	migrationList := converter.ToTableCreateMigration(db1, db2)
	assert.Nil(t, migrationList)

	delete(db1.Map, structure.TableName("sample"))

	migrationList = converter.ToTableCreateMigration(db1, db2)

	sql := "CREATE TABLE sample (\n"
	sql += " `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',\n"
	sql += " `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name',\n"
	sql += " `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',\n"
	sql += " `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Modified Time',\n"
	sql += " PRIMARY KEY (`id`)\n"
	sql += ") ENGINE=InnoDB DEFAULT CHARASET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample table'"

	actuals := migrationList.List()
	assert.Len(t, actuals, 1)

	tests := []*expectedMigration{
		{
			actuals[0],
			sql,
			"DROP TABLE sample",
		},
	}
	assertEqualMigration(t, tests)
}

func assertEqualMigration(t *testing.T, targets []*expectedMigration) {
	for _, tt := range targets {
		actualUp := tt.migration.UpQuery()
		assert.Equal(t, tt.expectedUp, actualUp)

		actualDown := tt.migration.DownQuery()
		assert.Equal(t, tt.expectedDown, actualDown)
	}
}
