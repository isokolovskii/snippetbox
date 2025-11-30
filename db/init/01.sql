CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'web'@'%' IDENTIFIED WITH mysql_native_password BY 'pass';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, ALTER ON snippetbox.* TO 'web'@'%';
