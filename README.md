#数据库导出
1. mysql\bin添加进环境变量
2. mysqldump -h localhost -u root -p nuuinfo >d:mydb.sql

#数据库语句
CREATE SCHEMA `nuuinfo` ;


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