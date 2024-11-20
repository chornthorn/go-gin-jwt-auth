ALTER TABLE users
    DROP COLUMN role,
    DROP COLUMN is_active;

DROP TYPE user_role;