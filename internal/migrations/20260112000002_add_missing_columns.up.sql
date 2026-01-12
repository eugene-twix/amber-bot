-- Add created_by column to teams
ALTER TABLE teams ADD COLUMN IF NOT EXISTS created_by BIGINT;

-- Add location and created_by columns to tournaments
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS location VARCHAR(255) DEFAULT '';
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS created_by BIGINT;
