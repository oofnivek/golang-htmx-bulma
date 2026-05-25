CREATE TABLE `roles` (
  `id` varchar(10) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `old_id` varchar(36) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
