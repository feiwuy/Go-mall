-- mysqldump导出
-- 主库 + 测试库（测试库是 main_test.go:29 切换的目标）
CREATE DATABASE IF NOT EXISTS `go_mall_db`      DEFAULT CHARACTER SET utf8mb4;
CREATE DATABASE IF NOT EXISTS `go_mall_test_db` DEFAULT CHARACTER SET utf8mb4;

USE `go_mall_db`;

CREATE TABLE IF NOT EXISTS `users` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username`   varchar(64) NOT NULL,
  `password`   varchar(60) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `todo_lists` (
  `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title`      longtext,
  `content`    longtext,
  `completed`  tinyint(1) DEFAULT NULL,
  `user_id`    bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_todo_lists_deleted_at` (`deleted_at`),
  KEY `idx_todo_lists_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 测试库直接复制主库结构（含索引、软删列），避免两份 DDL 以后改漂移
CREATE TABLE IF NOT EXISTS `go_mall_test_db`.`users`      LIKE `go_mall_db`.`users`;
CREATE TABLE IF NOT EXISTS `go_mall_test_db`.`todo_lists` LIKE `go_mall_db`.`todo_lists`;
