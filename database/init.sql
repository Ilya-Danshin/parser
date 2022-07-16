DROP DATABASE IF EXISTS parser;
CREATE DATABASE parser;

\c parser;

DROP TABLE IF EXISTS goods;
CREATE TABLE goods
(
    id INT PRIMARY KEY,
    name TEXT,
    url TEXT,
    url_img TEXT,
    price FLOAT
);

