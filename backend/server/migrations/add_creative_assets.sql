-- Migration: Add creative_assets table
-- Description: Stores AI-generated images, diagrams, and videos with base64 data
-- Date: 2024-12-11

CREATE TABLE IF NOT EXISTS creative_assets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    asset_type VARCHAR(50) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    b64_data TEXT,
    url VARCHAR(512),
    prompt TEXT,
    provider VARCHAR(50),
    format VARCHAR(20),
    width INT,
    height INT,
    size_bytes BIGINT,
    diagram_code TEXT,
    diagram_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_creative_assets_user_id ON creative_assets(user_id);
CREATE INDEX IF NOT EXISTS idx_creative_assets_entity ON creative_assets(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_creative_assets_asset_type ON creative_assets(asset_type);
CREATE INDEX IF NOT EXISTS idx_creative_assets_purpose ON creative_assets(purpose);

-- Composite index for entity + purpose lookups (used by GetAssetByEntityAndPurpose)
CREATE INDEX IF NOT EXISTS idx_creative_assets_entity_purpose ON creative_assets(entity_type, entity_id, purpose);

-- Add comment to table
COMMENT ON TABLE creative_assets IS 'Stores AI-generated creative assets (images, diagrams, videos) with base64 data and metadata';

