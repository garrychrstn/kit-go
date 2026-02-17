BEGIN;
CREATE TABLE IF NOT EXISTS rent_contracts (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    of_rent UUID NOT NULL references rents(id),
    of_store UUID NOT NULL references stores(id),
    of_user UUID NOT NULL references users(id),
    fileId VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);
COMMIT;
