CREATE TYPE user_role AS ENUM ('user', 'admin', 'super_admin');

ALTER TABLE users
    ADD COLUMN role user_role NOT NULL DEFAULT 'user',
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;

CREATE INDEX idx_users_role ON users(role);