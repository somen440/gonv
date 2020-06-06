CREATE TABLE `sample` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Sample ID',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Samole Name',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Modified Time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
