CREATE DATABASE IF NOT EXISTS `easy_chat` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `easy_chat`;

CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
    `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `nickname` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `status` tinyint NOT NULL DEFAULT 0 COMMENT 'bitmask: 0 normal, 1 disabled, 2 pending, 4 deleted',
    `user_type` tinyint NOT NULL DEFAULT 0 COMMENT '0 human, 1 agent',
    `sex` tinyint DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- goctl model mysql ddl -src="./deploy/sql/user.sql" -dir="./apps/user/models" -c
