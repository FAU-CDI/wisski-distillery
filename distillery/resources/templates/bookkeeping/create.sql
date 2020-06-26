-- create the database if it doesn't exist yet
CREATE DATABASE IF NOT EXISTS `${DATABASE}`;

-- use the database we just created
USE `${DATABASE}`;

-- create the bookkeeping table
CREATE TABLE IF NOT EXISTS `${TABLE}`(
    -- automatically created fields
    pk INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- slug of the website
    slug TEXT NOT NULL UNIQUE,

    -- local file path
    filesystem_base TEXT NOT NULL,
    
    -- sql access credentials
    sql_database TEXT NOT NULL,
    sql_user TEXT NOT NULL,
    sql_password TEXT NOT NULL,

    -- graphdb credentials
    graphdb_repository TEXT NOT NULL,
    graphdb_user TEXT NOT NULL,
    graphdb_password TEXT NOT NULL
);