package migration

// @TODO: create an index on product, user, saving to improve ILIKE queries - most like a gin index using trigram
// @TODO: create btree index for other product fields

/*
BEGIN;
CREATE EXTENSION IF NOT EXISTS postgres_fdw;
DROP SERVER IF EXISTS core CASCADE;
CREATE SERVER core FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'core', port '5432');
CREATE USER MAPPING IF NOT EXISTS FOR CURRENT_USER SERVER core OPTIONS (user 'postgres', password 'postgres');
DROP SCHEMA IF EXISTS core CASCADE;
CREATE SCHEMA core;
IMPORT FOREIGN SCHEMA public FROM SERVER core INTO core;
COMMIT;
*/
