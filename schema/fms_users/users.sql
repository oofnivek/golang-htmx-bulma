CREATE TABLE `users` (
  `email` varchar(50) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `first_name` varchar(45) NOT NULL,
  `last_name` varchar(45) NOT NULL,
  `mobile` varchar(50) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `designation` varchar(250) NOT NULL,
  `department` varchar(250) NOT NULL,
  `is_enabled` tinyint(1) NOT NULL,
  `created_at_utc` datetime(3) NOT NULL,
  `updated_at_utc` datetime(3) NOT NULL,
  `role_id` varchar(10) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `old_role_id` varchar(36) DEFAULT NULL,
  `password_reset_request_id` binary(16) DEFAULT NULL,
  `password_hash` char(84) CHARACTER SET ascii COLLATE ascii_general_ci DEFAULT NULL,
  `old_cms_user_id` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`email`),
  KEY `ix_users_role_id` (`role_id`),
  CONSTRAINT `fk_users_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
