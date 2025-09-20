CREATE DATABASE IF NOT EXISTS `easy_chat` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `easy_chat`;

CREATE TABLE IF NOT EXISTS `friends` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `remark` varchar(255) DEFAULT NULL,
   `add_source` tinyint DEFAULT NULL,
   `status` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '0正常 1删除 2拉黑 3免打扰',
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uniq_user_friend` (`user_id`,`friend_uid`),
   KEY `idx_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `friend_requests` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `req_msg` varchar(255) DEFAULT NULL,
   `req_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `handle_result` tinyint DEFAULT NULL,
   `handle_msg` varchar(255) DEFAULT NULL,
   `handled_at` timestamp NULL DEFAULT NULL,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `groups` (
   `id` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
   `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
   `icon` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
   `status` tinyint DEFAULT NULL,
   `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `group_type` int(11) NOT NULL,
   `is_verify` boolean NOT NULL,
   `notification` varchar(255) DEFAULT NULL,
   `notification_uid` varchar(64) DEFAULT NULL,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `group_members` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `role_level` tinyint NOT NULL,
   `join_time` timestamp NULL DEFAULT NULL,
   `join_source` tinyint DEFAULT NULL,
   `inviter_uid` varchar(64) DEFAULT NULL,
   `operator_uid` varchar(64) DEFAULT NULL,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `group_requests` (
   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
   `req_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
   `req_msg` varchar(255) DEFAULT NULL,
   `req_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
   `join_source` tinyint DEFAULT NULL,
   `inviter_user_id` varchar(64) DEFAULT NULL,
   `handle_user_id` varchar(64) DEFAULT NULL,
   `handle_time` timestamp NULL DEFAULT NULL,
   `handle_result` tinyint DEFAULT NULL,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- goctl model mysql ddl -src="./deploy/sql/social.sql" -dir="./apps/social/socialmodels" -c
