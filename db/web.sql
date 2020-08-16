CREATE DATABASE IF NOT EXISTS `web` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `web`.`userinfo` (
    `uid` INT(10) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL DEFAULT '',
    `departname` VARCHAR(64) NOT NULL DEFAULT '',
    `created` BIGINT(20) NOT NULL DEFAULT 0,
    PRIMARY KEY (`uid`)
);

CREATE TABLE IF NOT EXISTS `web`.`userdetail` (
    `uid` INT(10) NOT NULL DEFAULT 0,
    `intro` TEXT NOT NULL,
    `profile` TEXT NOT NULL,
    PRIMARY KEY (`uid`)
);