BEGIN;
    CREATE TABLE roles (
        id UUID DEFAULT uuidv7() PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        of_store UUID NOT NULL REFERENCES stores(id),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );
    CREATE TABLE permissions (
        id UUID DEFAULT uuidv7() PRIMARY KEY,
        action VARCHAR(255) NOT NULL UNIQUE,
        description VARCHAR(1000)
    );
    CREATE TABLE rel_role_permission (
        of_role UUID NOT NULL REFERENCES roles(id),
        of_permission UUID NOT NULL REFERENCES permissions(id),
        PRIMARY KEY (of_role, of_permission)
    );
    CREATE TABLE users (
        id UUID DEFAULT uuidv7() PRIMARY KEY,
        username VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        phone_number VARCHAR(255) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        of_store UUID NOT NULL REFERENCES stores(id),
        of_role UUID NOT NULL REFERENCES roles(id),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );
    ALTER TABLE stores ADD COLUMN of_owner UUID DEFAULT uuidv7() REFERENCES users(id);
COMMIT
