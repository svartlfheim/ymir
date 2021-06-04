-- Create the migrations user
CREATE ROLE ymir_migrator
WITH 
    NOINHERIT
    LOGIN
    PASSWORD 'iammigrator';

-- Create the ymir db
CREATE SCHEMA ymir AUTHORIZATION ymir_migrator;

-- Create the app user
CREATE ROLE ymir_app
WITH 
    NOINHERIT
    LOGIN
    PASSWORD 'iamapp';

-- Ensure we can connect to the db and schema
GRANT CONNECT ON DATABASE postgres to ymir_app;
GRANT USAGE ON SCHEMA ymir to ymir_app;

-- Add default privileges for ymir_app
-- When ymir_migrator creates a table, ymir_app gets privileges on it
ALTER DEFAULT PRIVILEGES
FOR USER ymir_migrator
IN SCHEMA ymir
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO ymir_app;

