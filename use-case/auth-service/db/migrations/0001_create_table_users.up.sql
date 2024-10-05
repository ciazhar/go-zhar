CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Ensure the uuid-ossp extension is enabled

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
