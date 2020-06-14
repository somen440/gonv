package converter

import (
	"strings"
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

func setUp() (*Converter, func()) {
	return NewConverter(), func() {}
}

func TestToTableDropMigration(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	delete(db2.Map, structure.TableName("sample_log"))

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAnswer{
		DroppedTableList: []structure.TableName{
			structure.TableName("sample_log"),
		},
	})

	actuals := migrationList.List()
	assert.Len(t, actuals, 1)

	sql := "CREATE TABLE `sample_log` (\n"
	sql += " `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n"
	sql += " `month` tinyint(2) unsigned NOT NULL,\n"
	sql += " `sample_id` bigint(20) unsigned NOT NULL,\n"
	sql += " `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,\n"
	sql += " `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"
	sql += " `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n"
	sql += " PRIMARY KEY (`id`, `month`),\n"
	sql += " KEY `sample_id` (`sample_id`)\n"
	sql += ") ENGINE=InnoDB DEFAULT CHARASET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample log table'\n"
	sql += "/*!50100 PARTITION BY LIST(month)\n"
	sql += "(PARTITION p1 VALUES IN (1) ENGINE = InnoDB,\n"
	sql += " PARTITION p2 VALUES IN (2) ENGINE = InnoDB,\n"
	sql += " PARTITION p3 VALUES IN (3) ENGINE = InnoDB,\n"
	sql += " PARTITION p4 VALUES IN (4) ENGINE = InnoDB,\n"
	sql += " PARTITION p5 VALUES IN (5) ENGINE = InnoDB,\n"
	sql += " PARTITION p6 VALUES IN (6) ENGINE = InnoDB,\n"
	sql += " PARTITION p7 VALUES IN (7) ENGINE = InnoDB,\n"
	sql += " PARTITION p8 VALUES IN (8) ENGINE = InnoDB,\n"
	sql += " PARTITION p9 VALUES IN (9) ENGINE = InnoDB,\n"
	sql += " PARTITION p10 VALUES IN (10) ENGINE = InnoDB,\n"
	sql += " PARTITION p11 VALUES IN (11) ENGINE = InnoDB,\n"
	sql += " PARTITION p12 VALUES IN (12) ENGINE = InnoDB) */"
	tests := []*expectedMigration{
		{
			actuals[0],
			"DROP TABLE `sample_log`",
			sql,
		},
	}
	assertEqualMigration(t, tests)
}

func TestToTableDropMigrationWithNotFoundTable(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	delete(db1.Map, structure.TableName("sample_log"))

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAnswer{
		DroppedTableList: []structure.TableName{
			structure.TableName("sample_log"),
		},
	})
	assert.Nil(t, migrationList)
	assert.True(t, converter.HasError())
	assert.True(t, strings.Contains(converter.Err.Error(), "ToTableDropMigration not found table sample_log "))
}

func TestToTableDropMigrationWithFoundTable(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAnswer{
		DroppedTableList: []structure.TableName{
			structure.TableName("sample_log"),
		},
	})
	assert.Nil(t, migrationList)
	assert.True(t, converter.HasError())
	assert.True(t, strings.Contains(converter.Err.Error(), "ToTableDropMigration found table sample_log "))
}

func TestToTableAlterMigrationAll(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	ask := &TableAnswer{
		RenamedTableList: map[structure.TableName]structure.TableName{
			structure.TableName("sample_log"): structure.TableName("m_sample_log"),
		},
	}

	db2.Map[structure.TableName("m_sample_log")] = db2.Map[structure.TableName("sample_log")]
	db2.Map[structure.TableName("m_sample_log")].Table = "m_sample_log"
	db2.Map[structure.TableName("m_sample_log")].Comment = "m_sample_log"

	db2.Map[structure.TableName("m_sample_log")].IndexStructureList[structure.IndexKey("name")] = structure.NewIndexStructure("name", "BTREE", false, []string{"name"}, 1)

	db2.Map[structure.TableName("m_sample_log")].ColumnStructureList[structure.ColumnField("sample_id")].Type = "int(20)"

	delete(db2.Map, structure.TableName("sample_log"))
	delete(db2.Map[structure.TableName("m_sample_log")].IndexStructureList, structure.IndexKey("PRIMARY"))
	delete(db2.Map[structure.TableName("m_sample_log")].ColumnStructureList, structure.ColumnField("modified"))

	migrationList := converter.ToTableAlterMigrationAll(db1, db2, ask)

	up := "ALTER TABLE `sample_log`\n"
	up += " RENAME TO `m_sample_log`,\n"
	up += " COMMENT 'm_sample_log',\n"
	up += " DROP PRIMARY KEY,\n"
	up += " DROP COLUMN `modified`,\n"
	up += " CHANGE `sample_id` `sample_id` int(20) unsigned NOT NULL,\n"
	up += " ADD INDEX `name` (`name`);"

	down := "ALTER TABLE `m_sample_log`\n"
	down += " RENAME TO `sample_log`,\n"
	down += " COMMENT 'sample log table',\n"
	down += " ADD PRIMARY KEY (`id`, `month`),\n"
	down += " ADD COLUMN `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n"
	down += " CHANGE `sample_id` `sample_id` bigint(20) unsigned NOT NULL,\n"
	down += " DROP INDEX `name`;"

	actuals := migrationList.List()
	assert.Len(t, actuals, 1)

	tests := []*expectedMigration{
		{
			actuals[0],
			up,
			down,
		},
	}
	assertEqualMigration(t, tests)
}

func TestToTableCreateMigration(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	delete(db1.Map, structure.TableName("sample"))

	migrationList := converter.ToTableCreateMigration(db1, db2)

	sql := "CREATE TABLE `sample` (\n"
	sql += " `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',\n"
	sql += " `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name',\n"
	sql += " `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',\n"
	sql += " `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Modified Time',\n"
	sql += " PRIMARY KEY (`id`)\n"
	sql += ") ENGINE=InnoDB DEFAULT CHARASET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample table'"

	actuals := migrationList.List()
	assert.Len(t, actuals, 1)

	tests := []*expectedMigration{
		{
			actuals[0],
			sql,
			"DROP TABLE `sample`",
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
