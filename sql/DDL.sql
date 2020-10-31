-- DB version : 10.4.6-MariaDB

-- DDL 初始化
-- Use CHARACTER utf8mb4 COLLATE utf8mb4_general_ci
DROP DATABASE IF EXISTS `sp_blog`;
CREATE DATABASE IF NOT EXISTS `sp_blog` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- 文章 Table
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '流水號主鍵',
  `created_at` datetime(3) DEFAULT NULL COMMENT '建立時間',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '修改時間',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '刪除時間',
  `title` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '標題',
  `desc` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描敘',
  `content` TEXT COMMENT '內容',
  `status` TINYINT(3) NOT NULL DEFAULT '1' COMMENT '狀態 0:禁用, 1:啟用',
  PRIMARY KEY (`id`),
  INDEX `idx_articles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章';

-- 標籤 Table
DROP TABLE IF EXISTS `article_labels`;
CREATE TABLe `article_labels` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '流水號主鍵',
    `created_at` datetime(3) DEFAULT NULL COMMENT '建立時間',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '修改時間',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '刪除時間',
    `articles_id` bigint unsigned NOT NULL COMMENT '文章的流水號主鍵',
    `tag` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '標籤',
    PRIMARY KEY (`id`),
    INDEX `idx_articles_id` (`articles_id`),
    INDEX `idx_tag` (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章的標籤';

-- 作者資訊 Table
DROP TABLE IF EXISTS `authors`;
CREATE TABle `authors` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '流水號主鍵',
    `created_at` datetime(3) DEFAULT NULL COMMENT '建立時間',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '修改時間',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '刪除時間',
    `title` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '標題',
    `content` TEXT COMMENT '內容',
    `status` TINYINT(3) NOT NULL DEFAULT '1' COMMENT '狀態 0:禁用, 1:啟用',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作者';
