-- Schema for golang-todo-list-api (MySQL 8+)

CREATE TABLE IF NOT EXISTS `user` (
  `user_id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_uuid` VARCHAR(36) NOT NULL,
  `nick_name` VARCHAR(255) NOT NULL,
  `password` VARBINARY(255) NOT NULL,
  `active` BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_uuid` (`user_uuid`),
  UNIQUE KEY `nick_name` (`nick_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `list_item` (
  `list_item_id` BIGINT NOT NULL AUTO_INCREMENT,
  `list_item_uuid` VARCHAR(36) NOT NULL,
  `user_id` BIGINT NOT NULL,
  `title` VARCHAR(512) NOT NULL,
  `description` TEXT,
  `active` BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (`list_item_id`),
  UNIQUE KEY `list_item_uuid` (`list_item_uuid`),
  KEY `list_item_user_id` (`user_id`),
  CONSTRAINT `list_item_user_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
