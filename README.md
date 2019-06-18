#数据库导出
1. mysql\bin添加进环境变量
2. mysqldump -h localhost -u root -p nuuinfo >d:mydb.sql

#数据库语句
CREATE SCHEMA `nuu_db` ;

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
  ) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='user profile table'