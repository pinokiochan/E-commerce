CREATE TABLE inventory (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK(price > 0),
    available integer DEFAULT 0 CHECK(available >= 0),
    isdeleted BOOLEAN DEFAULT 'FALSE',
    version integer NOT NULL DEFAULT 1
);