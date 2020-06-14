# gonv

![actions](https://github.com/somen440/gonv/workflows/ci/badge.svg)
[![codecov](https://codecov.io/gh/somen440/gonv/branch/master/graph/badge.svg)](https://codecov.io/gh/somen440/gonv)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

gonv create migration support tools. from https://github.com/howyi/conv

## Getting Started

```
$ go get -u github.com/somen440/gonv
```

## Running the tests

table schema

```sql
CREATE TABLE `sample` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Modified Time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample table'
```

file schema

```sql
CREATE TABLE `sample` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sample' COMMENT 'Sample Name',
  `description` text COMMENT 'Sample description',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Modified Time',
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='sample table'
```

diff result

```
❯ gonv diff -u root -p test -P 33066 test build/mysql/diff_schema 
...
...
INFO: migrations
*************************** migration up ***************************
ALTER TABLE `sample`
 ADD COLUMN `description` text COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Sample description',
 ADD INDEX `name` (`name`);
*************************** migration down ***************************
ALTER TABLE `sample`
 DROP COLUMN `description`,
 DROP INDEX `name`;
```

exec migration

```
mysql> ALTER TABLE `sample`
    ->  ADD COLUMN `description` text COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Sample description',
    ->  ADD INDEX `name` (`name`);
Query OK, 0 rows affected (0.04 sec)
Records: 0  Duplicates: 0  Warnings: 0
```

after result

```
INFO: migrations
no migrations.
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
