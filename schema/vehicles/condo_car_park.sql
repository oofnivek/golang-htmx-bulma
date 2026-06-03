CREATE TABLE `condo_car_park` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `condo_id` bigint NOT NULL,
  `car_park_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IX_condo_car_park_condo_id` (`condo_id`),
  CONSTRAINT `FK_condo_car_park_condo_condo_id` FOREIGN KEY (`condo_id`) REFERENCES `condo` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
