#数据库导出
1. mysql\bin添加进环境变量
2. mysqldump -h localhost -u root -p nuuinfo >d:mydb.sql

#数据库语句
CREATE SCHEMA `nuuinfo` ;
CREATE SCHEMA `nuu_db` ;


CREATE TABLE `nuuinfo`.`person` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` VARCHAR(45) NOT NULL,
  `user_name` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NULL,
  `gender` INT ZEROFILL NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `user_id_UNIQUE` (`user_id` ASC) VISIBLE,
  UNIQUE INDEX `user_name_UNIQUE` (`user_name` ASC) VISIBLE);
  
  
  
  CREATE TABLE `account` (
    `account_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `uuid` varchar(45) NOT NULL,
    `username` varchar(45) NOT NULL,
    `email` varchar(45) NOT NULL,
    `mobile` varchar(45) NOT NULL,
    `iso` char(5) NOT NULL,
    `password` varchar(45) NOT NULL,
    PRIMARY KEY (`account_id`),
    UNIQUE KEY `user_id_UNIQUE` (`uuid`),
    UNIQUE KEY `username_UNIQUE` (`username`),
    UNIQUE KEY `email_UNIQUE` (`email`),
    UNIQUE KEY `mobile_UNIQUE` (`mobile`)
  ) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='user profile table'