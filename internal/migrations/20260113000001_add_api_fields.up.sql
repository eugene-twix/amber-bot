-- Migration: Add fields for Mini App API
-- - soft-delete: deleted_at, deleted_by
-- - metadata: updated_at, updated_by, created_by (members)
-- - optimistic locking: version

-- =====================
-- TEAMS
-- =====================
ALTER TABLE teams ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS updated_by BIGINT;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS deleted_by BIGINT;
ALTER TABLE teams ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

-- Partial index for soft-delete queries
CREATE INDEX IF NOT EXISTS idx_teams_not_deleted ON teams(id) WHERE deleted_at IS NULL;

-- =====================
-- MEMBERS
-- =====================
ALTER TABLE members ADD COLUMN IF NOT EXISTS created_by BIGINT;
ALTER TABLE members ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
ALTER TABLE members ADD COLUMN IF NOT EXISTS updated_by BIGINT;
ALTER TABLE members ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE members ADD COLUMN IF NOT EXISTS deleted_by BIGINT;
ALTER TABLE members ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

-- Partial index for soft-delete queries
CREATE INDEX IF NOT EXISTS idx_members_not_deleted ON members(id) WHERE deleted_at IS NULL;
-- Composite index for team members lookup (excluding deleted)
CREATE INDEX IF NOT EXISTS idx_members_team_not_deleted ON members(team_id) WHERE deleted_at IS NULL;

-- =====================
-- TOURNAMENTS
-- =====================
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS updated_by BIGINT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS deleted_by BIGINT;
ALTER TABLE tournaments ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

-- Partial index for soft-delete queries
CREATE INDEX IF NOT EXISTS idx_tournaments_not_deleted ON tournaments(id) WHERE deleted_at IS NULL;
-- Index for date sorting (excluding deleted)
CREATE INDEX IF NOT EXISTS idx_tournaments_date_not_deleted ON tournaments(date DESC) WHERE deleted_at IS NULL;

-- =====================
-- RESULTS
-- =====================
ALTER TABLE results ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
ALTER TABLE results ADD COLUMN IF NOT EXISTS updated_by BIGINT;
ALTER TABLE results ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE results ADD COLUMN IF NOT EXISTS deleted_by BIGINT;
ALTER TABLE results ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

-- Partial index for soft-delete queries
CREATE INDEX IF NOT EXISTS idx_results_not_deleted ON results(id) WHERE deleted_at IS NULL;
-- Composite indexes for lookups (excluding deleted)
CREATE INDEX IF NOT EXISTS idx_results_tournament_not_deleted ON results(tournament_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_results_team_not_deleted ON results(team_id) WHERE deleted_at IS NULL;
