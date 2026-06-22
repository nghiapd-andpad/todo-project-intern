CREATE TABLE `audit_logs` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `event_name` VARCHAR(100) NOT NULL,
    `entity_type` VARCHAR(50) NOT NULL,
    `entity_id` BIGINT NOT NULL,
    `actor_id` BIGINT NOT NULL,
    `payload` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_audit_logs_entity` (`entity_type`, `entity_id`),
    INDEX `idx_audit_logs_actor_id` (`actor_id`),
    INDEX `idx_audit_logs_event_name` (`event_name`),
    INDEX `idx_audit_logs_created_at` (`created_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci;

CREATE TABLE `outbox_events` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `event_name` VARCHAR(100) NOT NULL,
    `routing_key` VARCHAR(100) NOT NULL,
    `payload` JSON NOT NULL,
    `status` VARCHAR(20) NOT NULL,
    `retry_count` INT NOT NULL DEFAULT 0,
    `last_error` TEXT NULL,
    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,
    `published_at` DATETIME(6) NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_outbox_events_status_created_at` (`status`, `created_at`),
    INDEX `idx_outbox_events_routing_key` (`routing_key`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_general_ci;