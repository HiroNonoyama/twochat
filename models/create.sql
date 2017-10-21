-- pg_ctl start && psql の後に

-- create database, user, table
CREATE DATABASE twochat;
CREATE USER twochat_client;
GRANT ALL ON DATABASE twochat TO twochat_client;
/c twochat twochat_client;


-- users, messages
CREATE TABLE users(
    id INT,
    name VARCHAR,
    icon_image VARCHAR,
    PRIMARY KEY (id)
);

CREATE TABLE messages(
    id INT,
    sender_id INT,
    message VARCHAR,
    created_at VARCHAR,
    PRIMARY KEY (id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);

