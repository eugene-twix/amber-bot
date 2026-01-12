-- Remove created_by column from teams
ALTER TABLE teams DROP COLUMN IF EXISTS created_by;

-- Remove location and created_by columns from tournaments
ALTER TABLE tournaments DROP COLUMN IF EXISTS location;
ALTER TABLE tournaments DROP COLUMN IF EXISTS created_by;
