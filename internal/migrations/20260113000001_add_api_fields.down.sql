-- Rollback: Remove API fields

-- Drop indexes first
DROP INDEX IF EXISTS idx_teams_not_deleted;
DROP INDEX IF EXISTS idx_members_not_deleted;
DROP INDEX IF EXISTS idx_members_team_not_deleted;
DROP INDEX IF EXISTS idx_tournaments_not_deleted;
DROP INDEX IF EXISTS idx_tournaments_date_not_deleted;
DROP INDEX IF EXISTS idx_results_not_deleted;
DROP INDEX IF EXISTS idx_results_tournament_not_deleted;
DROP INDEX IF EXISTS idx_results_team_not_deleted;

-- TEAMS
ALTER TABLE teams DROP COLUMN IF EXISTS updated_at;
ALTER TABLE teams DROP COLUMN IF EXISTS updated_by;
ALTER TABLE teams DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE teams DROP COLUMN IF EXISTS deleted_by;
ALTER TABLE teams DROP COLUMN IF EXISTS version;

-- MEMBERS
ALTER TABLE members DROP COLUMN IF EXISTS created_by;
ALTER TABLE members DROP COLUMN IF EXISTS updated_at;
ALTER TABLE members DROP COLUMN IF EXISTS updated_by;
ALTER TABLE members DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE members DROP COLUMN IF EXISTS deleted_by;
ALTER TABLE members DROP COLUMN IF EXISTS version;

-- TOURNAMENTS
ALTER TABLE tournaments DROP COLUMN IF EXISTS updated_at;
ALTER TABLE tournaments DROP COLUMN IF EXISTS updated_by;
ALTER TABLE tournaments DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE tournaments DROP COLUMN IF EXISTS deleted_by;
ALTER TABLE tournaments DROP COLUMN IF EXISTS version;

-- RESULTS
ALTER TABLE results DROP COLUMN IF EXISTS updated_at;
ALTER TABLE results DROP COLUMN IF EXISTS updated_by;
ALTER TABLE results DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE results DROP COLUMN IF EXISTS deleted_by;
ALTER TABLE results DROP COLUMN IF EXISTS version;
