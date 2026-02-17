BEGIN;

CREATE TABLE stores (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000),
    logo VARCHAR(1000),
    map VARCHAR(1000),
    coordinate VARCHAR(100),
    address VARCHAR(1000) NOT NULL,
    phone VARCHAR(1000) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    category VARCHAR(20)[] NOT NULL DEFAULT '{}',
    contacts JSONB DEFAULT '{}' NOT NULL,
    term_and_service TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE branches (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR(1000) NOT NULL,
    code VARCHAR(20),
    map VARCHAR(1000),
    address VARCHAR(1000),
    description VARCHAR(1000),
    of_store UUID REFERENCES stores(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

COMMIT;
