CREATE TABLE `fuel_card` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `card_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `fuel_company_id` bigint NOT NULL,
  `pin_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `vehicle_id` bigint DEFAULT NULL,
  `status` tinyint(1) NOT NULL,
  `created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `source_db` varchar(10) DEFAULT NULL,
  `old_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IX_fuel_card_company_id` (`fuel_company_id`),
  KEY `IX_fuel_card_vehicle_id` (`vehicle_id`),
  CONSTRAINT `FK_fuel_card_fuel_company_fuel_company_id` FOREIGN KEY (`fuel_company_id`) REFERENCES `fuel_company` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_fuel_card_vehicle_vehicle_id` FOREIGN KEY (`vehicle_id`) REFERENCES `vehicle` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
