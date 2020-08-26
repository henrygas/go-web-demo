CREATE DATABASE IF NOT EXISTS `web` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `web`.`userinfo` (
    `uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL DEFAULT '',
    `password` varchar(60) NOT NULL DEFAULT '',
    `departname` varchar(64) NOT NULL DEFAULT '',
    `created` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `web`.`userdetail` (
    `uid` int(10) unsigned NOT NULL DEFAULT 0,
    `intro` text NOT NULL,
    `profile` text NOT NULL,
    PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;