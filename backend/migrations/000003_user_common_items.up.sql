-- User common items with hierarchical structure using ltree

-- Enable ltree extension for hierarchical data
CREATE EXTENSION IF NOT EXISTS ltree;

-- User common items table
CREATE TABLE user_common_items (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) NOT NULL,
    path ltree NOT NULL,
    name VARCHAR(255) NOT NULL,
    product_uuid UUID REFERENCES product_details(uuid) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(login, path)
);

-- Indexes for efficient queries
CREATE INDEX idx_user_common_items_login ON user_common_items(login);
CREATE INDEX idx_user_common_items_path_gist ON user_common_items USING GIST(path);
CREATE INDEX idx_user_common_items_login_path ON user_common_items(login, path);
