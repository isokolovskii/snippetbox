CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'web'@'%' IDENTIFIED WITH mysql_native_password BY 'pass';
GRANT ALL PRIVILEGES ON snippetbox.* TO 'web'@'%';
