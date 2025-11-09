-- Initial schema for c-meter

-- Set timezone
SET timezone = 'Europe/Moscow';

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Product details table
CREATE TABLE product_details (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    ccal BIGINT NOT NULL CHECK (ccal > 0),
    fats BIGINT NOT NULL CHECK (fats >= 0),
    proteins BIGINT NOT NULL CHECK (proteins >= 0),
    carbs BIGINT NOT NULL CHECK (carbs >= 0)
);

-- Records table
CREATE TABLE records (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_uuid UUID NOT NULL REFERENCES product_details(uuid) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    login VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_records_login ON records(login);
CREATE INDEX idx_records_created_at ON records(created_at);
CREATE INDEX idx_records_product_uuid ON records(product_uuid);
