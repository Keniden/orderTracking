CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE "orders" (
    track_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id SERIAL,
    title TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    from_location TEXT,
    to_location TEXT,
    status TEXT NOT NULL
);