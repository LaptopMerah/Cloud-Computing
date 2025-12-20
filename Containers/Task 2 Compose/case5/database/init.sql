-- Create database and table for link shortener
CREATE DATABASE IF NOT EXISTS `shorten_link` CHARACTER
SET
    utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `shorten_link`;

CREATE TABLE IF NOT EXISTS `links` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(64) NOT NULL,
    `url` TEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_code` (`code`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
