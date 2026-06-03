CREATE TABLE `condo_acknowledgement` (
  `user_id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
