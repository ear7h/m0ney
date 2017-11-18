DROP SCHEMA IF EXISTS money_test;
CREATE SCHEMA money_test;
USE money_test;


DROP TABLE IF EXISTS `partitions`;
CREATE TABLE `partitions` (
	`name`    VARCHAR(64) NOT NULL,
	`week_of` DATETIME    NOT NULL,
	UNIQUE (`week_of`),
	unique KEY(`name`)
)
	ENGINE = InnoDB
	DEFAULT CHARSET = utf8;

DROP TABLE IF EXISTS `runs`;
CREATE TABLE `runs` (
	`id`           INT(11)    NOT NULL AUTO_INCREMENT,
	`symbol`       VARCHAR(8) NOT NULL,
	`start`        DATETIME   NOT NULL,
	`end`          DATETIME   NOT NULL,
	`partition_name` VARCHAR(64)   NOT NULL,
	FOREIGN KEY (`partition_name`) REFERENCES `partitions` (`name`),
	PRIMARY KEY (`id`)
)
	ENGINE = InnoDB
	AUTO_INCREMENT = 0
	DEFAULT CHARSET = utf8;