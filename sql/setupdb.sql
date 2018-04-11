DROP DATABASE IF EXISTS belajar;
CREATE DATABASE IF NOT EXISTS belajar;

DROP USER IF EXISTS belajar;
CREATE USER belajar IDENTIFIED BY 'pass';
GRANT ALL ON belajar.* TO belajar;

USE belajar;

CREATE TABLE IF NOT EXISTS satu(
	name VARCHAR (20),
	age INT
);
