BEGIN;
CREATE TABLE products (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000),
    price INTEGER DEFAULT 0 NOT NULL,
    category VARCHAR(50),
    specification JSONB DEFAULT '{}' NOT NULL,
    of_store UUID REFERENCES stores(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
CREATE TABLE product_availabilities (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    of_product UUID REFERENCES products(id) NOT NULL,
    of_store UUID REFERENCES stores(id) NOT NULL,
    of_branch UUID REFERENCES branches(id) NOT NULL,
    price INTEGER DEFAULT 0 NOT NULL,
    quantity INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
COMMIT;
