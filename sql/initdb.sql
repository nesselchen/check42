CREATE DATABASE IF NOT EXISTS `check42` ;

USE `check42`;

CREATE TABLE IF NOT EXISTS `user` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(50) NOT NULL,
    `name` VARCHAR(140),
    primary key (`id`)
);

CREATE TABLE IF NOT EXISTS `todo`(
    `id` INT NOT NULL AUTO_INCREMENT,
    `owner` INT NOT NULL,
    `text` VARCHAR(140),
    `done` BOOLEAN DEFAULT 0,
    `due` DATETIME NULL,
    `created` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`owner`) REFERENCES user(`id`)
);

INSERT INTO user (`name`, `email`) VALUES ("admin", "admin@adm.in");