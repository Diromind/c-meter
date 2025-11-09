-- Rollback initial schema

DROP INDEX IF EXISTS idx_records_product_uuid;
DROP INDEX IF EXISTS idx_records_created_at;
DROP INDEX IF EXISTS idx_records_login;
DROP TABLE IF EXISTS records;
DROP TABLE IF EXISTS product_details;
