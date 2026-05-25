CREATE TABLE `vehicle_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `vehicle_type_id` bigint NOT NULL,
  `vehicle_make_id` bigint NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` tinyint(1) NOT NULL,
  `created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `old_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IX_vehicle_model_vehicle_make_id` (`vehicle_make_id`),
  KEY `IX_vehicle_model_vehicle_type_id` (`vehicle_type_id`),
  CONSTRAINT `FK_vehicle_model_vehicle_make_vehicle_make_id` FOREIGN KEY (`vehicle_make_id`) REFERENCES `vehicle_make` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_vehicle_model_vehicle_type_vehicle_type_id` FOREIGN KEY (`vehicle_type_id`) REFERENCES `vehicle_type` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
