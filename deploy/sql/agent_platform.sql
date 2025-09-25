USE `easy_chat`;

CREATE TABLE IF NOT EXISTS `agents` (
    `id` varchar(32) NOT NULL,
    `user_id` varchar(24) NOT NULL,
    `code` varchar(64) NOT NULL,
    `name` varchar(128) NOT NULL,
    `description` text DEFAULT NULL,
    `status` varchar(32) NOT NULL DEFAULT 'draft',
    `model` varchar(128) NOT NULL DEFAULT 'gpt-3.5-turbo',
    `prompt` text NOT NULL,
    `tools` json NOT NULL,
    `memory_strategy` varchar(64) NOT NULL DEFAULT 'recent',
    `config` json NOT NULL,
    `created_by` varchar(24) NOT NULL,
    `updated_by` varchar(24) DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_agents_code` (`code`),
    UNIQUE KEY `uniq_agents_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `agent_versions` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `agent_id` varchar(32) NOT NULL,
    `version` int NOT NULL,
    `prompt` text NOT NULL,
    `tools` json NOT NULL,
    `config` json NOT NULL,
    `created_by` varchar(24) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_agent_versions_agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
