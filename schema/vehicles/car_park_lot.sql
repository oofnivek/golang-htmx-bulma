CREATE TABLE `car_park_lot` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `car_park_id` bigint NOT NULL,
  `lot_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `level` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` tinyint(1) NOT NULL,
  `created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `old_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IX_car_park_lot_car_park_id` (`car_park_id`),
  CONSTRAINT `FK_car_park_lot_car_park_car_park_id` FOREIGN KEY (`car_park_id`) REFERENCES `car_park` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
