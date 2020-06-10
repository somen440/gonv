package converter

import (
	"testing"

	"github.com/somen440/gonv/migration"
	"github.com/somen440/gonv/structure"
	"github.com/stretchr/testify/assert"
)

func TestTableCreateMigrationWithEqual(t *testing.T) {
	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	converter := &Converter{}
	migrationList := converter.ConvertTableDropMigration(db1, db2)
	assert.Nil(t, migrationList)
}

func TestTableCreateMigration(t *testing.T) {
	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	delete(db1.Map, structure.TableName("sample"))

	converter := &Converter{}
	migrationList := converter.ConvertTableCreateMigration(db1, db2)

	sql := "CREATE TABLE sample (\n"
	sql += " `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',\n"
	sql += " `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name',\n"
	sql += " `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',\n"
	sql += " `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Modified Time',\n"
	sql += " PRIMARY KEY (`id`)\n"
	sql += ") ENGINE=InnoDB DEFAULT CHARASET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample table'"

	tests := []struct {
		migration    migration.Migration
		expectedUp   string
		expectedDown string
	}{
		{
			migrationList.List()[0],
			sql,
			"DROP TABLE sample",
		},
	}

	for _, tt := range tests {
		actualUp := tt.migration.UpQuery()
		assert.Equal(t, tt.expectedUp, actualUp)

		actualDown := tt.migration.DownQuery()
		assert.Equal(t, tt.expectedDown, actualDown)
	}
}
