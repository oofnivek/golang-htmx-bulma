CREATE TABLE `global_setting` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(50) NOT NULL,
  `value` varchar(250) NOT NULL,
  `remark` varchar(250) DEFAULT NULL,
  `country_code` varchar(2) DEFAULT NULL,
  `created_by` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_by` varchar(50) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `global_setting_pk` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
