-- Init script for local database in docker compose --
-- Create database for app --
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create user for app --
CREATE USER 'web'@'%' IDENTIFIED WITH mysql_native_password BY 'pass';
-- Grant privileges for user to manipulate database --
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, INDEX, ALTER ON snippetbox.* TO 'web'@'%';
