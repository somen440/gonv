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

func setUp() (*Converter, func()) {
	return NewConverter(), func() {}
}

func TestToTableDropMigration(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	delete(db2.Map, structure.TableName("sample_log"))

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAsk{
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
			"DROP TABLE sample_log",
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

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAsk{
		DroppedTableList: []structure.TableName{
			structure.TableName("sample_log"),
		},
	})
	assert.Nil(t, migrationList)
	assert.True(t, converter.HasError())

	expectedErr := "ToTableDropMigration not found table sample_log from map[sample:name: sample\ntype: table\ncomment: sample table\nengine: InnoDB\ndefault_charset: utf8mb4\ncollate: utf8mb4_unicode_ci\nproperties: \ncolumns:\n\tfield: id\n\t\ttype: bigint(20)\n\t\tdefault: \n\t\tcomment: Sample ID\n\t\tattributes:\n\t\t\t- auto_increment\n\t\t\t- unsigned\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID'\n\tfield: name\n\t\ttype: varchar(255)\n\t\tdefault: sample\n\t\tcomment: Sample Name\n\t\tattributes:\n\t\tcollation_name: utf8mb4_unicode_ci\n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name'\n\tfield: created\n\t\ttype: datetime\n\t\tdefault: CURRENT_TIMESTAMP\n\t\tcomment: Created Time\n\t\tattributes:\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time'\n\tfield: modified\n\t\ttype: datetime\n\t\tdefault: CURRENT_TIMESTAMP\n\t\tcomment: Modified Time\n\t\tattributes:\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Modified Time'\nindex:\n\tPRIMARY KEY (`id`)\n]"
	assert.Equal(t, expectedErr, converter.Err.Error())
}

func TestToTableDropMigrationWithFoundTable(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()

	migrationList := converter.ToTableDropMigration(db1, db2, &TableAsk{
		DroppedTableList: []structure.TableName{
			structure.TableName("sample_log"),
		},
	})
	assert.Nil(t, migrationList)
	assert.True(t, converter.HasError())

	expectedErr := "ToTableDropMigration found table sample_log from map[sample:name: sample\ntype: table\ncomment: sample table\nengine: InnoDB\ndefault_charset: utf8mb4\ncollate: utf8mb4_unicode_ci\nproperties: \ncolumns:\n\tfield: id\n\t\ttype: bigint(20)\n\t\tdefault: \n\t\tcomment: Sample ID\n\t\tattributes:\n\t\t\t- auto_increment\n\t\t\t- unsigned\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID'\n\tfield: name\n\t\ttype: varchar(255)\n\t\tdefault: sample\n\t\tcomment: Sample Name\n\t\tattributes:\n\t\tcollation_name: utf8mb4_unicode_ci\n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name'\n\tfield: created\n\t\ttype: datetime\n\t\tdefault: CURRENT_TIMESTAMP\n\t\tcomment: Created Time\n\t\tattributes:\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time'\n\tfield: modified\n\t\ttype: datetime\n\t\tdefault: CURRENT_TIMESTAMP\n\t\tcomment: Modified Time\n\t\tattributes:\n\t\tcollation_name: \n\t\tproperties:\n\t\tgeneration_expression: \n\t\tgenerate: `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Modified Time'\nindex:\n\tPRIMARY KEY (`id`)\n sample_log:name: sample_log\ntype: table\ncomment: sample log table\nengine: InnoDB\ndefault_charset: utf8mb4\ncollate: utf8mb4_unicode_ci\nproperties: \ncolumns:\npartitions:\nPARTITION BY LIST(month)\n(PARTITION p1 VALUES IN (1),\n PARTITION p2 VALUES IN (2),\n PARTITION p3 VALUES IN (3),\n PARTITION p4 VALUES IN (4),\n PARTITION p5 VALUES IN (5),\n PARTITION p6 VALUES IN (6),\n PARTITION p7 VALUES IN (7),\n PARTITION p8 VALUES IN (8),\n PARTITION p9 VALUES IN (9),\n PARTITION p10 VALUES IN (10),\n PARTITION p11 VALUES IN (11),\n PARTITION p12 VALUES IN (12))]"
	assert.Equal(t, expectedErr, converter.Err.Error())
}

func TestToTableCreateMigration(t *testing.T) {
	converter, _ := setUp()

	db1 := CreateMockDatabaseStructure()
	db2 := CreateMockDatabaseStructure()
	migrationList := converter.ToTableCreateMigration(db1, db2)
	assert.Nil(t, migrationList)

	delete(db1.Map, structure.TableName("sample"))

	migrationList = converter.ToTableCreateMigration(db1, db2)

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
