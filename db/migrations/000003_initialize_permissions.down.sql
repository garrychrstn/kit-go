BEGIN;
    DROP TABLE IF EXISTS rel_role_permission;
    DROP TABLE IF EXISTS permissions;
    ALTER TABLE stores DROP COLUMN IF EXISTS of_owner;
    DROP TABLE IF EXISTS users;
    DROP TABLE IF EXISTS roles;
COMMIT
