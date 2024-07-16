# This DDL (Data Definition Language) file contains all the scripts needed to define and manage
# the structure of the database for Gamabunta project. It includes commands to create the database,
# tables, indexes, and constraints that ensure data integrity and optimize performance.

# -- COMMAND TO CREATE THE DATABASE
CREATE DATABASE IF NOT EXISTS gamabunta;

# -- COMMAND TO CREATE THE user TABLE
CREATE TABLE IF NOT EXISTS user
(
    id       INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,
    username VARCHAR(500)                       NOT NULL,
    password VARCHAR(500)                       NOT NULL
);